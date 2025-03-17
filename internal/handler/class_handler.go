package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/matthieukhl/align-back/internal/domain"
	"github.com/rs/zerolog/log"
)

type ClassHandler struct {
	service domain.ClassService
}

// NewClassHandler creates a new class handler
func NewClassHandler(service domain.ClassService) *ClassHandler {
	return &ClassHandler{
		service: service,
	}
}

// GetAll handles GET /api/classes
func (h *ClassHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	classes, err := h.service.GetAll()
	if err != nil {
		log.Error().Err(err).Msg("failed to get all classes")
		http.Error(w, "Failed to get classes", http.StatusInternalServerError)
		return
	}

	respondwithJSON(w, http.StatusOK, classes)
}

// GetByID handles GET /api/classes/{id}
func (h *ClassHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "Missing class ID", http.StatusBadRequest)
		return
	}

	class, err := h.service.GetByID(id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("failed to get class by ID")
		http.Error(w, "Failed to get class", http.StatusInternalServerError)
		return
	}

	if class == nil {
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}

	respondwithJSON(w, http.StatusOK, class)
}

// Create handles POST /api/classes
func (h *ClassHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input domain.ClassInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if input.Name == "" || input.Location == "" || input.Type == "" || input.Equipment == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
	}

	err := h.service.Create(input)
	if err != nil {
		log.Error().Err(err).Interface("input", input).Msg("failed to create class")
		if err.Error() == "class name is already in use" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create class", http.StatusInternalServerError)
		return
	}

	respondwithJSON(w, http.StatusCreated, input)

}

// Update handles PUT /api/classes/{id}
func (h *ClassHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing class ID", http.StatusBadRequest)
		return
	}

	var input domain.ClassInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if input.Name == "" || input.Location == "" || input.Type == "" || input.Equipment == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
	}

	class, err := h.service.Update(id, input)
	if err != nil {
		log.Error().Err(err).Str("id", id).Interface("input", input).Msg("failed to update class")
		if err.Error() == "class not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err.Error() == "class name is already in use" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Failed to update class", http.StatusInternalServerError)
		return
	}

	respondwithJSON(w, http.StatusOK, class)
}

// Delete handles DELETE /api/classes/{id}
func (h *ClassHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing class ID", http.StatusBadRequest)
		return
	}

	err := h.service.Delete(id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("failed to delete class")
		if err.Error() == "class not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete class", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
