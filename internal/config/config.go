package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

// Config holds application configuration.
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Logging  LoggingConfig  `json:"logging"`
}

// ServerConfig holds server-related configuration.
type ServerConfig struct {
	Host         string `json:"host"`
	Environment  string `json:"environment"`
	Port         int    `json:"port"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}

// DatabaseConfig holds database-related configuration.
type DatabaseConfig struct {
	Type             string `json:"type"`
	ConnectionString string `json:"connection_string"`
	MigrationsPath   string `json:"migrations_path"`
	MaxOpenConns     int    `json:"max_open_conns"`
	MaxIdleConns     int    `json:"max_idle_conns"`
}

// LoggingConfig holds logging-related configuration.
type LoggingConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
	Output string `json:"output"`
}

// LoadConfig loads configuration from environment variables with defaults.
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "localhost"),
			Port:         getEnvAsInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvAsInt("SERVER_READ_TIMEOUT", 30),
			WriteTimeout: getEnvAsInt("SERVER_WRITE_TIMEOUT", 30),
			Environment:  getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			Type:             getEnv("DB_TYPE", "sqlite"),
			ConnectionString: getEnv("DB_CONNECTION_STRING", "./data/portfolio.db"),
			MaxOpenConns:     getEnvAsInt("DB_MAX_OPEN_CONNS", 10),
			MaxIdleConns:     getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			MigrationsPath:   getEnv("DB_MIGRATIONS_PATH", "./migrations"),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
			Output: getEnv("LOG_OUTPUT", "stdout"),
		},
	}
}

// Address returns the server address in the format host:port.
func (s *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// IsDevelopment returns true if the environment is development.
func (s *ServerConfig) IsDevelopment() bool {
	return s.Environment == "development"
}

// IsProduction returns true if the environment is production.
func (s *ServerConfig) IsProduction() bool {
	return s.Environment == "production"
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.Server.ReadTimeout <= 0 {
		return fmt.Errorf("invalid read timeout: %d", c.Server.ReadTimeout)
	}

	if c.Server.WriteTimeout <= 0 {
		return fmt.Errorf("invalid write timeout: %d", c.Server.WriteTimeout)
	}

	if c.Database.Type == "" {
		return errors.New("database type cannot be empty")
	}

	if c.Database.ConnectionString == "" {
		return errors.New("database connection string cannot be empty")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLogLevels[c.Logging.Level] {
		return fmt.Errorf("invalid log level: %s", c.Logging.Level)
	}

	return nil
}

// getEnv gets an environment variable with a default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

// getEnvAsInt gets an environment variable as an integer with a default value.
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}

	return defaultValue
}
