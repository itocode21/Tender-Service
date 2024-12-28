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
	// Create a test organization
	err = db.QueryRow(`INSERT INTO organizations (name, description, type, created_at, updated_at) VALUES ('Test Organization', 'Test Organization Description', $1, NOW(), NOW()) RETURNING id`, "LLC").Scan(&organizationID)
	if err != nil {
		panic(err)
	}

	err = db.QueryRow(`INSERT INTO tenders (name, description, organization_id, publication_date, end_date, version, status, created_at, updated_at) VALUES ('Test Tender', 'Test Tender Description', $1, NOW(), NOW(), 1, $2, NOW(), NOW()) RETURNING id`, organizationID, models.TenderStatusCreated).Scan(&tenderID)
	if err != nil {
		panic(err)
	}

	// Очистка таблиц перед каждым тестом
	_, err = db.Exec("DELETE FROM proposals")
	if err != nil {
		panic(err)
	}

}

func teardown() {
	_ = db.Close()
}

func TestProposalService_Create(t *testing.T) {
	setup()
	defer teardown()
	proposalRepo := repository.NewProposalRepository(db)
	tenderRepo := repository.NewTenderRepository(db)
	service := services.NewProposalService(proposalRepo, tenderRepo)

	testProposal := &models.Proposal{
		TenderID:       tenderID,
		OrganizationID: organizationID,
		Description:    "Test Proposal",
	}

	id, err := service.Create(testProposal)
	assert.NoError(t, err)
	assert.NotEqual(t, int64(0), id)
	assert.Equal(t, models.ProposalStatusCreated, testProposal.Status)

	testProposal = &models.Proposal{
		TenderID:       tenderID,
		OrganizationID: organizationID,
		Description:    "",
	}

	_, err = service.Create(testProposal)
	assert.Error(t, err)
	assert.EqualError(t, err, "описание предложения не может быть пустым")

	testProposal = &models.Proposal{
		TenderID:       0,
		OrganizationID: organizationID,
		Description:    "Test Proposal",
	}
	_, err = service.Create(testProposal)
	assert.Error(t, err)
	assert.EqualError(t, err, "необходимо указать тендер")

	testProposal = &models.Proposal{
		TenderID:       tenderID,
		OrganizationID: 0,
		Description:    "Test Proposal",
	}
	_, err = service.Create(testProposal)
	assert.Error(t, err)
	assert.EqualError(t, err, "необходимо указать организацию")
}

func TestProposalService_GetByID(t *testing.T) {
	setup()
	defer teardown()

	proposalRepo := repository.NewProposalRepository(db)
	tenderRepo := repository.NewTenderRepository(db)
	service := services.NewProposalService(proposalRepo, tenderRepo)

	testProposal := &models.Proposal{
		TenderID:       tenderID,
		OrganizationID: organizationID,
		Description:    "Test Proposal",
	}

	id, _ := proposalRepo.Create(testProposal)

	proposal, err := service.GetByID(id)
	assert.NoError(t, err)
	assert.Equal(t, testProposal.Description, proposal.Description)
	assert.Equal(t, testProposal.OrganizationID, proposal.OrganizationID)
	assert.Equal(t, testProposal.TenderID, proposal.TenderID)
	_, err = service.GetByID(id + 1)
	assert.Error(t, err)
	assert.EqualError(t, err, "proposal not found")
}

func TestProposalService_Update(t *testing.T) {
	setup()
	defer teardown()

	proposalRepo := repository.NewProposalRepository(db)
	tenderRepo := repository.NewTenderRepository(db)
	service := services.NewProposalService(proposalRepo, tenderRepo)

	testProposal := &models.Proposal{
		TenderID:       tenderID,
		OrganizationID: organizationID,
		Description:    "Test Proposal",
	}

	id, _ := proposalRepo.Create(testProposal)
	retrievedProposal, _ := proposalRepo.GetByID(id)

	updatedTestProposal := &models.Proposal{
		ID:             retrievedProposal.ID,
		TenderID:       retrievedProposal.TenderID,
		OrganizationID: retrievedProposal.OrganizationID,
		Description:    "Updated Test Proposal",
		Version:        retrievedProposal.Version + 1,
		UpdatedAt:      time.Now(),
	}

	err := service.Update(updatedTestProposal)
	assert.NoError(t, err)

	updatedProposal, _ := proposalRepo.GetByID(id)
	assert.Equal(t, updatedTestProposal.Description, updatedProposal.Description)
	assert.Equal(t, updatedTestProposal.Version, updatedProposal.Version)
	_, err = service.GetByID(id + 1)
	assert.Error(t, err)
	assert.EqualError(t, err, "proposal not found")

}

