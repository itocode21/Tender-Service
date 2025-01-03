package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/itocode21/Tender-Service/internal/models"
	"github.com/itocode21/Tender-Service/internal/services"
)

type OrganizationHandler struct {
	service services.OrganizationService
}

func NewOrganizationHandler(service services.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{service: service}
}

func (h *OrganizationHandler) CreateOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	var org models.Organization
	if err := json.NewDecoder(r.Body).Decode(&org); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(&org)
	if err != nil {
		if err.Error() == "имя организации не может быть пустым" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	org.Id = int64(id)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(org)
}

func (h *OrganizationHandler) GetOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	org, err := h.service.GetByID(id)
	if err != nil {
		if err.Error() == "organization not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(org)
}

func (h *OrganizationHandler) UpdateOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var org models.Organization
	if err := json.NewDecoder(r.Body).Decode(&org); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	org.Id = int64(id)

	err = h.service.Update(&org)
	if err != nil {
		if err.Error() == "имя организации не может быть пустым" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(org)
}

func (h *OrganizationHandler) DeleteOrganizationHandler(w http.ResponseWriter, r *http.Request) {
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

func (h *OrganizationHandler) ListOrganizationsHandler(w http.ResponseWriter, r *http.Request) {
	orgs, err := h.service.List()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orgs)
}

func (h *OrganizationHandler) AddResponsibleHandler(w http.ResponseWriter, r *http.Request) {
	var orgResp models.OrganizationResponsible
	if err := json.NewDecoder(r.Body).Decode(&orgResp); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	id, err := h.service.AddResponsible(&orgResp)
	if err != nil {
		if err.Error() == "организация и пользователь не могут быть nil" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	orgResp.Id = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(orgResp)
}

func (h *OrganizationHandler) RemoveResponsibleHandler(w http.ResponseWriter, r *http.Request) {
	var orgResp models.OrganizationResponsible
	if err := json.NewDecoder(r.Body).Decode(&orgResp); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err := h.service.RemoveResponsible(&orgResp)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
