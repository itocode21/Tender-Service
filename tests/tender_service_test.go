package tests

/*
import (
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/itocode21/Tender-Service/internal/models"
	"github.com/itocode21/Tender-Service/internal/repository"
	"github.com/itocode21/Tender-Service/internal/services"
	postgresqldb "github.com/itocode21/Tender-Service/postgresql-db"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var db *sql.DB
var organizationID int64

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
	err = db.QueryRow(`INSERT INTO organizations (name, description, type, created_at, updated_at) VALUES ('Test Organization', 'Test Organization Description', $1, NOW(), NOW()) RETURNING id`, "LLC").Scan(&organizationID)
	if err != nil {
		panic(err)
	}

	// Очистка таблиц перед каждым тестом
	_, err = db.Exec("DELETE FROM organization_responsibles")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("DELETE FROM proposals")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("DELETE FROM tenders")
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

func TestTenderService_Create(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewTenderRepository(db)
	service := services.NewTenderService(repo)

	tender := &models.Tender{
		Name:            "Test Tender",
		Description:     "Test Description",
		OrganizationID:  organizationID,
		PublicationDate: time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
	}
	// Act
	id, err := service.Create(tender)

	// Assert
	assert.NoError(t, err)
	assert.NotEqual(t, int64(0), id)

	// проверяем что статус был записан в базу данных
	retrievedTender, err := service.GetByID(id)
	assert.NoError(t, err)
	assert.Equal(t, models.TenderStatusCreated, retrievedTender.Status)

	// test with invalid name
	invalidTender := &models.Tender{
		Name:            "",
		Description:     "Test Description",
		OrganizationID:  organizationID,
		PublicationDate: time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
	}
	_, err = service.Create(invalidTender)
	assert.Error(t, err)
	assert.EqualError(t, err, "имя тендера не может быть пустым")

	// test with invalid organization
	invalidTender = &models.Tender{
		Name:            "Test Tender",
		Description:     "Test Description",
		OrganizationID:  0,
		PublicationDate: time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
	}
	_, err = service.Create(invalidTender)
	assert.Error(t, err)
	assert.EqualError(t, err, "идентификатор организации не может быть пустым")
}

func TestTenderService_GetByID(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewTenderRepository(db)
	service := services.NewTenderService(repo)

	tender := &models.Tender{
		Name:            "Test Tender",
		Description:     "Test Description",
		OrganizationID:  organizationID,
		PublicationDate: time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
		Version:         1,
		Status:          models.TenderStatusCreated,
	}
	id, err := service.Create(tender)
	assert.NoError(t, err)

	// Act
	retrievedTender, err := service.GetByID(id)
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, tender.Name, retrievedTender.Name)
	assert.Equal(t, tender.Description, retrievedTender.Description)
	assert.Equal(t, tender.OrganizationID, retrievedTender.OrganizationID)
	assert.Equal(t, tender.Version, retrievedTender.Version)
	assert.Equal(t, tender.Status, retrievedTender.Status)

	// Test non-existent tender
	_, err = service.GetByID(id + 1)
	assert.Error(t, err)
	assert.EqualError(t, err, "tender not found")
}

func TestTenderService_Update(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewTenderRepository(db)
	service := services.NewTenderService(repo)

	tender := &models.Tender{
		Name:            "Test Tender",
		Description:     "Test Description",
		OrganizationID:  organizationID,
		PublicationDate: time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
		Version:         1,
		Status:          models.TenderStatusCreated,
	}
	id, err := service.Create(tender)
	assert.NoError(t, err)
	retrievedTender, _ := service.GetByID(id)

	// Act
	retrievedTender.Name = "Updated Test Tender"
	retrievedTender.Description = "Updated Test Description"
	oldVersion := retrievedTender.Version
	err = service.Update(retrievedTender)

	// Assert
	assert.NoError(t, err)

	updatedTender, err := service.GetByID(id)
	assert.NoError(t, err)
	assert.Equal(t, retrievedTender.Name, updatedTender.Name)
	assert.Equal(t, retrievedTender.Description, updatedTender.Description)
	assert.Equal(t, oldVersion+1, updatedTender.Version)
	assert.Equal(t, retrievedTender.Status, updatedTender.Status)

	// Update non-existent tender - expect error
	retrievedTender.ID = id + 1
	err = service.Update(retrievedTender)
	assert.Error(t, err)
	assert.EqualError(t, err, "tender not found")

	//test with invalid name
	retrievedTender, _ = service.GetByID(id)
	retrievedTender.Name = ""
	err = service.Update(retrievedTender)
	assert.Error(t, err)
	assert.EqualError(t, err, "имя тендера не может быть пустым")
}

func TestTenderService_Delete(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewTenderRepository(db)
	service := services.NewTenderService(repo)

	tender := &models.Tender{
		Name:            "Test Tender",
		Description:     "Test Description",
		OrganizationID:  organizationID,
		PublicationDate: time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
		Version:         1,
		Status:          models.TenderStatusCreated,
	}
	id, err := service.Create(tender)
	assert.NoError(t, err)

	// Act
	err = service.Delete(id)
	assert.NoError(t, err)

	// Assert that the tender is deleted
	_, err = service.GetByID(id)
	assert.Error(t, err)
	assert.EqualError(t, err, "tender not found")

	// Delete non-existent tender - no error expected
	err = service.Delete(id)
	assert.NoError(t, err)
}
func TestTenderService_List(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewTenderRepository(db)
	service := services.NewTenderService(repo)

	// Create test tenders
	tenders := []*models.Tender{
		{
			Name:            "Test Tender 1",
			Description:     "Test Description 1",
			OrganizationID:  organizationID,
			PublicationDate: time.Now(),
			EndDate:         time.Now().Add(24 * time.Hour),
			Version:         1,
			Status:          models.TenderStatusCreated,
		},
		{
			Name:            "Test Tender 2",
			Description:     "Test Description 2",
			OrganizationID:  organizationID,
			PublicationDate: time.Now(),
			EndDate:         time.Now().Add(24 * time.Hour),
			Version:         1,
			Status:          models.TenderStatusCreated,
		},
	}
	for _, tender := range tenders {
		_, err := service.Create(tender)
		assert.NoError(t, err)
	}
	// Act
	retrievedTenders, err := service.List()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, retrievedTenders)
	assert.Len(t, retrievedTenders, len(tenders))

	// Test empty list
	_, err = db.Exec("DELETE FROM tenders")
	assert.NoError(t, err)
	retrievedTenders, err = service.List()
	assert.NoError(t, err)
	assert.Empty(t, retrievedTenders)

}

func TestTenderService_ListByOrganizationID(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewTenderRepository(db)
	service := services.NewTenderService(repo)
	organizationID2 := int64(2)

	err := db.QueryRow(`INSERT INTO organizations (name, description, type, created_at, updated_at) VALUES ('Test Organization 2', 'Test Organization Description 2', $1, NOW(), NOW()) RETURNING id`, "LLC").Scan(&organizationID2)
	if err != nil {
		panic(err)
	}
	// Create test tenders for different organizations
	tenders := []*models.Tender{
		{
			Name:            "Test Tender 1",
			Description:     "Test Description 1",
			OrganizationID:  organizationID,
			PublicationDate: time.Now(),
			EndDate:         time.Now().Add(24 * time.Hour),
			Version:         1,
			Status:          models.TenderStatusCreated,
		},
		{
			Name:            "Test Tender 2",
			Description:     "Test Description 2",
			OrganizationID:  organizationID,
			PublicationDate: time.Now(),
			EndDate:         time.Now().Add(24 * time.Hour),
			Version:         1,
			Status:          models.TenderStatusCreated,
		},
		{
			Name:            "Test Tender 3",
			Description:     "Test Description 3",
			OrganizationID:  organizationID2, // another organization
			PublicationDate: time.Now(),
			EndDate:         time.Now().Add(24 * time.Hour),
			Version:         1,
			Status:          models.TenderStatusCreated,
		},
	}
	for _, tender := range tenders {
		_, err := service.Create(tender)
		assert.NoError(t, err)
	}
	// Act
	retrievedTenders, err := service.ListByOrganizationID(organizationID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, retrievedTenders)
	assert.Len(t, retrievedTenders, 2)

	// Test empty list
	_, err = db.Exec("DELETE FROM tenders")
	assert.NoError(t, err)
	retrievedTenders, err = service.ListByOrganizationID(organizationID)
	assert.NoError(t, err)
	assert.Empty(t, retrievedTenders)
}
func TestTenderService_Publish(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewTenderRepository(db)
	service := services.NewTenderService(repo)

	tender := &models.Tender{
		Name:            "Test Tender",
		Description:     "Test Description",
		OrganizationID:  organizationID,
		PublicationDate: time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
		Version:         1,
		Status:          models.TenderStatusCreated,
	}
	id, err := service.Create(tender)
	assert.NoError(t, err)

	err = service.Publish(id)
	assert.NoError(t, err)

	retrievedTender, err := service.GetByID(id)
	assert.NoError(t, err)
	assert.Equal(t, models.TenderStatusPublished, retrievedTender.Status)
	// Test error if tender already published
	err = service.Publish(id)
	assert.Error(t, err)
	assert.EqualError(t, err, "тендер уже опубликован")
	// test if tender not found
	err = service.Publish(id + 1)
	assert.Error(t, err)
	assert.EqualError(t, err, "tender not found")

}

func TestTenderService_Close(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewTenderRepository(db)
	service := services.NewTenderService(repo)

	tender := &models.Tender{
		Name:            "Test Tender",
		Description:     "Test Description",
		OrganizationID:  organizationID,
		PublicationDate: time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
		Version:         1,
		Status:          models.TenderStatusCreated,
	}
	id, err := service.Create(tender)
	assert.NoError(t, err)

	err = service.Close(id)
	assert.NoError(t, err)
	retrievedTender, err := service.GetByID(id)
	assert.NoError(t, err)
	assert.Equal(t, models.TenderStatusClosed, retrievedTender.Status)

	// Test error if tender already closed
	err = service.Close(id)
	assert.Error(t, err)
	assert.EqualError(t, err, "тендер уже закрыт")
	// test if tender not found
	err = service.Close(id + 1)
	assert.Error(t, err)
	assert.EqualError(t, err, "tender not found")
	//test with cancelled tender
	tender.Status = models.TenderStatusCancelled
	repo.Update(tender)
	err = service.Close(id)
	assert.Error(t, err)
	assert.EqualError(t, err, "невозможно закрыть отмененный тендер")
}
*/
