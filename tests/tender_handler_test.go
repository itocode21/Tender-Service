package tests

/*
import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"database/sql"

	"github.com/gorilla/mux"
	"github.com/itocode21/Tender-Service/internal/handlers"
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

func TestTenderHandler_CreateTenderHandler(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewTenderRepository(db)
	service := services.NewTenderService(repo)
	handler := handlers.NewTenderHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/tenders", handler.CreateTenderHandler).Methods("POST")

	tender := &models.Tender{
		Name:            "Test Tender",
		Description:     "Test Description",
		OrganizationID:  organizationID,
		PublicationDate: time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
	}

	jsonValue, _ := json.Marshal(tender)
	req, _ := http.NewRequest("POST", "/tenders", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var createdTender models.Tender
	json.NewDecoder(w.Body).Decode(&createdTender)
	assert.Equal(t, tender.Name, createdTender.Name)
	assert.Equal(t, tender.Description, createdTender.Description)
	assert.Equal(t, tender.OrganizationID, createdTender.OrganizationID)

	invalidTender := &models.Tender{
		Name:            "",
		Description:     "Test Description",
		OrganizationID:  organizationID,
		PublicationDate: time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
	}
	jsonValue, _ = json.Marshal(invalidTender)
	req, _ = http.NewRequest("POST", "/tenders", bytes.NewBuffer(jsonValue))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	invalidTender = &models.Tender{
		Name:            "Test Tender",
		Description:     "Test Description",
		OrganizationID:  0,
		PublicationDate: time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
	}
	jsonValue, _ = json.Marshal(invalidTender)
	req, _ = http.NewRequest("POST", "/tenders", bytes.NewBuffer(jsonValue))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTenderHandler_GetTenderHandler(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewTenderRepository(db)
	service := services.NewTenderService(repo)
	handler := handlers.NewTenderHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/tenders/{id}", handler.GetTenderHandler).Methods("GET")

	tender := &models.Tender{
		Name:            "Test Tender",
		Description:     "Test Description",
		OrganizationID:  organizationID,
		PublicationDate: time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
	}
	id, _ := service.Create(tender)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/tenders/%d", id), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var retrievedTender models.Tender
	json.NewDecoder(w.Body).Decode(&retrievedTender)
	assert.Equal(t, tender.Name, retrievedTender.Name)
	assert.Equal(t, tender.Description, retrievedTender.Description)
	assert.Equal(t, tender.OrganizationID, retrievedTender.OrganizationID)

	req, _ = http.NewRequest("GET", fmt.Sprintf("/tenders/%d", id+1), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	req, _ = http.NewRequest("GET", "/tenders/invalid-id", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestTenderHandler_UpdateTenderHandler(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewTenderRepository(db)
	service := services.NewTenderService(repo)
	handler := handlers.NewTenderHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/tenders/{id}", handler.UpdateTenderHandler).Methods("PUT")

	tender := &models.Tender{
		Name:            "Test Tender",
		Description:     "Test Description",
		OrganizationID:  organizationID,
		PublicationDate: time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
	}
	id, _ := service.Create(tender)

	updatedTender := &models.Tender{
		Name:            "Updated Test Tender",
		Description:     "Updated Test Description",
		OrganizationID:  organizationID,
		PublicationDate: time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
	}
	jsonValue, _ := json.Marshal(updatedTender)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/tenders/%d", id), bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var retrievedTender models.Tender
	json.NewDecoder(w.Body).Decode(&retrievedTender)
	assert.Equal(t, updatedTender.Name, retrievedTender.Name)
	assert.Equal(t, updatedTender.Description, retrievedTender.Description)
	assert.Equal(t, updatedTender.OrganizationID, retrievedTender.OrganizationID)

	invalidTender := &models.Tender{
		Name:            "",
		Description:     "Updated Test Description",
		OrganizationID:  organizationID,
		PublicationDate: time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
	}
	jsonValue, _ = json.Marshal(invalidTender)
	req, _ = http.NewRequest("PUT", fmt.Sprintf("/tenders/%d", id), bytes.NewBuffer(jsonValue))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	jsonValue, _ = json.Marshal(updatedTender)
	req, _ = http.NewRequest("PUT", "/tenders/invalid-id", bytes.NewBuffer(jsonValue))
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTenderHandler_DeleteTenderHandler(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewTenderRepository(db)
	service := services.NewTenderService(repo)
	handler := handlers.NewTenderHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/tenders/{id}", handler.DeleteTenderHandler).Methods("DELETE")

	tender := &models.Tender{
		Name:            "Test Tender",
		Description:     "Test Description",
		OrganizationID:  organizationID,
		PublicationDate: time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
	}
	id, _ := service.Create(tender)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/tenders/%d", id), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	req, _ = http.NewRequest("DELETE", "/tenders/invalid-id", nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTenderHandler_ListTendersHandler(t *testing.T) {
	setup()
	defer teardown()
	repo := repository.NewTenderRepository(db)
	service := services.NewTenderService(repo)
	handler := handlers.NewTenderHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/tenders", handler.ListTendersHandler).Methods("GET")
	tenders := []*models.Tender{
		{
			Name:            "Test Tender 1",
			Description:     "Test Description 1",
			OrganizationID:  organizationID,
			PublicationDate: time.Now(),
			EndDate:         time.Now().Add(24 * time.Hour),
		},
		{
			Name:            "Test Tender 2",
			Description:     "Test Description 2",
			OrganizationID:  organizationID,
			PublicationDate: time.Now(),
			EndDate:         time.Now().Add(24 * time.Hour),
		},
	}
	for _, tender := range tenders {
		service.Create(tender)
	}

	req, _ := http.NewRequest("GET", "/tenders", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var retrievedTenders []*models.Tender
	json.NewDecoder(w.Body).Decode(&retrievedTenders)
	assert.Len(t, retrievedTenders, len(tenders))
}

func TestTenderHandler_GetTendersByOrganizationHandler(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewTenderRepository(db)
	service := services.NewTenderService(repo)
	handler := handlers.NewTenderHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/organizations/{organization_id}/tenders", handler.GetTendersByOrganizationHandler).Methods("GET")

	tenders := []*models.Tender{
		{
			Name:            "Test Tender 1",
			Description:     "Test Description 1",
			OrganizationID:  organizationID,
			PublicationDate: time.Now(),
			EndDate:         time.Now().Add(24 * time.Hour),
		},
		{
			Name:            "Test Tender 2",
			Description:     "Test Description 2",
			OrganizationID:  organizationID,
			PublicationDate: time.Now(),
			EndDate:         time.Now().Add(24 * time.Hour),
		},
	}

	for _, tender := range tenders {
		service.Create(tender)
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("/organizations/%d/tenders", organizationID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var retrievedTenders []*models.Tender
	json.NewDecoder(w.Body).Decode(&retrievedTenders)
	assert.Len(t, retrievedTenders, len(tenders))

	req, _ = http.NewRequest("GET", "/organizations/invalid-id/tenders", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTenderHandler_PublishTenderHandler(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewTenderRepository(db)
	service := services.NewTenderService(repo)
	handler := handlers.NewTenderHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/tenders/{id}/publish", handler.PublishTenderHandler).Methods("PUT")

	tender := &models.Tender{
		Name:            "Test Tender",
		Description:     "Test Description",
		OrganizationID:  organizationID,
		PublicationDate: time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
	}
	id, _ := service.Create(tender)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/tenders/%d/publish", id), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("PUT", fmt.Sprintf("/tenders/%d/publish", id), nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	req, _ = http.NewRequest("PUT", "/tenders/invalid-id/publish", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTenderHandler_CloseTenderHandler(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewTenderRepository(db)
	service := services.NewTenderService(repo)
	handler := handlers.NewTenderHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/tenders/{id}/close", handler.CloseTenderHandler).Methods("PUT")

	tender := &models.Tender{
		Name:            "Test Tender",
		Description:     "Test Description",
		OrganizationID:  organizationID,
		PublicationDate: time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
	}
	id, _ := service.Create(tender)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/tenders/%d/close", id), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("PUT", fmt.Sprintf("/tenders/%d/close", id), nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	req, _ = http.NewRequest("PUT", "/tenders/invalid-id/close", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}
*/
