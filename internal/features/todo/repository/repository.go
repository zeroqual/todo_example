package todo_repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateTodo(ctx context.Context, name, description string) (string, error)
	// GetAllTodos(ctx context.Context) ([]*core_models.Todo, error)
	// GetTodoByID(ctx context.Context) (*core_models.Todo, error)
	// DeleteTodoByID(ctx context.Context) error
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) PostgresRepository {
	return PostgresRepository{
		pool: pool,
	}
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
