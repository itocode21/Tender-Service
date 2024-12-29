package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/itocode21/Tender-Service/internal/models"
	"github.com/itocode21/Tender-Service/internal/services"
)

type TenderHandler struct {
	service services.TenderService
}

func NewTenderHandler(service services.TenderService) *TenderHandler {
	return &TenderHandler{service: service}
}

func (h *TenderHandler) CreateTenderHandler(w http.ResponseWriter, r *http.Request) {
	var tender models.Tender
	if err := json.NewDecoder(r.Body).Decode(&tender); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(&tender)
	if err != nil {
		if errors.Is(err, errors.New("название тендера не может быть пустым")) || errors.Is(err, errors.New("идентификатор организации не может быть пустым")) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Bad Requeqst", http.StatusBadRequest)
		return
	}

	tender.ID = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tender)
}

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
		if errors.Is(err, errors.New("tender not found")) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Not Found 404", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tender)
}

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
		if errors.Is(err, errors.New("имя тендера не может быть пустым")) || errors.Is(err, errors.New("tender not found")) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Bad Requeqst", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tender)
}

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

func (h *TenderHandler) GetTendersByOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["organization_id"]

	organizationID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	tenders, err := h.service.ListByOrganizationID(organizationID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tenders)

}

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
		if errors.Is(err, errors.New("тендер уже опубликован")) || errors.Is(err, errors.New("tender not found")) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

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
		if errors.Is(err, errors.New("тендер уже закрыт")) || errors.Is(err, errors.New("невозможно закрыть отмененный тендер")) || errors.Is(err, errors.New("tender not found")) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TenderHandler) CancelTenderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = h.service.Close(id)
	if err != nil {
		if errors.Is(err, errors.New("тендер уже отменен")) || errors.Is(err, errors.New("невозможно отменить закрытый тендер")) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
