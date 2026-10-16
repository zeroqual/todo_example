package todo_repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	core_errors "github.com/zeroqual/todo_example/internal/core/errors"
	core_models "github.com/zeroqual/todo_example/internal/core/models"
)

type Repository interface {
	CreateTodo(ctx context.Context, name, description string) (string, error)
	GetTodoByID(ctx context.Context, todoID string) (*core_models.Todo, error)
	SetTodoCompleted(ctx context.Context, todoID string) (bool, error)
	DeleteTodoByID(ctx context.Context, todoID string) error
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) PostgresRepository {
	return PostgresRepository{
		pool: pool,
	}
}

func (r *PostgresRepository) DeleteTodoByID(ctx context.Context, todoID string) error {
	tag, err := r.pool.Exec(
		ctx,
		`
			DELETE FROM todos
			WHERE todo_id = $1
		`,
		todoID,
	)

	if tag.RowsAffected() == 0 {
		return core_errors.ErrTodoNotFound
	}

	return err
}

func (r *PostgresRepository) SetTodoCompleted(ctx context.Context, todoID string) (bool, error) {

	rows, err := r.pool.Exec(
		ctx,
		`
			UPDATE todos
			SET is_completed = true
			WHERE is_completed = false AND todo_id = $1
		`,
		todoID,
	)

	if err != nil {
		return false, err
	}

	return rows.RowsAffected() > 0, nil
}

func (r *PostgresRepository) GetTodoByID(ctx context.Context, todoID string) (*core_models.Todo, error) {
	todo := &core_models.Todo{}

	err := r.pool.QueryRow(
		ctx,
		`
			SELECT todo_id, todo_name,todo_description,is_completed, created_at FROM todos
			WHERE todo_id = $1
		`,
		todoID,
	).Scan(
		&todo.TodoID,
		&todo.TodoName,
		&todo.TodoDescription,
		&todo.IsCompleted,
		&todo.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (r *PostgresRepository) CreateTodo(ctx context.Context, name, description string) (string, error) {
	var id string
	err := r.pool.QueryRow(
		ctx,
		`
			INSERT INTO todos(todo_name, todo_description)
			VALUES($1,$2)
			RETURNING todo_id
		`,
		name, description,
	).Scan(&id)

	return id, err
}
