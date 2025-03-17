package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/matthieukhl/align-back/internal/domain"
	"github.com/rs/zerolog/log"
)

type BillingHandler struct {
	service domain.BillingService
}

// NewBillingHandler creates a new billing handler
func NewBillingHandler(service domain.BillingService) *BillingHandler {
	return &BillingHandler{
		service: service,
	}
}

// GetAll handles GET /api/billings
func (h *BillingHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	billings, err := h.service.GetAll()
	if err != nil {
		log.Error().Err(err).Msg("failed to get billings")
		http.Error(w, "Failed to get billings", http.StatusInternalServerError)
		return
	}

	respondwithJSON(w, http.StatusOK, billings)
}

// GetById handles GET /api/billings/{id}
func (h *BillingHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "Missing billing ID", http.StatusBadRequest)
		return
	}

	billing, err := h.service.GetByID(id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("failed to get billing by ID")
		http.Error(w, "Failed to get billing by ID", http.StatusInternalServerError)
		return
	}

	if billing == nil {
		http.Error(w, "Billing not found", http.StatusNotFound)
		return
	}

	respondwithJSON(w, http.StatusOK, billing)
}

// Create handles POST /api/billings
func (h *BillingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input domain.BillingInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if input.ClientID == "" || input.PackageID == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	err := h.service.Create(input)
	if err != nil {
		log.Error().Err(err).Interface("input", input).Msg("failed to create billing")
		http.Error(w, "Failed to create billing", http.StatusInternalServerError)
		return
	}

	respondwithJSON(w, http.StatusCreated, input)
}

// GetRecent handles /api/billings/recent
func (h *BillingHandler) GetRecent(w http.ResponseWriter, r *http.Request) {
	limit := 10

	billings, err := h.service.GetRecent(limit)
	if err != nil {
		log.Error().Err(err).Int("limit", limit).Msg("failed to get recent billings")
		http.Error(w, "Failed to get recent billings", http.StatusInternalServerError)
		return
	}

	respondwithJSON(w, http.StatusOK, billings)
}

// r.Put("/{id}", billingHandler.Update)
// r.Delete("/{id}", billingHandler.Delete)
// r.Get("/client/{clientId}", billingHandler.GetByClientID)

// Update handles PUT /api/billings/{id}
func (h *BillingHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing billing ID", http.StatusBadRequest)
		return
	}

	var input domain.BillingInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	// Validate input
	if input.ClientID == "" || input.PackageID == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	billing, err := h.service.Update(id, input)
	if err != nil {
		log.Error().Err(err).Str("id", id).Interface("input", input).Msg("failed to update billing")
		if err.Error() == "billing not found" {
			http.Error(w, "Billing not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update billing", http.StatusInternalServerError)
		return
	}

	respondwithJSON(w, http.StatusOK, billing)
}
