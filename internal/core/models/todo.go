package core_models

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	TodoID          uuid.UUID `json:"todo_id"`
	TodoName        string    `json:"todo_name"`
	TodoDescription string    `json:"todo_description"`
	IsCompleted     bool      `json:"is_completed"`
	CreatedAt       time.Time `json:"created_at"`
}
