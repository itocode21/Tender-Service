package tests

import (
	"database/sql"
	"log"
	"testing"

	postgresqldb "github.com/itocode21/Tender-Service/postgresql-db"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func TestInitDB(t *testing.T) {
	connStr := "postgresql://ito21:1899@localhost:5432/TENDER" // Убедитесь, что строка подключения правильная
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatalf("Ошибка подключения к базе: %v", err)
	}
	defer db.Close()

	createdDB, err := postgresqldb.InitDB()
	if err != nil {
		t.Fatalf("Ошибка вызова postgresqldb.InitDB: %v", err)
	}

	if createdDB == nil {
		t.Fatal("Ожидаемое подключения к базе пустое")
	}

	err = createdDB.Ping()
	if err != nil {
		t.Fatalf("Ошибка пинга базы: %v", err)
	}

	var exists bool

	// Проверка существования таблицы employees
	err = createdDB.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'employees')").Scan(&exists)
	if err != nil {
		t.Fatalf("Не удалось проверить, существует ли таблица employees: %v", err)
	}
	if !exists {
		t.Fatal("Таблицы employees не существует")
	}

	// Проверка существования таблицы organizations
	err = createdDB.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'organizations')").Scan(&exists)
	if err != nil {
		t.Fatalf("Не удалось проверить, существует ли таблица organizations: %v", err)
	}
	if !exists {
		t.Fatal("Таблицы organizations не существует")
	}

	// Проверка существования таблицы organization_responsibles
	err = createdDB.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'organization_responsibles')").Scan(&exists)
	if err != nil {
		t.Fatalf("Не удалось проверить, существует ли таблица organization_responsibles: %v", err)
	}
	if !exists {
		t.Fatal("Таблицы organization_responsibles не существует")
	}
	log.Println("Тест пройден: соединение с базой установлено и стабильно, таблицы employees, organizations, organization_responsibles успешно созданы.")
}