func TestProposalService_Delete(t *testing.T) {
	setup()
	defer teardown()

	proposalRepo := repository.NewProposalRepository(db)
	tenderRepo := repository.NewTenderRepository(db)
	service := services.NewProposalService(proposalRepo, tenderRepo)

	testProposal := &models.Proposal{
		TenderID:       tenderID,
		OrganizationID: organizationID,
		Description:    "Test Proposal",
	}

	id, _ := proposalRepo.Create(testProposal)

	err := service.Delete(id)
	assert.NoError(t, err)

	_, err = proposalRepo.GetByID(id)
	assert.Error(t, err)
	assert.EqualError(t, err, "proposal not found")

	err = service.Delete(id + 1)
	assert.NoError(t, err)
}

func TestProposalService_List(t *testing.T) {
	setup()
	defer teardown()

	proposalRepo := repository.NewProposalRepository(db)
	tenderRepo := repository.NewTenderRepository(db)
	service := services.NewProposalService(proposalRepo, tenderRepo)

	proposals := []*models.Proposal{
		{
			TenderID:       tenderID,
			OrganizationID: organizationID,
			Description:    "Test Proposal 1",
		},
		{
			TenderID:       tenderID,
			OrganizationID: organizationID,
			Description:    "Test Proposal 2",
		},
	}

	for _, proposal := range proposals {
		proposalRepo.Create(proposal)
	}

	retrievedProposals, err := service.List()
	assert.NoError(t, err)
	assert.NotNil(t, retrievedProposals)
	assert.Len(t, retrievedProposals, len(proposals))

	_, err = db.Exec("DELETE FROM proposals")
	assert.NoError(t, err)
	retrievedProposals, err = service.List()
	assert.NoError(t, err)
	assert.Empty(t, retrievedProposals)
}

func TestProposalService_Publish(t *testing.T) {
	setup()
	defer teardown()

	proposalRepo := repository.NewProposalRepository(db)
	tenderRepo := repository.NewTenderRepository(db)
	service := services.NewProposalService(proposalRepo, tenderRepo)

	testProposal := &models.Proposal{
		TenderID:       tenderID,
		OrganizationID: organizationID,
		Description:    "Test Proposal",
		Status:         models.ProposalStatusCreated,
	}
	id, _ := proposalRepo.Create(testProposal)

	err := service.Publish(id)
	assert.NoError(t, err)

	retrievedProposal, _ := proposalRepo.GetByID(id)
	assert.Equal(t, models.ProposalStatusPublished, retrievedProposal.Status)

	err = service.Publish(id)
	assert.Error(t, err)
	assert.EqualError(t, err, "предложение уже опубликовано")

	err = service.Publish(id + 1)
	assert.Error(t, err)
	assert.EqualError(t, err, "предложение не найдено: proposal not found")
}

func TestProposalService_Accept(t *testing.T) {
	setup()
	defer teardown()

	proposalRepo := repository.NewProposalRepository(db)
	tenderRepo := repository.NewTenderRepository(db)
	service := services.NewProposalService(proposalRepo, tenderRepo)
	testTender, _ := tenderRepo.GetByID(tenderID)

	testProposal := &models.Proposal{
		TenderID:       tenderID,
		OrganizationID: organizationID,
		Description:    "Test Proposal",
		Status:         models.ProposalStatusCreated,
	}
	id, _ := proposalRepo.Create(testProposal)

	err := service.Accept(id)
	assert.NoError(t, err)

	retrievedProposal, _ := proposalRepo.GetByID(id)
	retrievedTender, _ := tenderRepo.GetByID(tenderID)

	assert.Equal(t, models.ProposalStatusAccepted, retrievedProposal.Status)
	assert.Equal(t, models.TenderStatusClosed, retrievedTender.Status)

	err = service.Accept(id)
	assert.Error(t, err)
	assert.EqualError(t, err, "невозможно принять предложение для закрытого тендера")

	testTender.Status = models.TenderStatusCancelled
	testProposal.Status = models.ProposalStatusCreated
	_ = tenderRepo.Update(testTender)
	err = service.Accept(id)
	assert.Error(t, err)
	assert.EqualError(t, err, "невозможно принять предложение для отмененного тендера")

	testTender.Status = models.TenderStatusClosed
	testProposal.Status = models.ProposalStatusCreated
	_ = tenderRepo.Update(testTender)
	err = service.Accept(id)
	assert.Error(t, err)
	assert.EqualError(t, err, "невозможно принять предложение для закрытого тендера")
	err = service.Accept(id + 1)
	assert.Error(t, err)
	assert.EqualError(t, err, "предложение не найдено: proposal not found")
}

