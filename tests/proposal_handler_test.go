package tests

/*
import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/itocode21/Tender-Service/internal/handlers"
	"github.com/itocode21/Tender-Service/internal/models"
	"github.com/itocode21/Tender-Service/internal/repository"
	"github.com/itocode21/Tender-Service/internal/services"
	postgresqldb "github.com/itocode21/Tender-Service/postgresql-db"
)

var db *sql.DB

func setup() {
	godotenv.Load()
	conn, err := postgresqldb.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database connection: %v", err)
	}
	db = conn
}

func teardown() {
	db.Close()
}

func clearDatabase(db *sql.DB) {
	_, err := db.Exec("DELETE FROM organization_responsibles")
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
	_, err = db.Exec("DELETE FROM organizations")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("DELETE FROM employees")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("ALTER SEQUENCE employees_id_seq RESTART WITH 1")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("ALTER SEQUENCE organizations_id_seq RESTART WITH 1")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("ALTER SEQUENCE tenders_id_seq RESTART WITH 1")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("ALTER SEQUENCE proposals_id_seq RESTART WITH 1")
	if err != nil {
		panic(err)
	}
}

func TestProposalHandlers(t *testing.T) {
	setup()
	defer teardown()

	proposalRepo := repository.NewProposalRepository(db)
	tenderRepo := repository.NewTenderRepository(db)
	proposalService := services.NewProposalService(proposalRepo, tenderRepo)
	handler := handlers.NewProposalHandler(proposalService)

	// Create a new mux router
	router := mux.NewRouter()
	router.HandleFunc("/proposals", handler.CreateProposalHandler).Methods("POST")
	router.HandleFunc("/proposals/{id}", handler.GetProposalHandler).Methods("GET")
	router.HandleFunc("/proposals/{id}", handler.UpdateProposalHandler).Methods("PUT")
	router.HandleFunc("/proposals/{id}", handler.DeleteProposalHandler).Methods("DELETE")
	router.HandleFunc("/proposals", handler.ListProposalsHandler).Methods("GET")
	router.HandleFunc("/tenders/{tender_id}/proposals", handler.GetProposalsByTenderHandler).Methods("GET")
	router.HandleFunc("/proposals/{id}/publish", handler.PublishProposalHandler).Methods("PUT")
	router.HandleFunc("/proposals/{id}/accept", handler.AcceptProposalHandler).Methods("PUT")
	router.HandleFunc("/proposals/{id}/reject", handler.RejectProposalHandler).Methods("PUT")
	router.HandleFunc("/proposals/{id}/cancel", handler.CancelProposalHandler).Methods("PUT")

	t.Run("CreateProposalHandler_Success", func(t *testing.T) {
		clearDatabase(db)

		org := models.Organization{
			Id:          1,
			Name:        "Test Org",
			Description: "Test organization",
			Type:        "LLC",
		}
		orgRepo := repository.NewOrganizationRepository(db)
		_, err := orgRepo.Create(&org)
		if err != nil {
			panic(err)
		}

		emp := models.User{
			Id:        1,
			Username:  "test_employee_" + strconv.Itoa(int(time.Now().UnixNano())),
			FirstName: "Test",
			LastName:  "Employee",
		}

		empRepo := repository.NewEmployeeRepository(db)
		_, err = empRepo.Create(&emp)
		if err != nil {
			panic(err)
		}

		tender := models.Tender{
			ID:              1,
			OrganizationID:  org.Id,
			Name:            "Test tender",
			Description:     "Test tender",
			PublicationDate: time.Now(),
			EndDate:         time.Now().Add(time.Hour),
			Version:         1,
		}

		tenderRepo := repository.NewTenderRepository(db)
		_, err = tenderRepo.Create(&tender)
		if err != nil {
			panic(err)
		}

		proposal := models.Proposal{
			TenderID:       tender.ID,
			OrganizationID: org.Id,
			Description:    "Test proposal",
			Price:          100.00,
			Version:        1,
		}

		body, _ := json.Marshal(proposal)

		req := httptest.NewRequest("POST", "/proposals", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
		}
		var createdProposal models.Proposal
		err = json.Unmarshal(w.Body.Bytes(), &createdProposal)
		if err != nil {
			t.Errorf("error unmarshaling created proposal")
		}
		if createdProposal.ID == 0 {
			t.Error("Expected created proposal to have ID, got 0")
		}

	})

	t.Run("CreateProposalHandler_Error", func(t *testing.T) {
		clearDatabase(db)
		proposal := models.Proposal{
			Description: "",
		}
		body, _ := json.Marshal(proposal)

		req := httptest.NewRequest("POST", "/proposals", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}

	})

	t.Run("GetProposalHandler_Success", func(t *testing.T) {
		clearDatabase(db)
		org := models.Organization{
			Id:          1,
			Name:        "Test Org",
			Description: "Test organization",
			Type:        "LLC",
		}
		orgRepo := repository.NewOrganizationRepository(db)
		_, err := orgRepo.Create(&org)
		if err != nil {
			panic(err)
		}

		emp := models.User{
			Id:        1,
			Username:  "test_employee_" + strconv.Itoa(int(time.Now().UnixNano())),
			FirstName: "Test",
			LastName:  "Employee",
		}

		empRepo := repository.NewEmployeeRepository(db)
		_, err = empRepo.Create(&emp)
		if err != nil {
			panic(err)
		}

		tender := models.Tender{
			ID:              1,
			OrganizationID:  org.Id,
			Name:            "Test tender",
			Description:     "Test tender",
			PublicationDate: time.Now(),
			EndDate:         time.Now().Add(time.Hour),
			Version:         1,
		}

		tenderRepo := repository.NewTenderRepository(db)
		_, err = tenderRepo.Create(&tender)
		if err != nil {
			panic(err)
		}

		proposal := models.Proposal{
			TenderID:       tender.ID,
			OrganizationID: org.Id,
			Description:    "Test proposal",
			Price:          100.00,
			Version:        1,
		}
		createdID, _ := proposalService.Create(&proposal)
		req := httptest.NewRequest("GET", fmt.Sprintf("/proposals/%d", createdID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %v, got %v", http.StatusOK, w.Code)
		}

		var getProposal models.Proposal
		err = json.Unmarshal(w.Body.Bytes(), &getProposal)
		if err != nil {
			t.Errorf("error unmarshaling get proposal")
		}
		if getProposal.ID != createdID {
			t.Errorf("expected id %v, got %v", createdID, getProposal.ID)
		}
	})
	t.Run("GetProposalHandler_Error", func(t *testing.T) {
		clearDatabase(db)
		req := httptest.NewRequest("GET", "/proposals/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %v, got %v", http.StatusNotFound, w.Code)
		}

	})
	t.Run("UpdateProposalHandler_Success", func(t *testing.T) {
		clearDatabase(db)
		org := models.Organization{
			Id:          1,
			Name:        "Test Org",
			Description: "Test organization",
			Type:        "LLC",
		}
		orgRepo := repository.NewOrganizationRepository(db)
		_, err := orgRepo.Create(&org)
		if err != nil {
			panic(err)
		}

		emp := models.User{
			Id:        1,
			Username:  "test_employee_" + strconv.Itoa(int(time.Now().UnixNano())),
			FirstName: "Test",
			LastName:  "Employee",
		}

		empRepo := repository.NewEmployeeRepository(db)
		_, err = empRepo.Create(&emp)
		if err != nil {
			panic(err)
		}

		tender := models.Tender{
			ID:              1,
			OrganizationID:  org.Id,
			Name:            "Test tender",
			Description:     "Test tender",
			PublicationDate: time.Now(),
			EndDate:         time.Now().Add(time.Hour),
			Version:         1,
		}

		tenderRepo := repository.NewTenderRepository(db)
		_, err = tenderRepo.Create(&tender)
		if err != nil {
			panic(err)
		}
		testProposal := models.Proposal{
			TenderID:       tender.ID,
			OrganizationID: org.Id,
			Description:    "Test proposal",
			Price:          100.00,
			Version:        1,
		}
		createdID, _ := proposalRepo.Create(&testProposal)

		updatedProposal := models.Proposal{
			ID:             createdID,
			TenderID:       tender.ID,
			OrganizationID: org.Id,
			Description:    "Test proposal updated",
			Price:          110.00,
		}

		body, _ := json.Marshal(updatedProposal)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/proposals/%d", updatedProposal.ID), bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status %v, got %v", http.StatusOK, w.Code)
		}
	})

	t.Run("UpdateProposalHandler_Error", func(t *testing.T) {
		clearDatabase(db)
		org := models.Organization{
			Id:          1,
			Name:        "Test Org",
			Description: "Test organization",
			Type:        "LLC",
		}
		orgRepo := repository.NewOrganizationRepository(db)
		_, err := orgRepo.Create(&org)
		if err != nil {
			panic(err)
		}

		emp := models.User{
			Id:        1,
			Username:  "test_employee_" + strconv.Itoa(int(time.Now().UnixNano())),
			FirstName: "Test",
			LastName:  "Employee",
		}

		empRepo := repository.NewEmployeeRepository(db)
		_, err = empRepo.Create(&emp)
		if err != nil {
			panic(err)
		}

		tender := models.Tender{
			ID:              1,
			OrganizationID:  org.Id,
			Name:            "Test tender",
			Description:     "Test tender",
			PublicationDate: time.Now(),
			EndDate:         time.Now().Add(time.Hour),
			Version:         1,
		}

		tenderRepo := repository.NewTenderRepository(db)
		_, err = tenderRepo.Create(&tender)
		if err != nil {
			panic(err)
		}
		testProposal := models.Proposal{
			TenderID:       tender.ID,
			OrganizationID: org.Id,
			Description:    "Test proposal",
			Price:          100.00,
			Version:        1,
		}
		createdID, _ := proposalRepo.Create(&testProposal)

		updatedProposal := models.Proposal{
			ID:          createdID,
			Description: "",
		}
		body, _ := json.Marshal(updatedProposal)

		req := httptest.NewRequest("PUT", fmt.Sprintf("/proposals/%d", createdID), bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %v, got %v", http.StatusBadRequest, w.Code)
		}
		var errorResp map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &errorResp)
		if err != nil {
			t.Fatalf("Cannot unmarshal error response: %v", err)
		}
		if _, ok := errorResp["error"]; !ok {
			t.Error("Expected error message in response, but got none")
		}
	})
	t.Run("DeleteProposalHandler_Success", func(t *testing.T) {
		clearDatabase(db)
		org := models.Organization{
			Id:          1,
			Name:        "Test Org",
			Description: "Test organization",
			Type:        "LLC",
		}
		orgRepo := repository.NewOrganizationRepository(db)
		_, err := orgRepo.Create(&org)
		if err != nil {
			panic(err)
		}

		emp := models.User{
			Id:        1,
			Username:  "test_employee_" + strconv.Itoa(int(time.Now().UnixNano())),
			FirstName: "Test",
			LastName:  "Employee",
		}

		empRepo := repository.NewEmployeeRepository(db)
		_, err = empRepo.Create(&emp)
		if err != nil {
			panic(err)
		}

		tender := models.Tender{
			ID:              1,
			OrganizationID:  org.Id,
			Name:            "Test tender",
			Description:     "Test tender",
			PublicationDate: time.Now(),
			EndDate:         time.Now().Add(time.Hour),
			Version:         1,
		}

		tenderRepo := repository.NewTenderRepository(db)
		_, err = tenderRepo.Create(&tender)
		if err != nil {
			panic(err)
		}
		testProposal := models.Proposal{
			TenderID:       tender.ID,
			OrganizationID: org.Id,
			Description:    "Test proposal",
			Price:          100.00,
			Version:        1,
		}
		createdID, _ := proposalRepo.Create(&testProposal)

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/proposals/%d", createdID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status %v, got %v", http.StatusNoContent, w.Code)
		}
	})
	t.Run("ListProposalsHandler_Success", func(t *testing.T) {
		clearDatabase(db)
		org := models.Organization{
			Id:          1,
			Name:        "Test Org",
			Description: "Test organization",
			Type:        "LLC",
		}
		orgRepo := repository.NewOrganizationRepository(db)
		_, err := orgRepo.Create(&org)
		if err != nil {
			panic(err)
		}

		emp := models.User{
			Id:        1,
			Username:  "test_employee_" + strconv.Itoa(int(time.Now().UnixNano())),
			FirstName: "Test",
			LastName:  "Employee",
		}

		empRepo := repository.NewEmployeeRepository(db)
		_, err = empRepo.Create(&emp)
		if err != nil {
			panic(err)
		}

		tender := models.Tender{
			ID:              1,
			OrganizationID:  org.Id,
			Name:            "Test tender",
			Description:     "Test tender",
			PublicationDate: time.Now(),
			EndDate:         time.Now().Add(time.Hour),
			Version:         1,
		}

		tenderRepo := repository.NewTenderRepository(db)
		_, err = tenderRepo.Create(&tender)
		if err != nil {
			panic(err)
		}

		proposal1 := models.Proposal{
			TenderID:       tender.ID,
			OrganizationID: org.Id,
			Description:    "Test proposal 1" + strconv.Itoa(int(time.Now().UnixNano())),
			Price:          100.00,
			Version:        1,
		}
		proposal2 := models.Proposal{
			TenderID:       tender.ID,
			OrganizationID: org.Id,
			Description:    "Test proposal 2" + strconv.Itoa(int(time.Now().UnixNano())),
			Price:          100.00,
			Version:        1,
		}
		_, _ = proposalRepo.Create(&proposal1)
		_, _ = proposalRepo.Create(&proposal2)

		req := httptest.NewRequest("GET", "/proposals", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %v, got %v", http.StatusOK, w.Code)
		}
		var proposals []models.Proposal
		err = json.Unmarshal(w.Body.Bytes(), &proposals)
		if err != nil {
			t.Errorf("error unmarshaling list proposals %v", err)
		}
		if len(proposals) != 2 {
			t.Fatalf("expected proposals len to be %v, got %v", 2, len(proposals))
		}
		if proposals[0].Description != proposal1.Description {
			t.Errorf("expected first proposal description to be '%v', got '%v'", proposal1.Description, proposals[0].Description)
		}
		if proposals[1].Description != proposal2.Description {
			t.Errorf("expected second proposal description to be '%v', got '%v'", proposal2.Description, proposals[1].Description)
		}
	})
	t.Run("GetProposalsByTenderHandler_Success", func(t *testing.T) {
		clearDatabase(db)
		org := models.Organization{
			Id:          1,
			Name:        "Test Org",
			Description: "Test organization",
			Type:        "LLC",
		}
		orgRepo := repository.NewOrganizationRepository(db)
		_, err := orgRepo.Create(&org)
		if err != nil {
			panic(err)
		}

		emp := models.User{
			Id:        1,
			Username:  "test_employee_" + strconv.Itoa(int(time.Now().UnixNano())),
			FirstName: "Test",
			LastName:  "Employee",
		}

		empRepo := repository.NewEmployeeRepository(db)
		_, err = empRepo.Create(&emp)
		if err != nil {
			panic(err)
		}

		tender := models.Tender{
			ID:              1,
			OrganizationID:  org.Id,
			Name:            "Test tender",
			Description:     "Test tender",
			PublicationDate: time.Now(),
			EndDate:         time.Now().Add(time.Hour),
			Version:         1,
		}

		tenderRepo := repository.NewTenderRepository(db)
		_, err = tenderRepo.Create(&tender)
		if err != nil {
			panic(err)
		}
		proposal1 := models.Proposal{
			TenderID:       tender.ID,
			OrganizationID: org.Id,
			Description:    "Test proposal 1 " + strconv.Itoa(int(time.Now().UnixNano())),
			Price:          100.00,
			Version:        1,
		}
		proposal2 := models.Proposal{
			TenderID:       tender.ID,
			OrganizationID: org.Id,
			Description:    "Test proposal 2 " + strconv.Itoa(int(time.Now().UnixNano())),
			Price:          100.00,
			Version:        1,
		}
		_, _ = proposalRepo.Create(&proposal1)
		_, _ = proposalRepo.Create(&proposal2)

		req := httptest.NewRequest("GET", fmt.Sprintf("/tenders/%d/proposals", tender.ID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %v, got %v", http.StatusOK, w.Code)
		}
		var proposals []models.Proposal
		err = json.Unmarshal(w.Body.Bytes(), &proposals)
		if err != nil {
			t.Errorf("error unmarshaling get proposal %v", err)
		}
		if len(proposals) != 2 {
			t.Errorf("expected proposals len to be %v, got %v", 2, len(proposals))
		}
		if proposals[0].Description != proposal1.Description {
			t.Errorf("expected first proposal description to be '%v', got '%v'", proposal1.Description, proposals[0].Description)
		}
		if proposals[1].Description != proposal2.Description {
			t.Errorf("expected second proposal description to be '%v', got '%v'", proposal2.Description, proposals[1].Description)
		}
	})
	t.Run("GetProposalsByTenderHandler_Error", func(t *testing.T) {
		clearDatabase(db)
		// Crucial: Create a tender (or get the ID of a tender if it already exists)
		org := models.Organization{
			Id:          1,
			Name:        "Test Org",
			Description: "Test organization",
			Type:        "LLC",
		}
		orgRepo := repository.NewOrganizationRepository(db)
		_, err := orgRepo.Create(&org)
		if err != nil {
			panic(err)
		}

		tenderRepo := repository.NewTenderRepository(db)
		_, err = tenderRepo.Create(&models.Tender{OrganizationID: org.Id})
		if err != nil {
			panic(err)
		}
		req := httptest.NewRequest("GET", "/tenders/999/proposals", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %v, got %v", http.StatusBadRequest, w.Code)
		}
		var errorResp map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &errorResp)
		if err != nil {
			t.Fatalf("Cannot unmarshal error response: %v", err)
		}
		if _, ok := errorResp["error"]; !ok {
			t.Error("Expected error message in response, but got none")
		}
	})
	t.Run("PublishProposalHandler_Success", func(t *testing.T) {
		clearDatabase(db)
		org := models.Organization{
			Id:          1,
			Name:        "Test Org",
			Description: "Test organization",
			Type:        "LLC",
		}
		orgRepo := repository.NewOrganizationRepository(db)
		_, err := orgRepo.Create(&org)
		if err != nil {
			panic(err)
		}

		emp := models.User{
			Id:        1,
			Username:  "test_employee_" + strconv.Itoa(int(time.Now().UnixNano())),
			FirstName: "Test",
			LastName:  "Employee",
		}

		empRepo := repository.NewEmployeeRepository(db)
		_, err = empRepo.Create(&emp)
		if err != nil {
			panic(err)
		}

		tender := models.Tender{
			ID:              1,
			OrganizationID:  org.Id,
			Name:            "Test tender",
			Description:     "Test tender",
			PublicationDate: time.Now(),
			EndDate:         time.Now().Add(time.Hour),
			Version:         1,
		}

		tenderRepo := repository.NewTenderRepository(db)
		_, err = tenderRepo.Create(&tender)
		if err != nil {
			panic(err)
		}

		proposal := models.Proposal{
			TenderID:       tender.ID,
			OrganizationID: org.Id,
			Description:    "Test proposal",
			Price:          100.00,
			Version:        1,
		}
		createdID, _ := proposalRepo.Create(&proposal)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/proposals/%d/publish", createdID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %v, got %v", http.StatusOK, w.Code)
		}
		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Error decoding response: %v", err)
		}

		if response["message"] != "proposal published" {
			t.Errorf("Expected message to be 'proposal published', got '%v'", response["message"])
		}

	})
	t.Run("PublishProposalHandler_Error", func(t *testing.T) {
		clearDatabase(db)
		req := httptest.NewRequest("PUT", "/proposals/999/publish", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %v, got %v", http.StatusBadRequest, w.Code)
		}

		var errorResp map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &errorResp)
		if err != nil {
			t.Fatalf("Cannot unmarshal error response: %v", err)
		}
		if _, ok := errorResp["error"]; !ok {
			t.Error("Expected error message in response, but got none")
		}
	})
	t.Run("AcceptProposalHandler_Success", func(t *testing.T) {
		clearDatabase(db)
		org := models.Organization{
			Id:          1,
			Name:        "Test Org",
			Description: "Test organization",
			Type:        "LLC",
		}
		orgRepo := repository.NewOrganizationRepository(db)
		_, err := orgRepo.Create(&org)
		if err != nil {
			panic(err)
		}

		emp := models.User{
			Id:        1,
			Username:  "test_employee_" + strconv.Itoa(int(time.Now().UnixNano())),
			FirstName: "Test",
			LastName:  "Employee",
		}

		empRepo := repository.NewEmployeeRepository(db)
		_, err = empRepo.Create(&emp)
		if err != nil {
			panic(err)
		}

		tender := models.Tender{
			ID:              1,
			OrganizationID:  org.Id,
			Name:            "Test tender",
			Description:     "Test tender",
			PublicationDate: time.Now(),
			EndDate:         time.Now().Add(time.Hour),
			Version:         1,
			Status:          (models.TenderStatusCreated),
		}

		tenderRepo := repository.NewTenderRepository(db)
		_, err = tenderRepo.Create(&tender)
		if err != nil {
			panic(err)
		}

		proposal := models.Proposal{
			TenderID:       tender.ID,
			OrganizationID: org.Id,
			Description:    "Test proposal",
			Price:          100.00,
			Version:        1,
			Status:         "created",
		}
		createdID, _ := proposalRepo.Create(&proposal)

		req := httptest.NewRequest("PUT", fmt.Sprintf("/proposals/%d/accept", createdID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %v, got %v", http.StatusOK, w.Code)
		}
		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Error decoding response: %v", err)
		}

		if response["message"] != "proposal accepted" {
			t.Errorf("Expected message to be 'proposal accepted', got '%v'", response["message"])
		}

		updatedTender, _ := tenderRepo.GetByID(tender.ID)
		if updatedTender.Status != (models.TenderStatusClosed) {
			t.Errorf("Expected tender status to be '%v', but got '%v'", models.TenderStatusClosed, updatedTender.Status)
		}
	})
	t.Run("AcceptProposalHandler_Error", func(t *testing.T) {
		clearDatabase(db)
		req := httptest.NewRequest("PUT", "/proposals/999/accept", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %v, got %v", http.StatusBadRequest, w.Code)
		}

		var errorResp map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &errorResp)
		if err != nil {
			t.Fatalf("Cannot unmarshal error response: %v", err)
		}
		if _, ok := errorResp["error"]; !ok {
			t.Error("Expected error message in response, but got none")
		}
	})
	t.Run("RejectProposalHandler_Success", func(t *testing.T) {
		clearDatabase(db)
		org := models.Organization{
			Id:          1,
			Name:        "Test Org",
			Description: "Test organization",
			Type:        "LLC",
		}
		orgRepo := repository.NewOrganizationRepository(db)
		_, err := orgRepo.Create(&org)
		if err != nil {
			panic(err)
		}

		emp := models.User{
			Id:        1,
			Username:  "test_employee_" + strconv.Itoa(int(time.Now().UnixNano())),
			FirstName: "Test",
			LastName:  "Employee",
		}

		empRepo := repository.NewEmployeeRepository(db)
		_, err = empRepo.Create(&emp)
		if err != nil {
			panic(err)
		}

		tender := models.Tender{
			ID:              1,
			OrganizationID:  org.Id,
			Name:            "Test tender",
			Description:     "Test tender",
			PublicationDate: time.Now(),
			EndDate:         time.Now().Add(time.Hour),
			Version:         1,
		}

		tenderRepo := repository.NewTenderRepository(db)
		_, err = tenderRepo.Create(&tender)
		if err != nil {
			panic(err)
		}

		proposal := models.Proposal{
			TenderID:       tender.ID,
			OrganizationID: org.Id,
			Description:    "Test proposal",
			Price:          100.00,
			Version:        1,
			Status:         "created",
		}
		createdID, _ := proposalRepo.Create(&proposal)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/proposals/%d/reject", createdID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %v, got %v", http.StatusOK, w.Code)
		}
		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Error decoding response: %v", err)
		}

		if response["message"] != "proposal rejected" {
			t.Errorf("Expected message to be 'proposal rejected', got '%v'", response["message"])
		}
	})
	t.Run("RejectProposalHandler_Error", func(t *testing.T) {
		clearDatabase(db)
		req := httptest.NewRequest("PUT", "/proposals/999/reject", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %v, got %v", http.StatusBadRequest, w.Code)
		}

		var errorResp map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &errorResp)
		if err != nil {
			t.Fatalf("Cannot unmarshal error response: %v", err)
		}
		if _, ok := errorResp["error"]; !ok {
			t.Error("Expected error message in response, but got none")
		}
	})
	t.Run("RejectProposalHandler_Success", func(t *testing.T) {
		clearDatabase(db)
		org := models.Organization{
			Id:          1,
			Name:        "Test Org",
			Description: "Test organization",
			Type:        "LLC",
		}
		orgRepo := repository.NewOrganizationRepository(db)
		_, err := orgRepo.Create(&org)
		if err != nil {
			panic(err)
		}

		emp := models.User{
			Id:        1,
			Username:  "test_employee_" + strconv.Itoa(int(time.Now().UnixNano())),
			FirstName: "Test",
			LastName:  "Employee",
		}

		empRepo := repository.NewEmployeeRepository(db)
		_, err = empRepo.Create(&emp)
		if err != nil {
			panic(err)
		}

		tender := models.Tender{
			ID:              1,
			OrganizationID:  org.Id,
			Name:            "Test tender",
			Description:     "Test tender",
			PublicationDate: time.Now(),
			EndDate:         time.Now().Add(time.Hour),
			Version:         1,
		}

		tenderRepo := repository.NewTenderRepository(db)
		_, err = tenderRepo.Create(&tender)
		if err != nil {
			panic(err)
		}

		proposal := models.Proposal{
			TenderID:       tender.ID,
			OrganizationID: org.Id,
			Description:    "Test proposal",
			Price:          100.00,
			Version:        1,
			Status:         "created",
		}
		createdID, _ := proposalRepo.Create(&proposal)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/proposals/%d/reject", createdID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %v, got %v", http.StatusOK, w.Code)
		}
		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Error decoding response: %v", err)
		}

		if response["message"] != "proposal rejected" {
			t.Errorf("Expected message to be 'proposal rejected', got '%v'", response["message"])
		}
	})
	t.Run("RejectProposalHandler_Error", func(t *testing.T) {
		clearDatabase(db)
		req := httptest.NewRequest("PUT", "/proposals/999/reject", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %v, got %v", http.StatusBadRequest, w.Code)
		}

		var errorResp map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &errorResp)
		if err != nil {
			t.Fatalf("Cannot unmarshal error response: %v", err)
		}
		if _, ok := errorResp["error"]; !ok {
			t.Error("Expected error message in response, but got none")
		}
	})
	t.Run("CancelProposalHandler_Success", func(t *testing.T) {
		clearDatabase(db)
		org := models.Organization{
			Id:          1,
			Name:        "Test Org",
			Description: "Test organization",
			Type:        "LLC",
		}
		orgRepo := repository.NewOrganizationRepository(db)
		_, err := orgRepo.Create(&org)
		if err != nil {
			panic(err)
		}

		emp := models.User{
			Id:        1,
			Username:  "test_employee_" + strconv.Itoa(int(time.Now().UnixNano())),
			FirstName: "Test",
			LastName:  "Employee",
		}

		empRepo := repository.NewEmployeeRepository(db)
		_, err = empRepo.Create(&emp)
		if err != nil {
			panic(err)
		}

		tender := models.Tender{
			ID:              1,
			OrganizationID:  org.Id,
			Name:            "Test tender",
			Description:     "Test tender",
			PublicationDate: time.Now(),
			EndDate:         time.Now().Add(time.Hour),
			Version:         1,
		}

		tenderRepo := repository.NewTenderRepository(db)
		_, err = tenderRepo.Create(&tender)
		if err != nil {
			panic(err)
		}

		proposal := models.Proposal{
			TenderID:       tender.ID,
			OrganizationID: org.Id,
			Description:    "Test proposal",
			Price:          100.00,
			Version:        1,
			Status:         "created",
		}
		createdID, _ := proposalRepo.Create(&proposal)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/proposals/%d/cancel", createdID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %v, got %v", http.StatusOK, w.Code)
		}
		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Error decoding response: %v", err)
		}

		if response["message"] != "proposal cancelled" {
			t.Errorf("Expected message to be 'proposal cancelled', got '%v'", response["message"])
		}
	})
	t.Run("CancelProposalHandler_Error", func(t *testing.T) {
		clearDatabase(db)
		req := httptest.NewRequest("PUT", "/proposals/999/cancel", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %v, got %v", http.StatusBadRequest, w.Code)
		}

		var errorResp map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &errorResp)
		if err != nil {
			t.Fatalf("Cannot unmarshal error response: %v", err)
		}
		if _, ok := errorResp["error"]; !ok {
			t.Error("Expected error message in response, but got none")
		}
	})
}
*/
