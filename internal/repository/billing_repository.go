package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/matthieukhl/align-back/internal/domain"
	"github.com/rs/zerolog/log"
)

type billingRepository struct {
	db *sqlx.DB
}

// NewBillingRepository creates a new billing repository
func NewBillingRepository(db *sqlx.DB) domain.BillingRepository {
	return &billingRepository{
		db: db,
	}
}

// Create creates a new billing.
func (r *billingRepository) Create(billing *domain.Billing) error {
	query := `
	INSERT INTO 
		billings (
			client_id
			, package_id
			, amount
			, price
			, credits
			, payment_date
		)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, billing.ClientID, billing.PackageID, billing.Amount, billing.Price, billing.Credits, billing.PaymentDate)
	if err != nil {
		log.Error().Err(err).Interface("billing", billing).Msg("failed to create billing")
		return fmt.Errorf("failed to create billing: %w", err)
	}

	return nil
}

// Delete implements domain.BillingRepository.
func (r *billingRepository) Delete(id string) error {
	query := `
	DELETE FROM
		billings
	WHERE 
		id = ?
	`

	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("failed to delete billing")
		return fmt.Errorf("failed to delete billing: %w", err)
	}

	return nil
}

// GetAll returns all billings
func (r *billingRepository) GetAll() ([]domain.Billing, error) {
	var billings []domain.Billing

	query := `
	SELECT 
		id
		, client_id
		, package_id
		, amount
		, price
		, credits
		, payment_date
	FROM 
		billings
	`

	err := r.db.Select(&billings, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error().Err(err).Msg("failed to get billings")
		return nil, fmt.Errorf("failed to get billings: %w", err)
	}
	return billings, nil
}

// GetAllWithDetails returns all billings with their details
func (r *billingRepository) GetAllWithDetails() ([]domain.BillingWithDetails, error) {
	panic("unimplemented")
}

// GetByClient implements domain.BillingRepository.
func (r *billingRepository) GetByClient(clientID string) ([]domain.Billing, error) {
	var billings []domain.Billing

	query := `
	SELECT 
		id
		, client_id
		, package_id
		, amount
		, price
		, credits
		, payment_date
	FROM 
		billings
	WHERE
		client_id = ?
	`

	err := r.db.Select(&billings, query, clientID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error().Err(err).Str("clientID", clientID).Msg("failed to get billings by client ID")
		return nil, fmt.Errorf("failed to get billings by client ID")
	}

	return billings, nil
}

// GetByID implements domain.BillingRepository.
func (r *billingRepository) GetByID(id string) (*domain.Billing, error) {
	var billing domain.Billing

	query := `
	SELECT 
		id
		, client_id
		, package_id
		, amount
		, price
		, credits
		, payment_date
	FROM
		billings
	WHERE 
		id = ?
	`

	err := r.db.Get(&billing, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error().Err(err).Str("id", id).Msg("failed to get billing by ID")
		return nil, fmt.Errorf("failed to get billing by ID: %w", err)
	}

	return &billing, nil
}

// GetRecent returns recent billings based on a limit.
func (r *billingRepository) GetRecent(limit int) ([]domain.Billing, error) {
	var billings []domain.Billing

	query := `
	SELECT 
		id
		, client_id
		, package_id
		, amount
		, price
		, credits
		, payment_date
	FROM
		billings
	ORDER BY 
		payment_date DESC
	LIMIT 
		?
	`

	err := r.db.Select(&billings, query, limit)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error().Err(err).Int("limit", limit).Msg("failed to get recent billings")
		return nil, fmt.Errorf("failed to get recent billings: %w", err)
	}

	return billings, nil
}

// GetWithDetails returns a billing by ID with its details.
func (r *billingRepository) GetWithDetails(id string) (*domain.BillingWithDetails, error) {
	panic("unimplemented")
}

// Update updates an existing billing.
func (r *billingRepository) Update(billing *domain.Billing) error {
	query := `
	UPDATE
		billings
	SET
		client_id = ?
		, package_id = ?
		, amount = ?
		, price = ?
		, credits = ?
		, payment_date = ?
	WHERE
		id = ?
	`

	_, err := r.db.Exec(query, billing.ClientID, billing.PackageID, billing.Amount, billing.Price, billing.Credits, billing.PaymentDate)
	if err != nil {
		log.Error().Err(err).Interface("billing", billing).Msg("failed to update billing")
		return fmt.Errorf("failed to update billing")
	}

	return nil
}
