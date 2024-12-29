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

	"github.com/gorilla/mux"
	"github.com/itocode21/Tender-Service/internal/handlers"
	"github.com/itocode21/Tender-Service/internal/models"
	"github.com/itocode21/Tender-Service/internal/repository"
	"github.com/itocode21/Tender-Service/internal/services"
	postgresqldb "github.com/itocode21/Tender-Service/postgresql-db"

	"database/sql"

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

func TestOrganizationHandler_CreateOrganizationHandler(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewOrganizationRepository(db)
	service := services.NewOrganizationService(repo)
	handler := handlers.NewOrganizationHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/organizations", handler.CreateOrganizationHandler).Methods("POST")

	org := &models.Organization{
		Name:        "Test Organization",
		Description: "Test Description",
		Type:        models.OrganizationTypeLLC,
	}

	jsonValue, _ := json.Marshal(org)
	req, _ := http.NewRequest("POST", "/organizations", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var createdOrg models.Organization
	json.NewDecoder(w.Body).Decode(&createdOrg)
	assert.Equal(t, org.Name, createdOrg.Name)
	assert.Equal(t, org.Description, createdOrg.Description)
	assert.Equal(t, org.Type, createdOrg.Type)

	//test for invalid input
	invalidOrg := &models.Organization{
		Name:        "",
		Description: "Test Description",
		Type:        models.OrganizationTypeLLC,
	}
	jsonValue, _ = json.Marshal(invalidOrg)
	req, _ = http.NewRequest("POST", "/organizations", bytes.NewBuffer(jsonValue))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestOrganizationHandler_GetOrganizationHandler(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewOrganizationRepository(db)
	service := services.NewOrganizationService(repo)
	handler := handlers.NewOrganizationHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/organizations/{id}", handler.GetOrganizationHandler).Methods("GET")

	org := &models.Organization{
		Name:        "Test Organization",
		Description: "Test Description",
		Type:        models.OrganizationTypeLLC,
	}
	id, _ := repo.Create(org)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/organizations/%d", id), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var retrievedOrg models.Organization
	json.NewDecoder(w.Body).Decode(&retrievedOrg)
	assert.Equal(t, org.Name, retrievedOrg.Name)
	assert.Equal(t, org.Description, retrievedOrg.Description)
	assert.Equal(t, org.Type, retrievedOrg.Type)

	req, _ = http.NewRequest("GET", fmt.Sprintf("/organizations/%d", id+1), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	req, _ = http.NewRequest("GET", "/organizations/invalid-id", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestOrganizationHandler_UpdateOrganizationHandler(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewOrganizationRepository(db)
	service := services.NewOrganizationService(repo)
	handler := handlers.NewOrganizationHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/organizations/{id}", handler.UpdateOrganizationHandler).Methods("PUT")

	org := &models.Organization{
		Name:        "Test Organization",
		Description: "Test Description",
		Type:        models.OrganizationTypeLLC,
	}
	id, _ := repo.Create(org)

	updatedOrg := &models.Organization{
		Name:        "Updated Test Organization",
		Description: "Updated Test Description",
		Type:        models.OrganizationTypeJSC,
	}
	jsonValue, _ := json.Marshal(updatedOrg)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/organizations/%d", id), bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var retrievedOrg models.Organization
	json.NewDecoder(w.Body).Decode(&retrievedOrg)
	assert.Equal(t, updatedOrg.Name, retrievedOrg.Name)
	assert.Equal(t, updatedOrg.Description, retrievedOrg.Description)
	assert.Equal(t, updatedOrg.Type, retrievedOrg.Type)

	invalidOrg := &models.Organization{
		Name:        "",
		Description: "Updated Test Description",
		Type:        models.OrganizationTypeJSC,
	}
	jsonValue, _ = json.Marshal(invalidOrg)
	req, _ = http.NewRequest("PUT", fmt.Sprintf("/organizations/%d", id), bytes.NewBuffer(jsonValue))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	jsonValue, _ = json.Marshal(updatedOrg)
	req, _ = http.NewRequest("PUT", "/organizations/invalid-id", bytes.NewBuffer(jsonValue))
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestOrganizationHandler_DeleteOrganizationHandler(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewOrganizationRepository(db)
	service := services.NewOrganizationService(repo)
	handler := handlers.NewOrganizationHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/organizations/{id}", handler.DeleteOrganizationHandler).Methods("DELETE")

	org := &models.Organization{
		Name:        "Test Organization",
		Description: "Test Description",
		Type:        models.OrganizationTypeLLC,
	}
	id, _ := repo.Create(org)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/organizations/%d", id), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	req, _ = http.NewRequest("DELETE", "/organizations/invalid-id", nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestOrganizationHandler_ListOrganizationsHandler(t *testing.T) {
	setup()
	defer teardown()
	repo := repository.NewOrganizationRepository(db)
	service := services.NewOrganizationService(repo)
	handler := handlers.NewOrganizationHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/organizations", handler.ListOrganizationsHandler).Methods("GET")

	orgs := []*models.Organization{
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

	for _, org := range orgs {
		repo.Create(org)
	}

	req, _ := http.NewRequest("GET", "/organizations", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var retrievedOrgs []*models.Organization
	json.NewDecoder(w.Body).Decode(&retrievedOrgs)
	assert.Len(t, retrievedOrgs, len(orgs))
}
*/
