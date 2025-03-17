package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/matthieukhl/align-back/internal/domain"
	"github.com/rs/zerolog/log"
)

type ClientHandler struct {
	service domain.ClientService
}

// NewClientHandler creates a new client handler
func NewClientHandler(service domain.ClientService) *ClientHandler {
	return &ClientHandler{
		service: service,
	}
}

// GetAll handles GET /api/clients
func (h *ClientHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	clients, err := h.service.GetAll()
	if err != nil {
		log.Error().Err(err).Msg("failed to get all clients")
		http.Error(w, "Failed to get clients", http.StatusInternalServerError)
		return
	}

	respondwithJSON(w, http.StatusOK, clients)
}

// GetByID handles /api/clients/{id}
func (h *ClientHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing client ID", http.StatusBadRequest)
		return
	}

	client, err := h.service.GetByID(id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("failed to get client by ID")
		http.Error(w, "Failed to get client", http.StatusInternalServerError)
		return
	}

	if client == nil {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	respondwithJSON(w, http.StatusOK, client)
}

// Create handles POST /api/clients
func (h *ClientHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input domain.ClientInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if input.FirstName == "" || input.LastName == "" || input.Email == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	err := h.service.Create(input)
	if err != nil {
		log.Error().Err(err).Interface("input", input).Msg("failed to create client")
		if err.Error() == "email is already in use" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create client", http.StatusInternalServerError)
		return
	}

	respondwithJSON(w, http.StatusCreated, input)
}

// Update handles PUT /api/clients/{id}
func (h *ClientHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing Client ID", http.StatusBadRequest)
		return
	}

	var input domain.ClientInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if input.FirstName == "" || input.LastName == "" || input.Email == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	client, err := h.service.Update(id, input)
	if err != nil {
		log.Error().Err(err).Str("id", id).Interface("input", input).Msg("failed to update client")
		if err.Error() == "client not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err.Error() == "email is already in use" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Failed to update client", http.StatusInternalServerError)
		return
	}

	respondwithJSON(w, http.StatusOK, client)
}

// Delete handles DELETE /api/clients/{id}
func (h *ClientHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing client ID", http.StatusBadRequest)
		return
	}

	err := h.service.Delete(id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("failed to delete client")
		if err.Error() == "client not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete client", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetLowGroupCredits handles GET /api/clients/low-group-credits
func (h *ClientHandler) GetLowGroupCredits(w http.ResponseWriter, r *http.Request) {
	threshold := 1 // Default threshold
	clients, err := h.service.GetLowGroupCredits(threshold)
	if err != nil {
		log.Error().Err(err).Msg("failed to get clients with low group credits")
		http.Error(w, "Failed to get clients with low group credits", http.StatusInternalServerError)
		return
	}

	respondwithJSON(w, http.StatusOK, clients)
}

// GetLowPrivate handles GET /api/clients/low-private-credits
func (h *ClientHandler) GetLowPrivateCredits(w http.ResponseWriter, r *http.Request) {
	threshold := 1 // Default threshold
	clients, err := h.service.GetLowPrivateCredits(threshold)
	if err != nil {
		log.Error().Err(err).Msg("failed to get clients with low private credits")
		http.Error(w, "Failed to get clients with low private credits", http.StatusInternalServerError)
		return
	}

	respondwithJSON(w, http.StatusOK, clients)
}

// Helper to respond with JSON
func respondwithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}
