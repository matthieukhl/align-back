package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/matthieukhl/align-back/internal/domain"
	"github.com/rs/zerolog/log"
)

type clientRepository struct {
	db *sqlx.DB
}

// NewClientRepository creates a new client repository
func NewClientRepository(db *sqlx.DB) domain.ClientRepository {
	return &clientRepository{
		db: db,
	}
}

// GetAll returns all clients
func (r *clientRepository) GetAll() ([]domain.Client, error) {
	var clients []domain.Client

	query := `SELECT * FROM clients ORDER BY lastname, firstname`

	err := r.db.Select(&clients, query)
	if err != nil {
		log.Error().Err(err).Msg("failed to get all clients")
		return nil, fmt.Errorf("failed to get all clients: %w", err)
	}

	return clients, nil
}

// GetByID returns a client by ID
func (r *clientRepository) GetByID(id string) (*domain.Client, error) {
	var client domain.Client

	query := `SELECT * FROM clients WHERE id = ?`

	err := r.db.Get(&client, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error().Err(err).Str("id", id).Msg("failed to get client by ID")
		return nil, fmt.Errorf("failed to get client by ID: %w", err)
	}

	return &client, nil
}

// Create creates a new client
func (r *clientRepository) Create(client *domain.Client) error {
	query := `
	INSERT INTO clients (
		id
		, firstname
		, lastname
		, phone
		, email
		, street_number
		, street_name
		, city
		, zip_code
		, country
		, credits)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, client.FirstName, client.LastName, client.Phone, client.Email, client.StreetNumber, client.StreetName, client.City, client.ZipCode, client.Country, client.Credits)
	if err != nil {
		log.Error().Err(err).Interface("client", client).Msg("failed to create client")
		return fmt.Errorf("failed to create client: %w", err)
	}

	return nil
}

// Update updates a client's information
func (r *clientRepository) Update(client *domain.Client) error {
	query := `
	UPDATE
		clients
	SET
		firstname = ?
		, lastname = ?
		, phone = ?
		, email = ?
		, street_number = ?
		, street_name = ?
		, city = ?
		, zip_code = ?
		, country = ?
		, credits = ?
	WHERE
		id = ?`

	_, err := r.db.Exec(query, client.FirstName, client.LastName, client.Phone, client.Email, client.StreetNumber, client.StreetName, client.City, client.ZipCode, client.Country, client.Credits, client.ID)
	if err != nil {
		log.Error().Err(err).Interface("client", client).Msg("failed to update client")
		return fmt.Errorf("failed to update client: %w", err)
	}

	return nil
}

// Delete deletes a client by ID
func (r *clientRepository) Delete(id string) error {
	query := `
	DELETE FROM
		clients
	WHERE 
		id = ?
	`

	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("failed to delete client")
		return fmt.Errorf("failed to delete client: %w", err)
	}

	return nil
}

// GetByEmail returns a client by its email
func (r *clientRepository) GetByEmail(email string) (*domain.Client, error) {
	var client domain.Client

	query := `
	SELECT * FROM clients WHERE email = ?
	`

	err := r.db.Get(&client, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error().Err(err).Str("email", email).Msg("Failed to get client by email")
		return nil, fmt.Errorf("failed to get client by email: %w", err)
	}

	return &client, nil
}

// GetLowCredits returns clients with credits below a threshold
func (r *clientRepository) GetLowCredits(threshold int) ([]domain.Client, error) {
	var clients []domain.Client

	query := `
	SELECT
		*
	FROM
		clients
	WHERE
		credits <= ?
	ORDER BY  
		credits ASC, lastname, firstname
	`

	err := r.db.Select(&clients, query, threshold)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error().Err(err).Int("threshold", threshold).Msg("failed to get clients with low credits")
		return nil, fmt.Errorf("failed to get clients with low credits: %w", err)
	}

	return clients, nil
}
