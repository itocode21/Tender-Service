package postgresqldb

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"

	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла")
	}
}

func InitDB() (*sql.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	fmt.Printf("DB_HOST: %s, DB_PORT: %s, DB_USER: %s, DB_NAME: %s\n", dbHost, dbPort, dbUser, dbName)
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	// Проверка подключения
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging the database: %w", err)
	}

	if err := createTable(db); err != nil {
		return nil, err
	}

	return db, nil
}

func createTable(db *sql.DB) error {
	// Создание типа organization_type, если он не существует
	createOrganizationType := `
    DO $$ BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'organization_type') THEN
            CREATE TYPE organization_type AS ENUM (
                'IE',
                'LLC',
                'JSC'
            );
        END IF;
    END $$;`

	_, err := db.Exec(createOrganizationType)
	if err != nil {
		return fmt.Errorf("ошибка при создании типа organization_type: %w", err)
	}

	// Создание таблицы employees
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

	// Создание таблицы organizations
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

	// Создание таблицы organization_responsibles
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
