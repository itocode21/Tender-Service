package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/itocode21/Tender-Service/internal/handlers"
	"github.com/itocode21/Tender-Service/internal/models"
	"github.com/itocode21/Tender-Service/internal/repository"
	"github.com/itocode21/Tender-Service/internal/services"
	postgresqldb "github.com/itocode21/Tender-Service/postgresql-db"
)

var db *sql.DB

func setup() {
	var err error
	db, err = postgresqldb.InitDB()
	if err != nil {
		panic(err)
	}

	// Очистка таблицы employees перед каждым тестом
	_, err = db.Exec("DELETE FROM employees")
	if err != nil {
		panic(err)
	}
}

func teardown() {
	_ = db.Close()
}

type MockEmployeeRepository struct{}

func (m *MockEmployeeRepository) Create(user *models.User) (int, error) {
	return 0, nil
}
func (m *MockEmployeeRepository) GetByID(id int) (*models.User, error) {
	return nil, errors.New("mock error")
}
func (m *MockEmployeeRepository) GetByUsername(username string) (*models.User, error) {
	return nil, nil
}

func TestRegisterHandler_Success(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewEmployeeRepository(db) // Используйте вашу реализацию репозитория
	service := services.NewEmployeeService(repo) // Передайте репозиторий в сервис
	handler := handlers.NewEmployeeHandler(service)

	user := models.User{
		Username:  "avitoUser",
		FirstName: "AvitoName",
		LastName:  "AvitoLastName",
	}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.RegisterHandler(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status %v, got %v", http.StatusCreated, res.StatusCode)
	}

	var responseUser models.User
	json.NewDecoder(res.Body).Decode(&responseUser)
	if responseUser.Username != user.Username || responseUser.FirstName != user.FirstName || responseUser.LastName != user.LastName {
		t.Errorf("expected user to be %v, got %v", user, responseUser)
	}
}

func TestRegisterHandler_BadRequest(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewEmployeeRepository(db)
	service := services.NewEmployeeService(repo)
	handler := handlers.NewEmployeeHandler(service)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer([]byte("invalid json")))
	w := httptest.NewRecorder()

	handler.RegisterHandler(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %v, got %v", http.StatusBadRequest, res.StatusCode)
	}
}

func TestRegisterHandler_InternalServerError(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewEmployeeRepository(db)
	service := services.NewEmployeeService(repo)
	handler := handlers.NewEmployeeHandler(service)

	// Создайте пользователя с недопустимыми данными, чтобы вызвать ошибку
	user := models.User{ // Например, пустое имя пользователя
		Username:  "",
		FirstName: "John",
		LastName:  "Doe",
	}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.RegisterHandler(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status %v, got %v", http.StatusInternalServerError, res.StatusCode)
	}
}

func TestGetUserHandler_Success(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewEmployeeRepository(db)
	service := services.NewEmployeeService(repo)
	handler := handlers.NewEmployeeHandler(service)

	user := models.User{
		Username:  "AvitoUsername",
		FirstName: "AvitoName",
		LastName:  "AvitoLastName",
	}
	registeredID, _ := service.Register(&user) // Получаем ID

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/user/%d", registeredID), nil) // Используем ID в URL
	req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", registeredID)})
	w := httptest.NewRecorder()

	handler.GetUserHandler(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, res.StatusCode)
	}

	var responseUser models.User
	json.NewDecoder(res.Body).Decode(&responseUser)
	if responseUser.Username != user.Username || responseUser.FirstName != user.FirstName || responseUser.LastName != user.LastName {
		t.Errorf("expected user to be %v, got %v", user, responseUser)
	}
}

func TestGetUserHandler_NotFound(t *testing.T) {
	setup()
	defer teardown()

	repo := repository.NewEmployeeRepository(db)
	service := services.NewEmployeeService(repo)
	handler := handlers.NewEmployeeHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/user/999", nil) // ID, который не существует
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	w := httptest.NewRecorder()

	handler.GetUserHandler(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected status %v, got %v", http.StatusNotFound, res.StatusCode)
	}
}

func TestGetUserHandler_InternalServerError(t *testing.T) {
	setup()
	defer teardown()

	mockRepo := &MockEmployeeRepository{}
	service := services.NewEmployeeService(mockRepo)
	handler := handlers.NewEmployeeHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/user/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.GetUserHandler(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status %v, got %v", http.StatusInternalServerError, res.StatusCode)
	}
}
