package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/itocode21/Tender-Service/internal/models"
	"github.com/itocode21/Tender-Service/internal/services"
)

type ProposalHandler struct {
	service services.ProposalService
}

func NewProposalHandler(service services.ProposalService) *ProposalHandler {
	return &ProposalHandler{service: service}
}

func (h *ProposalHandler) parseIDFromRequest(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (h *ProposalHandler) parseTenderIDFromRequest(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	idStr := vars["tender_id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (h *ProposalHandler) handleServiceError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	var statusCode int

	switch err.Error() {
	case "proposal not found", "tender not found":
		statusCode = http.StatusNotFound
	default:
		statusCode = http.StatusBadRequest
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

func (h *ProposalHandler) sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *ProposalHandler) CreateProposalHandler(w http.ResponseWriter, r *http.Request) {
	var proposal models.Proposal
	if err := json.NewDecoder(r.Body).Decode(&proposal); err != nil {
		h.sendJSONResponse(w, map[string]string{"error": "Invalid input"}, http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(&proposal)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	proposal.ID = id
	h.sendJSONResponse(w, proposal, http.StatusCreated)
}

func (h *ProposalHandler) GetProposalHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromRequest(r)
	if err != nil {
		h.sendJSONResponse(w, map[string]string{"error": "Invalid ID format"}, http.StatusBadRequest)
		return
	}

	proposal, err := h.service.GetByID(id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.sendJSONResponse(w, proposal, http.StatusOK)
}

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
		h.handleServiceError(w, err)
		return
	}

	h.sendJSONResponse(w, proposal, http.StatusOK)
}

func (h *ProposalHandler) DeleteProposalHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromRequest(r)
	if err != nil {
		h.sendJSONResponse(w, map[string]string{"error": "Invalid ID format"}, http.StatusBadRequest)
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProposalHandler) ListProposalsHandler(w http.ResponseWriter, r *http.Request) {
	proposals, err := h.service.List()
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.sendJSONResponse(w, proposals, http.StatusOK)
}

func (h *ProposalHandler) GetProposalsByTenderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tenderID, err := h.parseTenderIDFromRequest(r)
	if err != nil {
		h.sendJSONResponse(w, map[string]string{"error": "Invalid ID format"}, http.StatusBadRequest)
		return
	}
	log.Printf("handler: get proposals by tender %v", tenderID)
	proposals, err := h.service.GetByTenderID(tenderID)
	if err != nil {
		log.Printf("handler: got error from service while get proposals by tender %v, error: %v", tenderID, err)
		h.handleServiceError(w, err)
		return
	}
	log.Printf("handler: get proposals by tender %v successfully %v", tenderID, proposals)
	h.sendJSONResponse(w, proposals, http.StatusOK)
}

func (h *ProposalHandler) PublishProposalHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromRequest(r)
	if err != nil {
		h.sendJSONResponse(w, map[string]string{"error": "Invalid ID format"}, http.StatusBadRequest)
		return
	}

	err = h.service.Publish(id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.sendJSONResponse(w, map[string]string{"message": "proposal published"}, http.StatusOK)
}

func (h *ProposalHandler) AcceptProposalHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromRequest(r)
	if err != nil {
		h.sendJSONResponse(w, map[string]string{"error": "Invalid ID format"}, http.StatusBadRequest)
		return
	}

	err = h.service.Accept(id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.sendJSONResponse(w, map[string]string{"message": "proposal accepted"}, http.StatusOK)

}

func (h *ProposalHandler) RejectProposalHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromRequest(r)
	if err != nil {
		h.sendJSONResponse(w, map[string]string{"error": "Invalid ID format"}, http.StatusBadRequest)
		return
	}

	err = h.service.Reject(id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.sendJSONResponse(w, map[string]string{"message": "proposal rejected"}, http.StatusOK)
}

func (h *ProposalHandler) CancelProposalHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseIDFromRequest(r)
	if err != nil {
		h.sendJSONResponse(w, map[string]string{"error": "Invalid ID format"}, http.StatusBadRequest)
		return
	}

	err = h.service.Cancel(id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.sendJSONResponse(w, map[string]string{"message": "proposal cancelled"}, http.StatusOK)
}
