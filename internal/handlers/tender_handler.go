package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/itocode21/Tender-Service/internal/models"
	"github.com/itocode21/Tender-Service/internal/services"
)

// TenderHandler структура для обработки запросов к тендерам
type TenderHandler struct {
	service services.TenderService
}

// NewTenderHandler создает новый экземпляр TenderHandler
func NewTenderHandler(service services.TenderService) *TenderHandler {
	return &TenderHandler{service: service}
}

// CreateTenderHandler обрабатывает запрос на создание нового тендера
func (h *TenderHandler) CreateTenderHandler(w http.ResponseWriter, r *http.Request) {
	var tender models.Tender
	if err := json.NewDecoder(r.Body).Decode(&tender); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(&tender)
	if err != nil {
		if err.Error() == "название тендера не может быть пустым" || err.Error() == "необходимо указать организацию" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tender.ID = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tender)
}

// GetTenderHandler обрабатывает запрос на получение тендера по ID
func (h *TenderHandler) GetTenderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	tender, err := h.service.GetByID(id)
	if err != nil {
		if err.Error() == "tender not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tender)
}

// UpdateTenderHandler обрабатывает запрос на обновление тендера
func (h *TenderHandler) UpdateTenderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var tender models.Tender
	if err := json.NewDecoder(r.Body).Decode(&tender); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	tender.ID = id

	err = h.service.Update(&tender)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tender)
}

// DeleteTenderHandler обрабатывает запрос на удаление тендера
func (h *TenderHandler) DeleteTenderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListTendersHandler обрабатывает запрос на получение списка тендеров
func (h *TenderHandler) ListTendersHandler(w http.ResponseWriter, r *http.Request) {
	tenders, err := h.service.List()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tenders)
}

// GetTendersByOrganizationHandler обрабатывает запрос на получение списка тендеров по ID организации
func (h *TenderHandler) GetTendersByOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["organization_id"]

	organizationID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	tenders, err := h.service.GetByOrganizationID(organizationID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tenders)

}

// PublishTenderHandler устанавливает статус тендера как "published"
func (h *TenderHandler) PublishTenderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = h.service.Publish(id)
	if err != nil {
		if err.Error() == "тендер уже опубликован" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// CloseTenderHandler устанавливает статус тендера как "closed"
func (h *TenderHandler) CloseTenderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = h.service.Close(id)
	if err != nil {
		if err.Error() == "тендер уже закрыт" || err.Error() == "невозможно закрыть отмененный тендер" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// CancelTenderHandler устанавливает статус тендера как "cancelled"
func (h *TenderHandler) CancelTenderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = h.service.Cancel(id)
	if err != nil {
		if err.Error() == "тендер уже отменен" || err.Error() == "невозможно отменить закрытый тендер" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
