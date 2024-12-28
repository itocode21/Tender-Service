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

// parseIDFromRequest извлекает и парсит ID из URL
func (h *ProposalHandler) parseIDFromRequest(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// parseTenderIDFromRequest извлекает и парсит tender_id из URL
func (h *ProposalHandler) parseTenderIDFromRequest(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	idStr := vars["tender_id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// sendJSONResponse отправляет JSON ответ
func (h *ProposalHandler) sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// CreateProposalHandler обрабатывает запрос на создание нового предложения
func (h *ProposalHandler) CreateProposalHandler(w http.ResponseWriter, r *http.Request) {
	var proposal models.Proposal
	if err := json.NewDecoder(r.Body).Decode(&proposal); err != nil {
		h.sendJSONResponse(w, map[string]string{"error": "Invalid input"}, http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(&proposal)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		var statusCode int

		switch err.Error() {
		case "proposal not found":
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusBadRequest
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	proposal.ID = id
	h.sendJSONResponse(w, proposal, http.StatusCreated)
}

// GetProposalHandler обрабатывает запрос на получение предложения по ID
func (h *ProposalHandler) GetProposalHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromRequest(r)
	if err != nil {
		h.sendJSONResponse(w, map[string]string{"error": "Invalid ID format"}, http.StatusBadRequest)
		return
	}

	proposal, err := h.service.GetByID(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		var statusCode int
		switch err.Error() {
		case "proposal not found":
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusInternalServerError
		}
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	h.sendJSONResponse(w, proposal, http.StatusOK)
}

// UpdateProposalHandler обрабатывает запрос на обновление предложения
func (h *ProposalHandler) UpdateProposalHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromRequest(r)
	if err != nil {
		h.sendJSONResponse(w, map[string]string{"error": "Invalid ID format"}, http.StatusBadRequest)
		return
	}

	var proposal models.Proposal
	if err := json.NewDecoder(r.Body).Decode(&proposal); err != nil {
		h.sendJSONResponse(w, map[string]string{"error": "Invalid input"}, http.StatusBadRequest)
		return
	}
	proposal.ID = id

	err = h.service.Update(&proposal)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		var statusCode int

		switch err.Error() {
		case "proposal not found":
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusBadRequest
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	h.sendJSONResponse(w, proposal, http.StatusOK)
}

// DeleteProposalHandler обрабатывает запрос на удаление предложения
func (h *ProposalHandler) DeleteProposalHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromRequest(r)
	if err != nil {
		h.sendJSONResponse(w, map[string]string{"error": "Invalid ID format"}, http.StatusBadRequest)
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		var statusCode int

		switch err.Error() {
		case "proposal not found":
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusInternalServerError
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListProposalsHandler обрабатывает запрос на получение списка предложений
func (h *ProposalHandler) ListProposalsHandler(w http.ResponseWriter, r *http.Request) {
	proposals, err := h.service.List()
	if err != nil {
		return
	}

	h.sendJSONResponse(w, proposals, http.StatusOK)
}

// GetProposalsByTenderHandler обрабатывает запрос на получение списка предложений по ID тендера
func (h *ProposalHandler) GetProposalsByTenderHandler(w http.ResponseWriter, r *http.Request) {
	tenderID, err := h.parseTenderIDFromRequest(r)
	if err != nil {
		h.sendJSONResponse(w, map[string]string{"error": "Invalid ID format"}, http.StatusBadRequest)
		return
	}

	proposals, err := h.service.GetByTenderID(tenderID)
	if err != nil {
		return
	}
	h.sendJSONResponse(w, proposals, http.StatusOK)
}

// PublishProposalHandler устанавливает статус предложения как "published"
func (h *ProposalHandler) PublishProposalHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromRequest(r)
	if err != nil {
		h.sendJSONResponse(w, map[string]string{"error": "Invalid ID format"}, http.StatusBadRequest)
		return
	}

	err = h.service.Publish(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		var statusCode int

		switch err.Error() {
		case "proposal not found":
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusBadRequest
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	h.sendJSONResponse(w, map[string]string{"message": "proposal published"}, http.StatusOK)
}

// AcceptProposalHandler устанавливает статус предложения как "accepted"
func (h *ProposalHandler) AcceptProposalHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromRequest(r)
	if err != nil {
		h.sendJSONResponse(w, map[string]string{"error": "Invalid ID format"}, http.StatusBadRequest)
		return
	}

	err = h.service.Accept(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		var statusCode int

		switch err.Error() {
		case "proposal not found":
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusBadRequest
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	h.sendJSONResponse(w, map[string]string{"message": "proposal accepted"}, http.StatusOK)

}

// RejectProposalHandler устанавливает статус предложения как "rejected"
func (h *ProposalHandler) RejectProposalHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromRequest(r)
	if err != nil {
		h.sendJSONResponse(w, map[string]string{"error": "Invalid ID format"}, http.StatusBadRequest)
		return
	}

	err = h.service.Reject(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		var statusCode int

		switch err.Error() {
		case "proposal not found":
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusBadRequest
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	h.sendJSONResponse(w, map[string]string{"message": "proposal rejected"}, http.StatusOK)
}

// CancelProposalHandler устанавливает статус предложения как "cancelled"
func (h *ProposalHandler) CancelProposalHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromRequest(r)
	if err != nil {
		h.sendJSONResponse(w, map[string]string{"error": "Invalid ID format"}, http.StatusBadRequest)
		return
	}

	err = h.service.Cancel(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		var statusCode int
		switch err.Error() {
		case "proposal not found":
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusBadRequest
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	h.sendJSONResponse(w, map[string]string{"message": "proposal cancelled"}, http.StatusOK)
}
