package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/matthieukhl/align-back/internal/domain"
	"github.com/rs/zerolog/log"
)

type PackageHandler struct {
	service domain.PackageService
}

// NewPackageHandler creates a new package handler.
func NewPackageHandler(service domain.PackageService) *PackageHandler {
	return &PackageHandler{
		service: service,
	}
}

// r.Route("/packages", func(r chi.Router) {
// 	r.Get("/", packageHandler.GetAll)
// 	r.Post("/", packageHandler.Create)
// 	r.Get("/{id}", packageHandler.GetByID)
// 	r.Put("/{id}", packageHandler.Update)
// 	r.Delete("/{id}", packageHandler.Delete)
// })

// GetAll handles GET /api/packages
func (h *PackageHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	packages, err := h.service.GetAll()
	if err != nil {
		log.Error().Err(err).Msg("failed to get all clients")
		http.Error(w, "failed to get clients", http.StatusInternalServerError)
		return
	}

	respondwithJSON(w, http.StatusOK, packages)
}

// GetByID handles GET /api/packages/{id}
func (h *PackageHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing package ID", http.StatusBadRequest)
		return
	}

	pkg, err := h.service.GetByID(id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("failed to get package by ID")
		http.Error(w, "Failed to get package", http.StatusInternalServerError)
		return
	}

	if pkg == nil {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	respondwithJSON(w, http.StatusOK, pkg)
}

// Create handles POST /api/packages
func (h *PackageHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input domain.PackageInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if input.Name == "" || input.NumberOfSessions == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	err := h.service.Create(input)
	if err != nil {
		log.Error().Err(err).Interface("input", input).Msg("failed to create package")
		if err.Error() == "package already exists" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create package", http.StatusInternalServerError)
		return
	}

	respondwithJSON(w, http.StatusCreated, input)
}

// Update handles PUT /api/packages/{id}
func (h *PackageHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing package ID", http.StatusBadRequest)
		return
	}

	var input domain.PackageInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if input.Name == "" || input.NumberOfSessions == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	pkg, err := h.service.Update(id, input)
	if err != nil {
		log.Error().Err(err).Str("id", id).Interface("input", input).Msg("failed to update package")
		if err.Error() == "package not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err.Error() == "package name already exists" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Failed to update package", http.StatusInternalServerError)
		return
	}

	respondwithJSON(w, http.StatusOK, pkg)
}

// Delete handles DELETE /api/packages/{id}
func (h *PackageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "Missing package ID", http.StatusBadRequest)
		return
	}

	err := h.service.Delete(id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("failed to delete client")
		if err.Error() == "package not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete package", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
