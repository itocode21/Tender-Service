package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/itocode21/Tender-Service/internal/models"
	"github.com/itocode21/Tender-Service/internal/services"
)

// ProposalHandler структура для обработки запросов к предложениям
type ProposalHandler struct {
	service services.ProposalService
}

// NewProposalHandler создает новый экземпляр ProposalHandler
func NewProposalHandler(service services.ProposalService) *ProposalHandler {
	return &ProposalHandler{service: service}
}

// CreateProposalHandler обрабатывает запрос на создание нового предложения
func (h *ProposalHandler) CreateProposalHandler(w http.ResponseWriter, r *http.Request) {
	var proposal models.Proposal
	if err := json.NewDecoder(r.Body).Decode(&proposal); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(&proposal)
	if err != nil {
		if err.Error() == "описание предложения не может быть пустым" || err.Error() == "необходимо указать тендер" || err.Error() == "необходимо указать организацию" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	proposal.ID = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(proposal)
}

// GetProposalHandler обрабатывает запрос на получение предложения по ID
func (h *ProposalHandler) GetProposalHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	proposal, err := h.service.GetByID(id)
	if err != nil {
		if err.Error() == "proposal not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(proposal)
}

// UpdateProposalHandler обрабатывает запрос на обновление предложения
func (h *ProposalHandler) UpdateProposalHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var proposal models.Proposal
	if err := json.NewDecoder(r.Body).Decode(&proposal); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	proposal.ID = id
	err = h.service.Update(&proposal)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(proposal)
}

// DeleteProposalHandler обрабатывает запрос на удаление предложения
func (h *ProposalHandler) DeleteProposalHandler(w http.ResponseWriter, r *http.Request) {
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

// ListProposalsHandler обрабатывает запрос на получение списка предложений
func (h *ProposalHandler) ListProposalsHandler(w http.ResponseWriter, r *http.Request) {
	proposals, err := h.service.List()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(proposals)
}

// GetProposalsByTenderHandler обрабатывает запрос на получение списка предложений по ID тендера
func (h *ProposalHandler) GetProposalsByTenderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["tender_id"]

	tenderID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	proposals, err := h.service.GetByTenderID(tenderID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(proposals)
}

// PublishProposalHandler устанавливает статус предложения как "published"
func (h *ProposalHandler) PublishProposalHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = h.service.Publish(id)
	if err != nil {
		if err.Error() == "предложение уже опубликовано" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// AcceptProposalHandler устанавливает статус предложения как "accepted"
func (h *ProposalHandler) AcceptProposalHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = h.service.Accept(id)
	if err != nil {
		if err.Error() == "предложение уже принято" || err.Error() == "невозможно принять отклоненное предложение" || err.Error() == "невозможно принять отмененное предложение" || err.Error() == "невозможно принять предложение для закрытого тендера" || err.Error() == "невозможно принять предложение для отмененного тендера" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// RejectProposalHandler устанавливает статус предложения как "rejected"
func (h *ProposalHandler) RejectProposalHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = h.service.Reject(id)
	if err != nil {
		if err.Error() == "предложение уже отклонено" || err.Error() == "невозможно отклонить принятое предложение" || err.Error() == "невозможно отклонить отмененное предложение" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// CancelProposalHandler устанавливает статус предложения как "cancelled"
func (h *ProposalHandler) CancelProposalHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = h.service.Cancel(id)
	if err != nil {
		if err.Error() == "предложение уже отменено" || err.Error() == "невозможно отменить принятое предложение" || err.Error() == "невозможно отменить отклоненное предложение" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
