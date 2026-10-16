package todo_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	core_errors "github.com/zeroqual/todo_example/internal/core/errors"
	core_models "github.com/zeroqual/todo_example/internal/core/models"
	todo_dto "github.com/zeroqual/todo_example/internal/features/todo"
	todo_repository "github.com/zeroqual/todo_example/internal/features/todo/repository"
)

type Service struct {
	repo todo_repository.Repository
}

func NewService(repo todo_repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) DeleteTodoByID(ctx context.Context, todoID string) error {
	if err := s.repo.DeleteTodoByID(ctx, todoID); err != nil {
		return fmt.Errorf("delete todo: %w", err)
	}
	return nil
}

func (s *Service) SetTodoCompleted(ctx context.Context, todoID string) error {

	updated, err := s.repo.SetTodoCompleted(ctx, todoID)
	if err != nil {
		return fmt.Errorf("failed to set todo completed: %v", err)
	}

	if !updated {
		return core_errors.ErrTodoNotUpdated
	}
	return nil
}

func (s *Service) GetTodoByID(ctx context.Context, todoID string) (*core_models.Todo, error) {

	todo, err := s.repo.GetTodoByID(ctx, todoID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w:%v", core_errors.ErrTodoNotFound, err)
		}
		return nil, fmt.Errorf("unexpected error while get todo by id: %v", err)
	}

	return todo, nil
}

func (s *Service) CreateTodo(ctx context.Context, dto todo_dto.CreateTodoRequest) (string, error) {
	todoID, err := s.repo.CreateTodo(ctx, dto.TodoName, dto.TodoDescription)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == "23505" {
			return "", fmt.Errorf("%w:%v", core_errors.ErrUniqueTodo, err)
		}
		return "", fmt.Errorf("unexpected error: %v", err)
	}

	return todoID, nil
}