func TestProposalService_Reject(t *testing.T) {
	setup()
	defer teardown()
	proposalRepo := repository.NewProposalRepository(db)
	tenderRepo := repository.NewTenderRepository(db)
	service := services.NewProposalService(proposalRepo, tenderRepo)

	testProposal := &models.Proposal{
		TenderID:       tenderID,
		OrganizationID: organizationID,
		Description:    "Test Proposal",
		Status:         models.ProposalStatusCreated,
	}
	id, _ := proposalRepo.Create(testProposal)
	err := service.Reject(id)
	assert.NoError(t, err)

	retrievedProposal, _ := proposalRepo.GetByID(id)
	assert.Equal(t, models.ProposalStatusRejected, retrievedProposal.Status)

	err = service.Reject(id)
	assert.Error(t, err)
	assert.EqualError(t, err, "предложение уже отклонено")

	testProposal.Status = models.ProposalStatusAccepted
	_ = proposalRepo.Update(testProposal)
	err = service.Reject(id)
	assert.Error(t, err)
	assert.EqualError(t, err, "невозможно отклонить принятое предложение")

	err = service.Reject(id + 1)
	assert.Error(t, err)
	assert.EqualError(t, err, "предложение не найдено: proposal not found")
}

func TestProposalService_Cancel(t *testing.T) {
	setup()
	defer teardown()

	proposalRepo := repository.NewProposalRepository(db)
	tenderRepo := repository.NewTenderRepository(db)
	service := services.NewProposalService(proposalRepo, tenderRepo)

	testProposal := &models.Proposal{
		TenderID:       tenderID,
		OrganizationID: organizationID,
		Description:    "Test Proposal",
		Status:         models.ProposalStatusCreated,
	}
	id, _ := proposalRepo.Create(testProposal)

	err := service.Cancel(id)
	assert.NoError(t, err)

	retrievedProposal, _ := proposalRepo.GetByID(id)
	assert.Equal(t, models.ProposalStatusCancelled, retrievedProposal.Status)

	err = service.Cancel(id)
	assert.Error(t, err)
	assert.EqualError(t, err, "предложение уже отменено")

	testProposal.Status = models.ProposalStatusAccepted
	_ = proposalRepo.Update(testProposal)
	err = service.Cancel(id)
	assert.Error(t, err)
	assert.EqualError(t, err, "невозможно отменить принятое предложение")

	err = service.Cancel(id + 1)
	assert.Error(t, err)
	assert.EqualError(t, err, "предложение не найдено: proposal not found")

}

func TestProposalService_GetByTenderID(t *testing.T) {
	setup()
	defer teardown()

	proposalRepo := repository.NewProposalRepository(db)
	tenderRepo := repository.NewTenderRepository(db)
	service := services.NewProposalService(proposalRepo, tenderRepo)
	tenderID2 := int64(2)

	err := db.QueryRow(`INSERT INTO tenders (name, description, organization_id, publication_date, end_date, version, status, created_at, updated_at) VALUES ('Test Tender 2', 'Test Tender Description 2', $1, NOW(), NOW(), 1, $2, NOW(), NOW()) RETURNING id`, organizationID, models.TenderStatusCreated).Scan(&tenderID2)
	if err != nil {
		panic(err)
	}
	proposals := []*models.Proposal{
		{
			TenderID:       tenderID,
			OrganizationID: organizationID,
			Description:    "Test Proposal 1",
		},
		{
			TenderID:       tenderID,
			OrganizationID: organizationID,
			Description:    "Test Proposal 2",
		},
		{
			TenderID:       tenderID2,
			OrganizationID: organizationID,
			Description:    "Test Proposal 3",
		},
	}

	for _, proposal := range proposals {
		proposalRepo.Create(proposal)
	}
	retrievedProposals, err := service.GetByTenderID(tenderID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedProposals)
	assert.Len(t, retrievedProposals, 2)

	_, err = service.GetByTenderID(tenderID + 100)
	assert.NoError(t, err)
}
*/
