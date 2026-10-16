package todo_handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	todo_dto "github.com/zeroqual/todo_example/internal/features/todo"
	todo_service "github.com/zeroqual/todo_example/internal/features/todo/service"
)

type Handler struct {
	service *todo_service.Service
}

func NewHandler(service *todo_service.Service) Handler {
	return Handler{
		service: service,
	}
}

func (h *Handler) CreateTodo(c *gin.Context) {
	var dto todo_dto.CreateTodoRequest

	if err := c.ShouldBindJSON(&dto); err != nil {
		slog.Error("failed to bind json", "error", err)
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	//

	todoID, err := h.service.CreateTodo(c.Request.Context(), dto)
	if err != nil {
		switch {
		case errors.Is(err, todo_service.ErrUniqueTodo):
			c.JSON(http.StatusConflict, "todo already exists")
			return
		default:
			c.JSON(http.StatusInternalServerError, "internal server error")
			return
		}
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, gin.H{"todo_id": todoID})
}

func (h *Handler) GetTodos(c *gin.Context) {

}
