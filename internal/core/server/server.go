package core_server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	core_config "github.com/zeroqual/todo_example/internal/core/config"
	todo_handler "github.com/zeroqual/todo_example/internal/features/todo/handler"
	todo_repository "github.com/zeroqual/todo_example/internal/features/todo/repository"
	todo_service "github.com/zeroqual/todo_example/internal/features/todo/service"
)

func loader(pool *pgxpool.Pool) todo_handler.Handler {
	// repo
	todoRepository := todo_repository.NewPostgresRepository(pool)

	//services
	todoService := todo_service.NewService(&todoRepository)

	//handlers
	todoHandler := todo_handler.NewHandler(todoService)

	return todoHandler
}

func Start(cfg *core_config.Config, pool *pgxpool.Pool) error {
	r := gin.New()

	todoHandler := loader(pool)

	publicRoutes := r.Group("/todo")
	{
		//routes here
		publicRoutes.POST("/create", todoHandler.CreateTodo)
		publicRoutes.GET("/get", todoHandler.GetTodoByID)
		publicRoutes.POST("/complete", todoHandler.SetTodoCompleted)
		publicRoutes.DELETE("/delete", todoHandler.DeleteTodo)
	}

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler: r,
	}

	serverError := make(chan error, 1)

	go func() {
		slog.Info("start server", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverError <- err
		}
	}()

	closer := make(chan os.Signal, 1)

	signal.Notify(closer, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-serverError:
		pool.Close()
		return fmt.Errorf("failed to start server: %v", err)
	case sig := <-closer:
		slog.Info("got signal", "signal", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			pool.Close()
			return server.Close()
		}

		// gracefull shutdown done
		pool.Close()

		slog.Info("successfully shutdown server")
	}

	return nil
}
