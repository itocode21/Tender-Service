package tests

/*
import (
	"database/sql"
	"testing"

	"github.com/itocode21/Tender-Service/internal/models"
	"github.com/itocode21/Tender-Service/internal/repository"
	postgresqldb "github.com/itocode21/Tender-Service/postgresql-db"
	_ "github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

var db1 *sql.DB

func setup() {
	var err error
	db, err = postgresqldb.InitDB()
	if err != nil {
		panic(err)
	}

	// Очистка таблицы employees перед каждым тестом
	_, err = db.Exec("DELETE FROM employees")
	if err != nil {
		panic(err)
	}
}
func teardown() {
	_ = db.Close()
}

func TestEmployeeRepository_Create(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewEmployeeRepository(db)

	employee := &models.User{
		Username:  "AvitoTest",
		FirstName: "AvitoName",
		LastName:  "AvitoLastName",
	}

	id, err := repo.Create(employee)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM employees WHERE username=$1", employee.Username).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestEmployeeRepository_GetByID(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewEmployeeRepository(db)

	employee := &models.User{
		Username:  "TestAvitoUsername",
		FirstName: "AvitoName",
		LastName:  "AvitolastName",
	}
	id, err := repo.Create(employee)
	assert.NoError(t, err)

	fetchedEmployee, err := repo.GetByID(id)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedEmployee)
	assert.Equal(t, employee.Username, fetchedEmployee.Username)
}

func TestEmployeeRepository_GetByUsername(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewEmployeeRepository(db)

	employee := &models.User{
		Username:  "TestAvitoUsername",
		FirstName: "AvitoName",
		LastName:  "AvitolastName",
	}

	_, err := repo.Create(employee)
	assert.NoError(t, err)

	fetchedEmployee, err := repo.GetByUsername(employee.Username)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedEmployee)
	assert.Equal(t, employee.Username, fetchedEmployee.Username)
}
*/
