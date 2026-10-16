package todo_handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	core_errors "github.com/zeroqual/todo_example/internal/core/errors"
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

func (h *Handler) DeleteTodo(c *gin.Context) {
	var dto todo_dto.GetTodoRequest

	if err := c.ShouldBindJSON(&dto); err != nil {
		slog.Error("failed to bind json", "error", err)
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	err := h.service.DeleteTodoByID(c.Request.Context(), dto.TodoID)
	if err != nil {
		slog.Error("failed to delete todo", "error", err)
		switch {
		case errors.Is(err, core_errors.ErrTodoNotFound):
			c.JSON(http.StatusBadRequest, "todo not found")
			return
		default:
			c.JSON(http.StatusInternalServerError, "internal server error")
			return
		}
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusFound, gin.H{"success": true})
}

func (h *Handler) SetTodoCompleted(c *gin.Context) {
	var dto todo_dto.GetTodoRequest

	if err := c.ShouldBindJSON(&dto); err != nil {
		slog.Error("failed to bind json", "error", err)
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	err := h.service.SetTodoCompleted(c.Request.Context(), dto.TodoID)
	if err != nil {
		slog.Error("failed to set todo completed", "error", err)
		switch {
		case errors.Is(err, core_errors.ErrTodoNotUpdated):
			c.JSON(http.StatusBadRequest, "todo not updated")
			return
		default:
			c.JSON(http.StatusInternalServerError, "internal server error")
			return
		}
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusFound, gin.H{"success": true})
}

func (h *Handler) GetTodoByID(c *gin.Context) {
	var dto todo_dto.GetTodoRequest

	if err := c.ShouldBindJSON(&dto); err != nil {
		slog.Error("failed to bind json", "error", err)
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	//

	todo, err := h.service.GetTodoByID(c.Request.Context(), dto.TodoID)
	if err != nil {
		slog.Error("failed to get todo", "error", err)
		switch {
		case errors.Is(err, core_errors.ErrTodoNotFound):
			c.JSON(http.StatusNotFound, "todo not found")
			return
		default:
			c.JSON(http.StatusInternalServerError, "internal server error")
			return
		}
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, gin.H{"todo": todo})
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
		slog.Error("failed to create todo", "error", err)
		switch {
		case errors.Is(err, core_errors.ErrUniqueTodo):
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
