package todo_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	todo_dto "github.com/zeroqual/todo_example/internal/features/todo"
	todo_repository "github.com/zeroqual/todo_example/internal/features/todo/repository"
)

var (
	ErrUniqueTodo = errors.New("todo already exists")
)

type Service struct {
	repo todo_repository.Repository
}

func NewService(repo todo_repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateTodo(ctx context.Context, dto todo_dto.CreateTodoRequest) (string, error) {
	//тут оп сути должна быть бизнес логика, но проверка происходит на уровне дтошки так что плохие данные не пройдут и тут их можно не проверить
	//так что тут сервис это просто заглушка

	todoID, err := s.repo.CreateTodo(ctx, dto.TodoName, dto.TodoDescription)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == "23505" {
			return "", fmt.Errorf("%w:%v", ErrUniqueTodo, err)
		}
		return "", fmt.Errorf("unexpected error: %v", err)
	}

	return todoID, nil
}
