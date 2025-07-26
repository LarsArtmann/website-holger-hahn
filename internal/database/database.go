// Package database provides database connection, initialization, and repository implementations
// using sqlc generated code for type-safe SQL operations with SQLite backend.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Config holds database configuration options.
type Config struct {
	// DatabasePath is the path to the SQLite database file
	DatabasePath string
	
	// MaxOpenConns is the maximum number of open connections
	MaxOpenConns int
	
	// MaxIdleConns is the maximum number of idle connections
	MaxIdleConns int
	
	// ConnMaxLifetime is the maximum amount of time a connection may be reused
	ConnMaxLifetime time.Duration
	
	// EnableWAL enables Write-Ahead Logging for better concurrency
	EnableWAL bool
	
	// BusyTimeout sets how long to wait for locked database
	BusyTimeout time.Duration
}

// DefaultConfig returns a default database configuration.
func DefaultConfig() *Config {
	return &Config{
		DatabasePath:    "./data/holger-hahn.db",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		EnableWAL:       true,
		BusyTimeout:     30 * time.Second,
	}
}

// DatabaseManager manages the database connection and provides access to repositories.
type DatabaseManager struct {
	db      *sql.DB
	queries *Queries
	config  *Config
}

// NewDatabaseManager creates a new database manager with the given configuration.
func NewDatabaseManager(config *Config) (*DatabaseManager, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Ensure database directory exists
	if err := ensureDirectoryExists(config.DatabasePath); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Build connection string with pragma settings
	connStr := buildConnectionString(config)

	// Open database connection
	db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create sqlc queries instance
	queries := New(db)

	return &DatabaseManager{
		db:      db,
		queries: queries,
		config:  config,
	}, nil
}

// buildConnectionString constructs SQLite connection string with pragma settings.
func buildConnectionString(config *Config) string {
	params := []string{
		"_busy_timeout=" + fmt.Sprintf("%d", int(config.BusyTimeout.Milliseconds())),
		"_journal_mode=WAL",
		"_synchronous=NORMAL",
		"_cache_size=1000",
		"_foreign_keys=true",
		"_temp_store=memory",
	}

	if !config.EnableWAL {
		// Override WAL mode if disabled
		params[1] = "_journal_mode=DELETE"
	}

	connStr := config.DatabasePath
	for i, param := range params {
		if i == 0 {
			connStr += "?" + param
		} else {
			connStr += "&" + param
		}
	}

	return connStr
}

// ensureDirectoryExists creates the directory for the database file if it doesn't exist.
func ensureDirectoryExists(dbPath string) error {
	dir := filepath.Dir(dbPath)
	if dir == "." {
		return nil // Current directory
	}

	return os.MkdirAll(dir, 0755)
}

// DB returns the underlying database connection.
func (dm *DatabaseManager) DB() *sql.DB {
	return dm.db
}

// Queries returns the sqlc generated queries instance.
func (dm *DatabaseManager) Queries() *Queries {
	return dm.queries
}

// WithTx executes a function within a database transaction.
func (dm *DatabaseManager) WithTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := dm.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	queries := dm.queries.WithTx(tx)
	
	if err := fn(queries); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error: %w, rollback error: %v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// HealthCheck performs a health check on the database connection.
func (dm *DatabaseManager) HealthCheck(ctx context.Context) error {
	return dm.db.PingContext(ctx)
}

// Close closes the database connection.
func (dm *DatabaseManager) Close() error {
	if dm.db != nil {
		return dm.db.Close()
	}
	return nil
}

// Migrate applies database migrations (placeholder for now).
func (dm *DatabaseManager) Migrate(ctx context.Context) error {
	// Read and execute schema file
	schemaPath := "./internal/database/schema/001_initial.sql"
	schemaSQL, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	if _, err := dm.db.ExecContext(ctx, string(schemaSQL)); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}

// GetStats returns database connection statistics.
func (dm *DatabaseManager) GetStats() sql.DBStats {
	return dm.db.Stats()
}