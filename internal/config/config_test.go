package config

import (
	"fmt"
	"os"
	"testing"

	"holger-hahn-website/internal/constants"
	"holger-hahn-website/internal/testutil"
)

func TestLoadConfig(t *testing.T) {
	t.Run("LoadConfig with defaults", func(t *testing.T) {
		// Clear environment variables to test defaults
		testutil.WithEnv(t, map[string]string{
			"SERVER_HOST":          "",
			"SERVER_PORT":          "",
			"SERVER_READ_TIMEOUT":  "",
			"SERVER_WRITE_TIMEOUT": "",
			"ENVIRONMENT":          "",
			"DB_TYPE":              "",
			"DB_CONNECTION_STRING": "",
			"DB_MAX_OPEN_CONNS":    "",
			"DB_MAX_IDLE_CONNS":    "",
			"DB_MIGRATIONS_PATH":   "",
			"LOG_LEVEL":            "",
			"LOG_FORMAT":           "",
			"LOG_OUTPUT":           "",
		}, func() {
			config := LoadConfig()

			// Test server defaults
			testutil.AssertEqual(t, "localhost", config.Server.Host)
			testutil.AssertEqual(t, constants.DefaultServerPort, config.Server.Port)
			testutil.AssertEqual(t, constants.DefaultReadTimeoutSeconds, config.Server.ReadTimeout)
			testutil.AssertEqual(t, constants.DefaultWriteTimeoutSeconds, config.Server.WriteTimeout)
			testutil.AssertEqual(t, "development", config.Server.Environment)

			// Test database defaults
			testutil.AssertEqual(t, "sqlite", config.Database.Type)
			testutil.AssertEqual(t, "./data/portfolio.db", config.Database.ConnectionString)
			testutil.AssertEqual(t, constants.DefaultMaxOpenConnections, config.Database.MaxOpenConns)
			testutil.AssertEqual(t, constants.DefaultMaxIdleConnections, config.Database.MaxIdleConns)
			testutil.AssertEqual(t, "./migrations", config.Database.MigrationsPath)

			// Test logging defaults
			testutil.AssertEqual(t, "info", config.Logging.Level)
			testutil.AssertEqual(t, "json", config.Logging.Format)
			testutil.AssertEqual(t, "stdout", config.Logging.Output)
		})
	})

	t.Run("LoadConfig with environment variables", func(t *testing.T) {
		testutil.WithEnv(t, map[string]string{
			"SERVER_HOST":          "0.0.0.0",
			"SERVER_PORT":          "3000",
			"SERVER_READ_TIMEOUT":  "45",
			"SERVER_WRITE_TIMEOUT": "60",
			"ENVIRONMENT":          "production",
			"DB_TYPE":              "postgres",
			"DB_CONNECTION_STRING": "postgres://user:pass@localhost/db",
			"DB_MAX_OPEN_CONNS":    "25",
			"DB_MAX_IDLE_CONNS":    "10",
			"DB_MIGRATIONS_PATH":   "./db/migrations",
			"LOG_LEVEL":            "debug",
			"LOG_FORMAT":           "text",
			"LOG_OUTPUT":           "file",
		}, func() {
			config := LoadConfig()

			// Test server configuration
			testutil.AssertEqual(t, "0.0.0.0", config.Server.Host)
			testutil.AssertEqual(t, 3000, config.Server.Port)
			testutil.AssertEqual(t, 45, config.Server.ReadTimeout)
			testutil.AssertEqual(t, 60, config.Server.WriteTimeout)
			testutil.AssertEqual(t, "production", config.Server.Environment)

			// Test database configuration
			testutil.AssertEqual(t, "postgres", config.Database.Type)
			testutil.AssertEqual(t, "postgres://user:pass@localhost/db", config.Database.ConnectionString)
			testutil.AssertEqual(t, 25, config.Database.MaxOpenConns)
			testutil.AssertEqual(t, 10, config.Database.MaxIdleConns)
			testutil.AssertEqual(t, "./db/migrations", config.Database.MigrationsPath)

			// Test logging configuration
			testutil.AssertEqual(t, "debug", config.Logging.Level)
			testutil.AssertEqual(t, "text", config.Logging.Format)
			testutil.AssertEqual(t, "file", config.Logging.Output)
		})
	})

	t.Run("LoadConfig with invalid environment values falls back to defaults", func(t *testing.T) {
		testutil.WithEnv(t, map[string]string{
			"SERVER_PORT":         "invalid_port",
			"SERVER_READ_TIMEOUT": "not_a_number",
			"DB_MAX_OPEN_CONNS":   "invalid_number",
		}, func() {
			config := LoadConfig()

			// Should fall back to defaults for invalid values
			testutil.AssertEqual(t, constants.DefaultServerPort, config.Server.Port)
			testutil.AssertEqual(t, constants.DefaultReadTimeoutSeconds, config.Server.ReadTimeout)
			testutil.AssertEqual(t, constants.DefaultMaxOpenConnections, config.Database.MaxOpenConns)
		})
	})
}

