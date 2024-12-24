package postgresqldb

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
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
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	// Проверка подключения, обработка sql.ErrNoRows
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки подключения к базе данных: %w", err)
	}

	// Создаем транзакцию
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer tx.Rollback() // Автоматический откат если ошибка

	if err := createTables(tx); err != nil {
		return nil, fmt.Errorf("ошибка при создании таблиц: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("ошибка коммита транзакции: %w", err)
	}

	return db, nil // Возвращаем подключение
}

func createTables(tx *sql.Tx) error {
	// Создание типа organization_type
	_, err := tx.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'organization_type') THEN
				CREATE TYPE organization_type AS ENUM ('IE', 'LLC', 'JSC');
			END IF;
		END $$;
	`)
	if err != nil {
		return fmt.Errorf("ошибка при создании типа organization_type: %w", err)
	}

	// Создание таблицы employees
	_, err = tx.Exec(`
        CREATE TABLE IF NOT EXISTS employees (
            id SERIAL PRIMARY KEY,
            username VARCHAR(50) UNIQUE NOT NULL,
            first_name VARCHAR(50),
            last_name VARCHAR(50),
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		return fmt.Errorf("ошибка при создании таблицы employees: %w", err)
	}

	// Создание таблицы organizations
	_, err = tx.Exec(`
        CREATE TABLE IF NOT EXISTS organizations (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            description TEXT,
            type organization_type,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		return fmt.Errorf("ошибка при создании таблицы organizations: %w", err)
	}

	// Создание таблицы organization_responsibles
	_, err = tx.Exec(`
        CREATE TABLE IF NOT EXISTS organization_responsibles (
            id SERIAL PRIMARY KEY,
            organization_id INT REFERENCES organizations(id) ON DELETE CASCADE,
            user_id INT REFERENCES employees(id) ON DELETE CASCADE
        )
    `)
	if err != nil {
		return fmt.Errorf("ошибка при создании таблицы organization_responsibles: %w", err)
	}

	return nil
}
