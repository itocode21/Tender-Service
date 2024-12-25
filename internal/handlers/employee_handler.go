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

	// Регистрация пользователя
	id, err := h.service.Register(&user)
	if err != nil {
		// Если ошибка связана с пустым именем пользователя, возвращаем статус 400
		if err.Error() == "username cannot be empty" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Для других ошибок возвращаем статус 500
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Заполнение полей ID и других данных пользователя
	user.Id = id                // Предполагается, что `id` - это ID, который получили от сервиса
	user.CreatedAt = time.Now() // Заполняем дату создания
	user.UpdatedAt = time.Now() // Заполняем дату обновления

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user) // Возвращаем полную информацию о пользователе
}

// GetUser Handler обрабатывает запрос на получение пользователя по ID
func (h *EmployeeHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Извлечение ID из URL
	vars := mux.Vars(r) // Используем mux для извлечения переменных из URL
	idStr := vars["id"] // Сохраняем ID как строку

	// Преобразование ID из строки в int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	// Получение пользователя из сервиса

	user, err := h.service.GetUserByID(id)
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
