package core_errors

import "errors"

var (
	ErrUniqueTodo     = errors.New("todo already exists")
	ErrTodoNotFound   = errors.New("todo not found")
	ErrTodoNotUpdated = errors.New("todo not updated")
)
