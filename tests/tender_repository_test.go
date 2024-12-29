package tests

/*
import (
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/itocode21/Tender-Service/internal/repository"
	postgresqldb "github.com/itocode21/Tender-Service/postgresql-db"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/itocode21/Tender-Service/internal/models"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var db *sql.DB
var organizationID int64
var organizationID2 int64

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
	// Создание организации и получение ее ID
	err = db.QueryRow(`INSERT INTO organizations (name, description, type, created_at, updated_at) VALUES ('Test Organization', 'Test Organization Description', $1, NOW(), NOW()) RETURNING id`, "LLC").Scan(&organizationID)
	if err != nil {
		panic(err)
	}
	err = db.QueryRow(`INSERT INTO organizations (name, description, type, created_at, updated_at) VALUES ('Test Organization 2', 'Test Organization Description 2', $1, NOW(), NOW()) RETURNING id`, "LLC").Scan(&organizationID2)
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

func TestTenderRepository_Create(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewTenderRepository(db)

	tender := &models.Tender{
		Name:            "Test Tender",
		Description:     "Test Description",
		OrganizationID:  organizationID,
		PublicationDate: time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
		Version:         1,
		Status:          models.TenderStatusCreated,
	}
	id, err := repo.Create(tender)

	assert.NoError(t, err)
	assert.NotEqual(t, int64(0), id)

	retrievedTender, err := repo.GetByID(id)
	assert.NoError(t, err)
	assert.Equal(t, models.TenderStatusCreated, retrievedTender.Status)
}

func TestTenderRepository_GetByID(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewTenderRepository(db)

	tender := &models.Tender{
		Name:            "Test Tender",
		Description:     "Test Description",
		OrganizationID:  organizationID,
		PublicationDate: time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
		Version:         1,
		Status:          models.TenderStatusCreated,
	}
	id, err := repo.Create(tender)
	assert.NoError(t, err)

	retrievedTender, err := repo.GetByID(id)
	assert.NoError(t, err)
	assert.Equal(t, tender.Name, retrievedTender.Name)
	assert.Equal(t, tender.Description, retrievedTender.Description)
	assert.Equal(t, tender.OrganizationID, retrievedTender.OrganizationID)
	assert.Equal(t, tender.Version, retrievedTender.Version)
	assert.Equal(t, tender.Status, retrievedTender.Status)

	_, err = repo.GetByID(id + 1)
	assert.Error(t, err)
	assert.EqualError(t, err, "tender not found")
}
func TestTenderRepository_Update(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewTenderRepository(db)

	tender := &models.Tender{
		Name:            "Test Tender",
		Description:     "Test Description",
		OrganizationID:  organizationID,
		PublicationDate: time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
		Version:         1,
		Status:          models.TenderStatusCreated,
	}
	id, err := repo.Create(tender)
	assert.NoError(t, err)

	tender.Name = "Updated Test Tender"
	tender.Description = "Updated Test Description"
	tender.Version = 2
	err = repo.Update(tender)

	assert.NoError(t, err)

	updatedTender, err := repo.GetByID(id)
	assert.NoError(t, err)
	assert.Equal(t, tender.Name, updatedTender.Name)
	assert.Equal(t, tender.Description, updatedTender.Description)
	assert.Equal(t, tender.Version, updatedTender.Version)
	assert.Equal(t, tender.Status, updatedTender.Status)

	tender.ID = id + 1
	err = repo.Update(tender)
	assert.NoError(t, err)
}

func TestTenderRepository_Delete(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewTenderRepository(db)

	tender := &models.Tender{
		Name:            "Test Tender",
		Description:     "Test Description",
		OrganizationID:  organizationID,
		PublicationDate: time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
		Version:         1,
		Status:          models.TenderStatusCreated,
	}
	id, err := repo.Create(tender)
	assert.NoError(t, err)

	err = repo.Delete(id)
	assert.NoError(t, err)

	_, err = repo.GetByID(id)
	assert.Error(t, err)
	assert.EqualError(t, err, "tender not found")

	err = repo.Delete(id)
	assert.NoError(t, err)
}
func TestTenderRepository_List(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewTenderRepository(db)

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
		_, err := repo.Create(tender)
		assert.NoError(t, err)
	}

	retrievedTenders, err := repo.List()

	assert.NoError(t, err)
	assert.NotNil(t, retrievedTenders)
	assert.Len(t, retrievedTenders, len(tenders))

	_, err = db.Exec("DELETE FROM tenders")
	assert.NoError(t, err)
	retrievedTenders, err = repo.List()
	assert.NoError(t, err)
	assert.Empty(t, retrievedTenders)

}

func TestTenderRepository_ListByOrganizationID(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewTenderRepository(db)

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
			OrganizationID:  organizationID2,
			PublicationDate: time.Now(),
			EndDate:         time.Now().Add(24 * time.Hour),
			Version:         1,
			Status:          models.TenderStatusCreated,
		},
	}
	for _, tender := range tenders {
		_, err := repo.Create(tender)
		assert.NoError(t, err)
	}
	retrievedTenders, err := repo.ListByOrganizationID(organizationID)

	assert.NoError(t, err)
	assert.NotNil(t, retrievedTenders)
	assert.Len(t, retrievedTenders, 2)

	_, err = db.Exec("DELETE FROM tenders")
	assert.NoError(t, err)
	retrievedTenders, err = repo.ListByOrganizationID(organizationID)
	assert.NoError(t, err)
	assert.Empty(t, retrievedTenders)
}
*/
