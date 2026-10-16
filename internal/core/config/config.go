package core_config

import "os"

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

func LoadConfig() *Config {
	return &Config{
		Server:   loadServerConfig(),
		Database: loadDbConfig(),
	}
}

type ServerConfig struct {
	Host string
	Port string
}

func loadServerConfig() ServerConfig {
	return ServerConfig{
		Host: getEnv("SERVER_HOST", ""),
		Port: getEnv("SERVER_PORT", ""),
	}
}

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

func loadDbConfig() DatabaseConfig {
	return DatabaseConfig{
		User:     getEnv("POSTGRES_USER", ""),
		Password: getEnv("POSTGRES_PASSWORD", ""),
		Host:     getEnv("POSTGRES_HOST", ""),
		Port:     getEnv("POSTGRES_PORT", ""),
		Name:     getEnv("POSTGRES_NAME", ""),
	}
}

func getEnv(key, def string) string {
	str := os.Getenv(key)
	if str == "" {
		return def
	}
	return str
}