func TestServerConfig_Address(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		expected string
		port     int
	}{
		{
			name:     "localhost with default port",
			host:     "localhost",
			port:     8080,
			expected: "localhost:8080",
		},
		{
			name:     "all interfaces with custom port",
			host:     "0.0.0.0",
			port:     3000,
			expected: "0.0.0.0:3000",
		},
		{
			name:     "custom host with standard port",
			host:     "api.example.com",
			port:     80,
			expected: "api.example.com:80",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := ServerConfig{
				Host: tt.host,
				Port: tt.port,
			}
			result := server.Address()
			testutil.AssertEqual(t, tt.expected, result)
		})
	}
}

func TestServerConfig_IsDevelopment(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		expected    bool
	}{
		{
			name:        "development environment",
			environment: "development",
			expected:    true,
		},
		{
			name:        "production environment",
			environment: "production",
			expected:    false,
		},
		{
			name:        "staging environment",
			environment: "staging",
			expected:    false,
		},
		{
			name:        "empty environment",
			environment: "",
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := ServerConfig{
				Environment: tt.environment,
			}
			result := server.IsDevelopment()
			testutil.AssertEqual(t, tt.expected, result)
		})
	}
}

func TestServerConfig_IsProduction(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		expected    bool
	}{
		{
			name:        "production environment",
			environment: "production",
			expected:    true,
		},
		{
			name:        "development environment",
			environment: "development",
			expected:    false,
		},
		{
			name:        "staging environment",
			environment: "staging",
			expected:    false,
		},
		{
			name:        "empty environment",
			environment: "",
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := ServerConfig{
				Environment: tt.environment,
			}
			result := server.IsProduction()
			testutil.AssertEqual(t, tt.expected, result)
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	t.Run("valid configuration", func(t *testing.T) {
		config := &Config{
			Server: ServerConfig{
				Host:         "localhost",
				Port:         8080,
				ReadTimeout:  30,
				WriteTimeout: 30,
				Environment:  "development",
			},
			Database: DatabaseConfig{
				Type:             "sqlite",
				ConnectionString: "./test.db",
				MaxOpenConns:     10,
				MaxIdleConns:     5,
				MigrationsPath:   "./migrations",
			},
			Logging: LoggingConfig{
				Level:  "info",
				Format: "json",
				Output: "stdout",
			},
		}

		err := config.Validate()
		testutil.AssertNoError(t, err)
	})

	t.Run("invalid server port - zero", func(t *testing.T) {
		config := LoadConfig()
		config.Server.Port = 0

		err := config.Validate()
		testutil.AssertError(t, err)
		testutil.AssertTrue(t, err.Error() == fmt.Sprintf("%s: %d", ErrInvalidServerPort, 0), "Expected invalid server port error")
	})

	t.Run("invalid server port - negative", func(t *testing.T) {
		config := LoadConfig()
		config.Server.Port = -1

		err := config.Validate()
		testutil.AssertError(t, err)
		testutil.AssertTrue(t, err.Error() == fmt.Sprintf("%s: %d", ErrInvalidServerPort, -1), "Expected invalid server port error")
	})

	t.Run("invalid server port - too high", func(t *testing.T) {
		config := LoadConfig()
		config.Server.Port = 70000

		err := config.Validate()
		testutil.AssertError(t, err)
		testutil.AssertTrue(t, err.Error() == fmt.Sprintf("%s: %d", ErrInvalidServerPort, 70000), "Expected invalid server port error")
	})

	t.Run("invalid read timeout", func(t *testing.T) {
		config := LoadConfig()
		config.Server.ReadTimeout = 0

		err := config.Validate()
		testutil.AssertError(t, err)
		testutil.AssertTrue(t, err.Error() == fmt.Sprintf("%s: %d", ErrInvalidReadTimeout, 0), "Expected invalid read timeout error")
	})

	t.Run("invalid write timeout", func(t *testing.T) {
		config := LoadConfig()
		config.Server.WriteTimeout = -5

		err := config.Validate()
		testutil.AssertError(t, err)
		testutil.AssertTrue(t, err.Error() == fmt.Sprintf("%s: %d", ErrInvalidWriteTimeout, -5), "Expected invalid write timeout error")
	})

	t.Run("empty database type", func(t *testing.T) {
		config := LoadConfig()
		config.Database.Type = ""

		err := config.Validate()
		testutil.AssertError(t, err)
		testutil.AssertEqual(t, ErrEmptyDatabaseType, err)
	})

	t.Run("empty database connection string", func(t *testing.T) {
		config := LoadConfig()
		config.Database.ConnectionString = ""

		err := config.Validate()
		testutil.AssertError(t, err)
		testutil.AssertEqual(t, ErrEmptyConnectionString, err)
	})

	t.Run("invalid log level", func(t *testing.T) {
		config := LoadConfig()
		config.Logging.Level = "invalid"

		err := config.Validate()
		testutil.AssertError(t, err)
		testutil.AssertTrue(t, err.Error() == fmt.Sprintf("%s: %s", ErrInvalidLogLevel, "invalid"), "Expected invalid log level error")
	})

	t.Run("valid log levels", func(t *testing.T) {
		validLevels := []string{"debug", "info", "warn", "error"}

		for _, level := range validLevels {
			t.Run("log level: "+level, func(t *testing.T) {
				config := LoadConfig()
				config.Logging.Level = level

				err := config.Validate()
				testutil.AssertNoError(t, err)
			})
		}
	})
}

func TestGetEnv(t *testing.T) {
	t.Run("existing environment variable", func(t *testing.T) {
		testutil.WithEnv(t, map[string]string{
			"TEST_VAR": "test_value",
		}, func() {
			result := getEnv("TEST_VAR", "default_value")
			testutil.AssertEqual(t, "test_value", result)
		})
	})

	t.Run("non-existing environment variable", func(t *testing.T) {
		// Make sure the variable doesn't exist
		os.Unsetenv("NON_EXISTING_VAR")

		result := getEnv("NON_EXISTING_VAR", "default_value")
		testutil.AssertEqual(t, "default_value", result)
	})

	t.Run("empty environment variable", func(t *testing.T) {
		testutil.WithEnv(t, map[string]string{
			"EMPTY_VAR": "",
		}, func() {
			result := getEnv("EMPTY_VAR", "default_value")
			testutil.AssertEqual(t, "default_value", result)
		})
	})
}

func TestGetEnvAsInt(t *testing.T) {
	t.Run("valid integer environment variable", func(t *testing.T) {
		testutil.WithEnv(t, map[string]string{
			"INT_VAR": "42",
		}, func() {
			result := getEnvAsInt("INT_VAR", 10)
			testutil.AssertEqual(t, 42, result)
		})
	})

	t.Run("invalid integer environment variable", func(t *testing.T) {
		testutil.WithEnv(t, map[string]string{
			"INVALID_INT_VAR": "not_a_number",
		}, func() {
			result := getEnvAsInt("INVALID_INT_VAR", 10)
			testutil.AssertEqual(t, 10, result)
		})
	})

	t.Run("non-existing environment variable", func(t *testing.T) {
		// Make sure the variable doesn't exist
		os.Unsetenv("NON_EXISTING_INT_VAR")

		result := getEnvAsInt("NON_EXISTING_INT_VAR", 10)
		testutil.AssertEqual(t, 10, result)
	})

	t.Run("empty environment variable", func(t *testing.T) {
		testutil.WithEnv(t, map[string]string{
			"EMPTY_INT_VAR": "",
		}, func() {
			result := getEnvAsInt("EMPTY_INT_VAR", 10)
			testutil.AssertEqual(t, 10, result)
		})
	})

	t.Run("negative integer", func(t *testing.T) {
		testutil.WithEnv(t, map[string]string{
			"NEGATIVE_INT_VAR": "-5",
		}, func() {
			result := getEnvAsInt("NEGATIVE_INT_VAR", 10)
			testutil.AssertEqual(t, -5, result)
		})
	})

	t.Run("zero integer", func(t *testing.T) {
		testutil.WithEnv(t, map[string]string{
			"ZERO_INT_VAR": "0",
		}, func() {
			result := getEnvAsInt("ZERO_INT_VAR", 10)
			testutil.AssertEqual(t, 0, result)
		})
	})
}

// BenchmarkLoadConfig benchmarks the config loading performance.
func BenchmarkLoadConfig(b *testing.B) {
	// Set up some environment variables
	os.Setenv("SERVER_HOST", "localhost")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("DB_TYPE", "sqlite")
	os.Setenv("DB_CONNECTION_STRING", "./test.db")

	b.ResetTimer()

	for range b.N {
		_ = LoadConfig()
	}
}

// BenchmarkConfigValidate benchmarks the config validation performance.
func BenchmarkConfigValidate(b *testing.B) {
	config := LoadConfig()

	b.ResetTimer()

	for range b.N {
		_ = config.Validate()
	}
}
