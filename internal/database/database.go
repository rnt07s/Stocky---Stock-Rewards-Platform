package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// Connect establishes a database connection
func Connect(dsn string, log *logrus.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Info("Database connection established")
	return db, nil
}

// RunMigrations executes SQL migration files
func RunMigrations(db *sql.DB, migrationsPath string, log *logrus.Logger) error {
	files, err := filepath.Glob(filepath.Join(migrationsPath, "*.sql"))
	if err != nil {
		return fmt.Errorf("failed to find migration files: %w", err)
	}

	if len(files) == 0 {
		log.Warn("No migration files found")
		return nil
	}

	for _, file := range files {
		log.Infof("Running migration: %s", filepath.Base(file))
		
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}
	}

	log.Info("All migrations completed successfully")
	return nil
}
