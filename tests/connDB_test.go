package tests

import (
	"log"
	"os"
	"testing"

	postgresqldb "github.com/itocode21/Tender-Service/postgresql-db"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}
}

func TestInitDB(t *testing.T) {
	// переменные окружения для строки подключения
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		t.Fatal("Не указаны необходимые переменные окружения для подключения к базе данных")
	}

	// Инициализация базы данных
	createdDB, err := postgresqldb.InitDB()
	if err != nil {
		t.Fatalf("Ошибка вызова postgresqldb.InitDB: %v", err)
	}
	defer createdDB.Close() // закрываем соединение

	if createdDB == nil {
		t.Fatal("Ожидаемое подключения к базе пустое")
	}

	err = createdDB.Ping()
	if err != nil {
		t.Fatalf("Ошибка пинга базы: %v", err)
	}

	// Проверка существования таблиц
	tables := []string{"employees", "organizations", "organization_responsibles"}
	for _, table := range tables {
		t.Run("CheckTableExists_"+table, func(t *testing.T) {
			var exists bool
			err = createdDB.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)", table).Scan(&exists)
			if err != nil {
				t.Fatalf("Не удалось проверить, существует ли таблица %s: %v", table, err)
			}
			if !exists {
				t.Fatalf("Таблицы %s не существует", table)
			}
		})
	}

	log.Println("Тест пройден: соединение с базой установлено и стабильно, таблицы успешно созданы.")
}
