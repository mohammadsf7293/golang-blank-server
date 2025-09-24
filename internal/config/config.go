package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBConfig DBConfig
	Server   ServerConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type ServerConfig struct {
	Port string
}

func New() *Config {
	return &Config{
		DBConfig: DBConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "3306"),
			User:     getEnvOrDefault("DB_USER", "app_user"),
			Password: getEnvOrDefault("DB_PASSWORD", "app_pass"),
			DBName:   getEnvOrDefault("DB_NAME", "blank_project"),
		},
		Server: ServerConfig{
			Port: getEnvOrDefault("SERVER_PORT", "8080"),
		},
	}
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}