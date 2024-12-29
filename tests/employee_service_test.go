package tests

/*
import (
	"database/sql"
	"testing"
	"time"

	"github.com/itocode21/Tender-Service/internal/models"
	"github.com/itocode21/Tender-Service/internal/repository"
	"github.com/itocode21/Tender-Service/internal/services"
	postgresqldb "github.com/itocode21/Tender-Service/postgresql-db"
	"github.com/stretchr/testify/assert"
)

var db *sql.DB

func setup() {
	var err error

	db, err = postgresqldb.InitDB()
	if err != nil {
		panic(err)
	}
	// Очистка таблицы organizations перед каждым тестом
	_, err = db.Exec("DELETE FROM organizations")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("DELETE FROM employees")
	if err != nil {
		panic(err)
	}

}

func teardown() {
	_ = db.Close()
}


func TestCreateEmployee(t *testing.T) {
	setup()
	defer teardown()

	empRepo := repository.NewEmployeeRepository(db)
	service := services.NewEmployeeService(empRepo)

	t.Run("valid employee", func(t *testing.T) {
		currentTime := time.Now()
		employee := &models.User{
			Username:  "testuser",
			FirstName: "test",
			LastName:  "user",
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		}
		id, err := service.Create(employee)

		assert.NoError(t, err)
		assert.NotEqual(t, int64(0), id)
	})

	t.Run("empty employee username", func(t *testing.T) {
		employee := &models.User{
			FirstName: "test",
			LastName:  "user",
		}
		_, err := service.Create(employee)
		assert.Error(t, err)
		assert.Equal(t, "имя пользователя не может быть пустым", err.Error())
	})

}

func TestGetEmployeeByID(t *testing.T) {
	setup()
	defer teardown()
	empRepo := repository.NewEmployeeRepository(db)
	service := services.NewEmployeeService(empRepo)

	t.Run("valid id", func(t *testing.T) {
		currentTime := time.Now()
		employee := &models.User{
			Username:  "testuser",
			FirstName: "test",
			LastName:  "user",
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		}
		id, err := service.Create(employee)
		assert.NoError(t, err)

		res, err := service.GetByID(int(id))

		assert.NoError(t, err)
		assert.Equal(t, employee.Username, res.Username)
		assert.Equal(t, employee.FirstName, res.FirstName)
		assert.Equal(t, employee.LastName, res.LastName)

	})

	t.Run("employee not found", func(t *testing.T) {
		_, err := service.GetByID(123456789)

		assert.Error(t, err)
		assert.Equal(t, "employee not found", err.Error())
	})
}

func TestUpdateEmployee(t *testing.T) {
	setup()
	defer teardown()

	empRepo := repository.NewEmployeeRepository(db)
	service := services.NewEmployeeService(empRepo)

	t.Run("valid employee", func(t *testing.T) {
		currentTime := time.Now()
		employee := &models.User{
			Username:  "testuser",
			FirstName: "test",
			LastName:  "user",
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		}
		id, err := service.Create(employee)
		assert.NoError(t, err)
		employee.Id = int(id)
		employee.FirstName = "new test"
		err = service.Update(employee)

		assert.NoError(t, err)
	})

}

func TestDeleteEmployee(t *testing.T) {
	setup()
	defer teardown()
	empRepo := repository.NewEmployeeRepository(db)
	service := services.NewEmployeeService(empRepo)

	t.Run("valid id", func(t *testing.T) {
		currentTime := time.Now()
		employee := &models.User{
			Username:  "testuser",
			FirstName: "test",
			LastName:  "user",
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		}
		id, err := service.Create(employee)
		assert.NoError(t, err)

		err = service.Delete(int(id))
		assert.NoError(t, err)

	})
}

func TestListEmployees(t *testing.T) {
	setup()
	defer teardown()
	empRepo := repository.NewEmployeeRepository(db)
	service := services.NewEmployeeService(empRepo)

	t.Run("valid request", func(t *testing.T) {
		currentTime := time.Now()
		employee1 := &models.User{
			Username:  "testuser1",
			FirstName: "test1",
			LastName:  "user1",
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		}
		employee2 := &models.User{
			Username:  "testuser2",
			FirstName: "test2",
			LastName:  "user2",
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		}
		_, err := service.Create(employee1)
		assert.NoError(t, err)
		_, err = service.Create(employee2)
		assert.NoError(t, err)

		employees, err := service.List()

		assert.NoError(t, err)
		assert.NotNil(t, employees)
		assert.Equal(t, 2, len(employees))
	})

}
*/
