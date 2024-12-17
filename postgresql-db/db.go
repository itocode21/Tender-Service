package postgresqldb

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4"

	"context"
)

func InitDB() (*sql.DB, error) {
	connStr := "postgres://ito21:1899@localhost:5432/TENDER"

	ctx := context.Background()

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка при проверке подключенияк дб: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	if err := createTable(db); err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ошибки при проверке подключения: %w", err)
	}

	return db, nil
}

func createTable(db *sql.DB) error {

	// Создание типа organization_type
	createOrganizationType := `
	CREATE TYPE organization_type AS ENUM (
		'IE',
		'LLC',
		'JSC'
	);`

	_, err := db.Exec(createOrganizationType)
	if err != nil {
		return fmt.Errorf("ошибка при создании типа organization_type: %w", err)
	}

	createEmployeeTable := `
	CREATE TABLE IF NOT EXISTS employees (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		first_name VARCHAR(50),
		last_name VARCHAR(50),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`

	_, err = db.Exec(createEmployeeTable)
	if err != nil {
		return fmt.Errorf("ошибка при создании таблицы employees: %w", err)
	}

	createOrganizationTable := `
	CREATE TABLE IF NOT EXISTS organizations (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT,
		type organization_type,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`

	_, err = db.Exec(createOrganizationTable)
	if err != nil {
		return fmt.Errorf("ошибка при создании таблицы organizations: %w", err)
	}

	createOrganizationResponsibleTable := `
	CREATE TABLE IF NOT EXISTS organization_responsibles (
		id SERIAL PRIMARY KEY,
		organization_id INT REFERENCES organizations(id) ON DELETE CASCADE,
		user_id INT REFERENCES employees(id) ON DELETE CASCADE);`

	_, err = db.Exec(createOrganizationResponsibleTable)
	if err != nil {
		return fmt.Errorf("ошибка при создании таблицы organization_responsibles: %w", err)
	}

	return nil
}
