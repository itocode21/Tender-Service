package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/itocode21/Tender-Service/internal/models"
	"github.com/itocode21/Tender-Service/internal/services"
)

type EmployeeHandler struct {
	service services.EmployeeService
}

func NewEmployeeHandler(service services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: service}
}

func (h *EmployeeHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Регистрация пользователя
	id, err := h.service.Create(&user)
	if err != nil {
		if err.Error() == "username cannot be empty" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	user.Id = id // Предполагается, что id это ID, который мы получили от сервиса
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// тут настроил пирамид,  но как более чисто реализовать(чтобы работало при этом) я не знаю.
func (h *EmployeeHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Извлечение id из URL
	vars := mux.Vars(r) // Используем mux для извлечения переменных из URL
	idStr := vars["id"] // Сохраняем id как строку

	// Преобразование ID из строки в int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	// Получение пользователя из сервиса
	user, err := h.service.GetByID(id)
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
