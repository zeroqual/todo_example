package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	core_config "github.com/zeroqual/todo_example/internal/core/config"
	core_database "github.com/zeroqual/todo_example/internal/core/database"
	core_server "github.com/zeroqual/todo_example/internal/core/server"
)

func main() {
	//init config
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	cfg := core_config.LoadConfig()
	fmt.Println(cfg)
	//init db pool

	pool, err := core_database.InitPostgres(&cfg.Database)
	if err != nil {
		log.Fatalf("failed to load pool: %v", err)
	}

	//start server
	if err := core_server.Start(cfg, pool); err != nil {
		log.Fatalf("server start error: %v", err)
	}
}
