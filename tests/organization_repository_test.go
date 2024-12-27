package tests

/*
import (
	"database/sql"
	"log"
	"testing"

	"github.com/itocode21/Tender-Service/internal/models"
	"github.com/itocode21/Tender-Service/internal/repository"
	postgresqldb "github.com/itocode21/Tender-Service/postgresql-db"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
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

	// Очистка таблиц перед каждым тестом
	_, err = db.Exec("DELETE FROM organization_responsibles")
	if err != nil {
		panic(err)
	}
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

func TestOrganizationRepository_Create(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewOrganizationRepository(db)

	org := &models.Organization{
		Name:        "Test Organization",
		Description: "Test Description",
		Type:        models.OrganizationTypeLLC,
	}

	// Act
	id, err := repo.Create(org)

	// Assert
	assert.NoError(t, err)
	assert.NotEqual(t, int64(0), id)

	retrievedOrg, err := repo.GetByID(id)
	assert.NoError(t, err)
	assert.Equal(t, org.Name, retrievedOrg.Name)
	assert.Equal(t, org.Description, retrievedOrg.Description)
	assert.Equal(t, org.Type, retrievedOrg.Type)

}
func TestOrganizationRepository_GetByID(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewOrganizationRepository(db)

	org := &models.Organization{
		Name:        "Test Organization",
		Description: "Test Description",
		Type:        models.OrganizationTypeLLC,
	}

	id, err := repo.Create(org)
	assert.NoError(t, err)
	// Act
	retrievedOrg, err := repo.GetByID(id)
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, org.Name, retrievedOrg.Name)
	assert.Equal(t, org.Description, retrievedOrg.Description)
	assert.Equal(t, org.Type, retrievedOrg.Type)

	// Test non-existent tender
	_, err = repo.GetByID(id + 1)
	assert.Error(t, err)
	assert.EqualError(t, err, "organization not found")
}
func TestOrganizationRepository_Update(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewOrganizationRepository(db)

	org := &models.Organization{
		Name:        "Test Organization",
		Description: "Test Description",
		Type:        models.OrganizationTypeLLC,
	}
	id, err := repo.Create(org)
	assert.NoError(t, err)

	// Act
	org.Name = "Updated Test Organization"
	org.Description = "Updated Test Description"
	org.Type = models.OrganizationTypeJSC
	err = repo.Update(org)

	// Assert
	assert.NoError(t, err)

	updatedOrg, err := repo.GetByID(id)
	assert.NoError(t, err)
	assert.Equal(t, org.Name, updatedOrg.Name)
	assert.Equal(t, org.Description, updatedOrg.Description)
	assert.Equal(t, org.Type, updatedOrg.Type)

	// Update non-existent tender - no error expected
	org.Id = id + 1
	err = repo.Update(org)
	assert.NoError(t, err)
}

func TestOrganizationRepository_Delete(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewOrganizationRepository(db)

	org := &models.Organization{
		Name:        "Test Organization",
		Description: "Test Description",
		Type:        models.OrganizationTypeLLC,
	}
	id, err := repo.Create(org)
	assert.NoError(t, err)

	// Act
	err = repo.Delete(id)
	assert.NoError(t, err)

	// Assert that the tender is deleted
	_, err = repo.GetByID(id)
	assert.Error(t, err)
	assert.EqualError(t, err, "organization not found")

	// Delete non-existent tender - no error expected
	err = repo.Delete(id)
	assert.NoError(t, err)
}

func TestOrganizationRepository_List(t *testing.T) {
	setup()
	defer teardown()
	repo := repository.NewOrganizationRepository(db)
	// Create test organizations
	organizations := []*models.Organization{
		{
			Name:        "Test Organization 1",
			Description: "Test Description 1",
			Type:        models.OrganizationTypeLLC,
		},
		{
			Name:        "Test Organization 2",
			Description: "Test Description 2",
			Type:        models.OrganizationTypeJSC,
		},
	}
	for _, org := range organizations {
		_, err := repo.Create(org)
		assert.NoError(t, err)
	}

	retrievedOrganizations, err := repo.List()

	assert.NoError(t, err)
	assert.NotNil(t, retrievedOrganizations)
	assert.Len(t, retrievedOrganizations, len(organizations))

	_, err = db.Exec("DELETE FROM organizations")
	assert.NoError(t, err)
	retrievedOrganizations, err = repo.List()
	assert.NoError(t, err)
	assert.Empty(t, retrievedOrganizations)

}
*/
