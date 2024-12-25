package tests

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/itocode21/Tender-Service/internal/models"
	"github.com/itocode21/Tender-Service/internal/repository"
	"github.com/itocode21/Tender-Service/internal/services"
	postgresqldb "github.com/itocode21/Tender-Service/postgresql-db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

type MockOrganizationRepository struct {
	mock.Mock
}

func (m *MockOrganizationRepository) Create(org *models.Organization) (int64, error) {
	args := m.Called(org)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockOrganizationRepository) GetByID(id int64) (*models.Organization, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Organization), args.Error(1)
}

func (m *MockOrganizationRepository) Update(org *models.Organization) error {
	args := m.Called(org)
	return args.Error(0)
}

func (m *MockOrganizationRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockOrganizationRepository) List() ([]*models.Organization, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Organization), args.Error(1)
}
func (m *MockOrganizationRepository) AddResponsible(orgResp *models.OrganizationResponsible) (int64, error) {
	args := m.Called(orgResp)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockOrganizationRepository) RemoveResponsible(orgResp *models.OrganizationResponsible) error {
	args := m.Called(orgResp)
	return args.Error(0)
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
	t.Run("repo return error", func(t *testing.T) {
		org := &models.Organization{Name: "Test Org", Description: "Test Description", Type: models.OrganizationTypeIE}
		mockRepo := new(MockOrganizationRepository)
		mockRepo.On("Create", org).Return(int64(0), errors.New("db error"))

		service := services.NewOrganizationService(mockRepo)

		_, err := service.Create(org)
		assert.Error(t, err)
		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())

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
		res, err := service.GetByID(id)

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
	t.Run("repo return error", func(t *testing.T) {
		mockRepo := new(MockOrganizationRepository)
		mockRepo.On("GetByID", int64(1)).Return(nil, errors.New("db error"))
		service := services.NewOrganizationService(mockRepo)

		_, err := service.GetByID(int64(1))

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
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
		org.Id = id
		org.Name = "new Test Org"
		err = service.Update(org)

		assert.NoError(t, err)
	})

	t.Run("repo return error", func(t *testing.T) {
		org := &models.Organization{Id: int64(1), Name: "Test Org", Description: "Test Description", Type: models.OrganizationTypeIE}
		mockRepo := new(MockOrganizationRepository)
		mockRepo.On("Update", org).Return(errors.New("db error"))

		service := services.NewOrganizationService(mockRepo)
		err := service.Update(org)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())

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

	t.Run("repo return error", func(t *testing.T) {
		mockRepo := new(MockOrganizationRepository)
		mockRepo.On("Delete", int64(1)).Return(errors.New("db error"))
		service := services.NewOrganizationService(mockRepo)
		err := service.Delete(int64(1))
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())

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

	t.Run("repo return error", func(t *testing.T) {
		mockRepo := new(MockOrganizationRepository)
		mockRepo.On("List").Return(nil, errors.New("db error"))
		service := services.NewOrganizationService(mockRepo)

		_, err := service.List()

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}

// TestAddResponsible тестирует добавление ответственного лица
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
			OrganizationID: &models.Organization{Id: id},
			UserID:         &models.User{Id: userId},
		}
		respId, err := service.AddResponsible(orgResp)

		assert.NoError(t, err)
		assert.NotEqual(t, int64(0), respId)
	})
	t.Run("repo return error", func(t *testing.T) {
		org := &models.Organization{Id: 1}
		user := &models.User{Id: 2}
		orgResp := &models.OrganizationResponsible{
			OrganizationID: org,
			UserID:         user,
		}
		mockRepo := new(MockOrganizationRepository)
		mockRepo.On("AddResponsible", mock.Anything).Return(int64(0), errors.New("db error"))
		service := services.NewOrganizationService(mockRepo)
		_, err := service.AddResponsible(orgResp)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})
	t.Run("organization or user is nil", func(t *testing.T) {
		orgResp := &models.OrganizationResponsible{}
		_, err := service.AddResponsible(orgResp)

		assert.Error(t, err)
		assert.Equal(t, "организация и пользователь не могут быть nil", err.Error())
	})
}

// TestRemoveResponsible тестирует удаление ответственного лица
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
			OrganizationID: &models.Organization{Id: id},
			UserID:         &models.User{Id: userId},
		}
		_, err = service.AddResponsible(orgResp)
		assert.NoError(t, err)

		err = service.RemoveResponsible(orgResp)

		assert.NoError(t, err)
	})

	t.Run("repo return error", func(t *testing.T) {
		org := &models.Organization{Id: 1}
		user := &models.User{Id: 2}
		orgResp := &models.OrganizationResponsible{
			OrganizationID: org,
			UserID:         user,
		}
		mockRepo := new(MockOrganizationRepository)
		mockRepo.On("RemoveResponsible", orgResp).Return(errors.New("db error"))
		service := services.NewOrganizationService(mockRepo)
		err := service.RemoveResponsible(orgResp)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}
