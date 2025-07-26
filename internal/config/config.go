// Package config provides configuration loading and validation for the portfolio website.
// It handles environment variable parsing, defaults, and configuration validation
// to ensure the application starts with valid settings.
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"holger-hahn-website/internal/constants"
)

// Static error variables to avoid dynamic error creation.
var (
	ErrInvalidServerPort     = errors.New("invalid server port")
	ErrInvalidReadTimeout    = errors.New("invalid read timeout")
	ErrInvalidWriteTimeout   = errors.New("invalid write timeout")
	ErrEmptyDatabaseType     = errors.New("database type cannot be empty")
	ErrEmptyConnectionString = errors.New("database connection string cannot be empty")
	ErrInvalidLogLevel       = errors.New("invalid log level")
)

// Config holds application configuration.
type Config struct {
	Logging  LoggingConfig  `json:"logging"`
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
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
			Port:         getEnvAsInt("SERVER_PORT", constants.DefaultServerPort),
			ReadTimeout:  getEnvAsInt("SERVER_READ_TIMEOUT", constants.DefaultReadTimeoutSeconds),
			WriteTimeout: getEnvAsInt("SERVER_WRITE_TIMEOUT", constants.DefaultWriteTimeoutSeconds),
			Environment:  getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			Type:             getEnv("DB_TYPE", "sqlite"),
			ConnectionString: getEnv("DB_CONNECTION_STRING", "./data/portfolio.db"),
			MaxOpenConns:     getEnvAsInt("DB_MAX_OPEN_CONNS", constants.DefaultMaxOpenConnections),
			MaxIdleConns:     getEnvAsInt("DB_MAX_IDLE_CONNS", constants.DefaultMaxIdleConnections),
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
	if c.Server.Port < constants.MinValidPort || c.Server.Port > constants.MaxValidPort {
		return fmt.Errorf("%w: %d", ErrInvalidServerPort, c.Server.Port)
	}

	if c.Server.ReadTimeout <= 0 {
		return fmt.Errorf("%w: %d", ErrInvalidReadTimeout, c.Server.ReadTimeout)
	}

	if c.Server.WriteTimeout <= 0 {
		return fmt.Errorf("%w: %d", ErrInvalidWriteTimeout, c.Server.WriteTimeout)
	}

	if c.Database.Type == "" {
		return ErrEmptyDatabaseType
	}

	if c.Database.ConnectionString == "" {
		return ErrEmptyConnectionString
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLogLevels[c.Logging.Level] {
		return fmt.Errorf("%w: %s", ErrInvalidLogLevel, c.Logging.Level)
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
