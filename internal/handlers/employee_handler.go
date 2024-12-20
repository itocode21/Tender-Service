package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/itocode21/Tender-Service/internal/models"
	"github.com/itocode21/Tender-Service/internal/services"
)

type EmployeeHandler struct {
	service services.EmployeeService
}

// NewEmployeeHandler создает новый экземпляр EmployeeHandler
func NewEmployeeHandler(service services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: service}
}

// RegisterHandler обрабатывает регистрацию пользователя
func (h *EmployeeHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.service.Register(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUser Handler обрабатывает запрос на получение пользователя по ID
func (h *EmployeeHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Извлечение ID из URL
	vars := mux.Vars(r) // Используем mux для извлечения переменных из URL
	id := vars["id"]

	// Получение пользователя из сервиса
	user, err := h.service.GetUserByID(id) // Используем сервис для получения пользователя
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Установка заголовка Content-Type и отправка ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
