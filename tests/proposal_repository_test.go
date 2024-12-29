package tests

/*
import (
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/itocode21/Tender-Service/internal/models"
	"github.com/itocode21/Tender-Service/internal/repository"
	postgresqldb "github.com/itocode21/Tender-Service/postgresql-db"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var db *sql.DB
var organizationID int64
var tenderID int64

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

	err = db.QueryRow(`INSERT INTO tenders (name, description, organization_id, publication_date, end_date, version, status, created_at, updated_at) VALUES ('Test Tender', 'Test Tender Description', $1, NOW(), NOW(), 1, $2, NOW(), NOW()) RETURNING id`, organizationID, models.TenderStatusCreated).Scan(&tenderID)
	if err != nil {
		panic(err)
	}

	// Очистка таблицы proposals перед каждым тестом
	_, err = db.Exec("DELETE FROM proposals")
	if err != nil {
		panic(err)
	}
}

func teardown() {
	_ = db.Close()
}

func TestProposalRepository_Create(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewProposalRepository(db)

	proposal := &models.Proposal{
		TenderID:        tenderID,
		OrganizationID:  organizationID,
		Description:     "Test Proposal",
		PublicationDate: time.Now(),
		Price:           100.00,
		Version:         1,
	}
	id, err := repo.Create(proposal)
	assert.NoError(t, err)
	assert.NotEqual(t, int64(0), id)

	retrievedProposal, err := repo.GetByID(id)
	assert.NoError(t, err)
	assert.Equal(t, models.ProposalStatusCreated, retrievedProposal.Status)
}

func TestProposalRepository_GetByID(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewProposalRepository(db)

	proposal := &models.Proposal{
		TenderID:        tenderID,
		OrganizationID:  organizationID,
		Description:     "Test Proposal",
		PublicationDate: time.Now(),
		Price:           100.00,
		Version:         1,
	}
	id, _ := repo.Create(proposal)

	retrievedProposal, err := repo.GetByID(id)

	assert.NoError(t, err)
	assert.Equal(t, proposal.Description, retrievedProposal.Description)
	assert.Equal(t, proposal.TenderID, retrievedProposal.TenderID)
	assert.Equal(t, proposal.OrganizationID, retrievedProposal.OrganizationID)
	assert.Equal(t, proposal.Price, retrievedProposal.Price)
	assert.Equal(t, proposal.Version, retrievedProposal.Version)

	_, err = repo.GetByID(id + 1)
	assert.Error(t, err)
	assert.EqualError(t, err, "proposal not found")
}

func TestProposalRepository_Update(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewProposalRepository(db)

	proposal := &models.Proposal{
		TenderID:        tenderID,
		OrganizationID:  organizationID,
		Description:     "Test Proposal",
		PublicationDate: time.Now(),
		Price:           100.00,
		Version:         1,
	}
	id, _ := repo.Create(proposal)
	retrievedProposal, _ := repo.GetByID(id)

	retrievedProposal.Description = "Updated Test Proposal"
	retrievedProposal.Price = 200.00
	err := repo.Update(retrievedProposal)
	assert.NoError(t, err)

	updatedProposal, _ := repo.GetByID(id)
	assert.Equal(t, retrievedProposal.Description, updatedProposal.Description)
	assert.Equal(t, retrievedProposal.Price, updatedProposal.Price)
	assert.Equal(t, retrievedProposal.Version, updatedProposal.Version)
}

func TestProposalRepository_Delete(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewProposalRepository(db)

	proposal := &models.Proposal{
		TenderID:        tenderID,
		OrganizationID:  organizationID,
		Description:     "Test Proposal",
		PublicationDate: time.Now(),
		Price:           100.00,
		Version:         1,
	}
	id, _ := repo.Create(proposal)
	err := repo.Delete(id)
	assert.NoError(t, err)

	_, err = repo.GetByID(id)
	assert.Error(t, err)
	assert.EqualError(t, err, "proposal not found")

	err = repo.Delete(id)
	assert.NoError(t, err)
}

func TestProposalRepository_List(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewProposalRepository(db)

	proposals := []*models.Proposal{
		{
			TenderID:        tenderID,
			OrganizationID:  organizationID,
			Description:     "Test Proposal 1",
			PublicationDate: time.Now(),
			Price:           100.00,
			Version:         1,
		},
		{
			TenderID:        tenderID,
			OrganizationID:  organizationID,
			Description:     "Test Proposal 2",
			PublicationDate: time.Now(),
			Price:           200.00,
			Version:         1,
		},
	}
	for _, proposal := range proposals {
		repo.Create(proposal)
	}

	retrievedProposals, err := repo.List()
	assert.NoError(t, err)
	assert.NotNil(t, retrievedProposals)
	assert.Len(t, retrievedProposals, len(proposals))

	_, err = db.Exec("DELETE FROM proposals")
	assert.NoError(t, err)
	retrievedProposals, err = repo.List()
	assert.NoError(t, err)
	assert.Empty(t, retrievedProposals)
}

func TestProposalRepository_GetByTenderID(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewProposalRepository(db)
	tenderID2 := int64(2)
	err := db.QueryRow(`INSERT INTO tenders (name, description, organization_id, publication_date, end_date, version, status, created_at, updated_at) VALUES ('Test Tender 2', 'Test Tender Description 2', $1, NOW(), NOW(), 1, $2, NOW(), NOW()) RETURNING id`, organizationID, models.TenderStatusCreated).Scan(&tenderID2)
	if err != nil {
		panic(err)
	}
	proposals := []*models.Proposal{
		{
			TenderID:        tenderID,
			OrganizationID:  organizationID,
			Description:     "Test Proposal 1",
			PublicationDate: time.Now(),
			Price:           100.00,
			Version:         1,
		},
		{
			TenderID:        tenderID,
			OrganizationID:  organizationID,
			Description:     "Test Proposal 2",
			PublicationDate: time.Now(),
			Price:           200.00,
			Version:         1,
		},
		{
			TenderID:        tenderID2,
			OrganizationID:  organizationID,
			Description:     "Test Proposal 3",
			PublicationDate: time.Now(),
			Price:           300.00,
			Version:         1,
		},
	}
	for _, proposal := range proposals {
		repo.Create(proposal)
	}

	retrievedProposals, err := repo.GetByTenderID(tenderID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedProposals)
	assert.Len(t, retrievedProposals, 2)

	_, err = db.Exec("DELETE FROM proposals")
	assert.NoError(t, err)
	retrievedProposals, err = repo.GetByTenderID(tenderID)
	assert.NoError(t, err)
	assert.Empty(t, retrievedProposals)
}
*/
