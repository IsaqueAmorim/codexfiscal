package config

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

type Migration struct {
	Database *sql.DB
}

func NewMigration(db *sql.DB) *Migration {
	return &Migration{
		Database: db,
	}
}

func (m *Migration) Run() error {
	if m.Database == nil {
		return fmt.Errorf("database connection is nil")
	}

	// Check if the NCM table exists
	var exists bool
	err := m.Database.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'ncm')").Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if NCM table exists: %w", err)
	}

	if exists {
		fmt.Println("NCM table already exists, skipping migration.")
		return nil
	}

	// Execute the migration script
	fmt.Println("Creating NCM table...")
	err = m.createNCMTable()
	if err != nil {
		return fmt.Errorf("error creating NCM table: %w", err)
	}

	fmt.Println("NCM table created successfully!")
	return nil
}

func (m *Migration) createNCMTable() error {
	// Read the migration SQL file
	migrationPath := filepath.Join("scripts", "migrations", "001_create_ncm_table.sql")
	sqlBytes, err := os.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("error reading migration file: %w", err)
	}

	sqlScript := string(sqlBytes)

	// Execute the SQL script
	_, err = m.Database.Exec(sqlScript)
	if err != nil {
		return fmt.Errorf("error executing migration script: %w", err)
	}

	return nil
}
