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
	"github.com/joho/godotenv"

	"log"

	"github.com/stretchr/testify/assert"
)

var db *sql.DB

func setup() {
	err1 := godotenv.Load()
	if err1 != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err1)
	}
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

func TestCreateOrganization(t *testing.T) {
	setup()
	defer teardown()

	orgRepo := repository.NewOrganizationRepository(db)
	service := services.NewOrganizationService(orgRepo)
	t.Run("empty organization name", func(t *testing.T) {
		org := &models.Organization{Description: "Test Description"}
		_, err := service.Create(org)
		assert.Error(t, err)
		assert.Equal(t, "имя организации не может быть пустым", err.Error())
	})
	t.Run("valid organization", func(t *testing.T) {
		org := &models.Organization{Name: "Test Org", Description: "Test Description", Type: models.OrganizationTypeIE}

		id, err := service.Create(org)
		assert.NoError(t, err)
		assert.NotEqual(t, int64(0), id)

	})

}

func TestGetOrganizationByID(t *testing.T) {
	setup()
	defer teardown()

	orgRepo := repository.NewOrganizationRepository(db)
	service := services.NewOrganizationService(orgRepo)
	t.Run("valid id", func(t *testing.T) {
		org := &models.Organization{Name: "Test Org", Description: "Test Description", Type: models.OrganizationTypeIE}
		id, err := service.Create(org)
		assert.NoError(t, err)
		res, err := service.GetByID(int64(id))

		assert.NoError(t, err)
		assert.Equal(t, org.Name, res.Name)
		assert.Equal(t, org.Description, res.Description)
		assert.Equal(t, org.Type, res.Type)
	})

	t.Run("organization not found", func(t *testing.T) {
		_, err := service.GetByID(int64(123456789))

		assert.Error(t, err)
		assert.Equal(t, "organization not found", err.Error())

	})

}
func TestUpdateOrganization(t *testing.T) {
	setup()
	defer teardown()

	orgRepo := repository.NewOrganizationRepository(db)
	service := services.NewOrganizationService(orgRepo)
	t.Run("valid organization", func(t *testing.T) {
		org := &models.Organization{Name: "Test Org", Description: "Test Description", Type: models.OrganizationTypeIE}
		id, err := service.Create(org)
		assert.NoError(t, err)
		org.Id = int64(id)
		org.Name = "new Test Org"
		err = service.Update(org)

		assert.NoError(t, err)
	})

}

func TestDeleteOrganization(t *testing.T) {
	setup()
	defer teardown()
	orgRepo := repository.NewOrganizationRepository(db)
	service := services.NewOrganizationService(orgRepo)
	t.Run("valid id", func(t *testing.T) {
		org := &models.Organization{Name: "Test Org", Description: "Test Description", Type: models.OrganizationTypeIE}
		id, err := service.Create(org)
		assert.NoError(t, err)

		err = service.Delete(id)

		assert.NoError(t, err)
	})
}
func TestListOrganizations(t *testing.T) {
	setup()
	defer teardown()

	orgRepo := repository.NewOrganizationRepository(db)
	service := services.NewOrganizationService(orgRepo)

	t.Run("valid request", func(t *testing.T) {
		org1 := &models.Organization{Name: "Test Org 1", Description: "Test Description 1", Type: models.OrganizationTypeIE}
		org2 := &models.Organization{Name: "Test Org 2", Description: "Test Description 2", Type: models.OrganizationTypeJSC}
		_, err := service.Create(org1)
		assert.NoError(t, err)
		_, err = service.Create(org2)
		assert.NoError(t, err)

		orgs, err := service.List()

		assert.NoError(t, err)
		assert.NotNil(t, orgs)
		assert.Equal(t, 2, len(orgs))

	})

}

func TestAddResponsible(t *testing.T) {
	setup()
	defer teardown()

	orgRepo := repository.NewOrganizationRepository(db)
	service := services.NewOrganizationService(orgRepo)

	t.Run("valid request", func(t *testing.T) {
		org := &models.Organization{Name: "Test Org", Type: models.OrganizationTypeIE}
		id, err := service.Create(org)
		assert.NoError(t, err)

		currentTime := time.Now()
		user := &models.User{
			Username:  "testuser",
			FirstName: "test",
			LastName:  "user",
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		}

		userRepo := repository.NewEmployeeRepository(db)
		userId, err := userRepo.Create(user)
		assert.NoError(t, err)

		orgResp := &models.OrganizationResponsible{
			OrganizationID: &models.Organization{Id: int64(id)},
			UserID:         &models.User{Id: int(userId)},
		}
		respId, err := service.AddResponsible(orgResp)

		assert.NoError(t, err)
		assert.NotEqual(t, int64(0), respId)
	})
	t.Run("organization or user is nil", func(t *testing.T) {
		orgResp := &models.OrganizationResponsible{}
		_, err := service.AddResponsible(orgResp)

		assert.Error(t, err)
		assert.Equal(t, "организация и пользователь не могут быть nil", err.Error())
	})
}

func TestRemoveResponsible(t *testing.T) {
	setup()
	defer teardown()

	orgRepo := repository.NewOrganizationRepository(db)
	service := services.NewOrganizationService(orgRepo)

	t.Run("valid request", func(t *testing.T) {
		org := &models.Organization{Name: "Test Org", Type: models.OrganizationTypeIE}
		id, err := service.Create(org)
		assert.NoError(t, err)

		currentTime := time.Now()
		user := &models.User{
			Username:  "testuser",
			FirstName: "test",
			LastName:  "user",
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		}
		userRepo := repository.NewEmployeeRepository(db)
		userId, err := userRepo.Create(user)
		assert.NoError(t, err)
		orgResp := &models.OrganizationResponsible{
			OrganizationID: &models.Organization{Id: int64(id)},
			UserID:         &models.User{Id: int(userId)},
		}
		_, err = service.AddResponsible(orgResp)
		assert.NoError(t, err)
		err = service.RemoveResponsible(orgResp)
		assert.NoError(t, err)
	})
}
*/
