package core_database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	core_config "github.com/zeroqual/todo_example/internal/core/config"
)

func InitPostgres(dbConfig *core_config.DatabaseConfig) (*pgxpool.Pool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbConfig.GetDsn())
	if err != nil {
		return nil, fmt.Errorf("failed to create new pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping pool: %v", err)
	}

	return pool, nil
}
