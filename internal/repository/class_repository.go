package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/matthieukhl/align-back/internal/domain"
	"github.com/rs/zerolog/log"
)

type classRepository struct {
	db *sqlx.DB
}

// NewClassRepository creates a new class repository
func NewClassRepository(db *sqlx.DB) domain.ClassRepository {
	return &classRepository{
		db: db,
	}
}

// Create creates a class.
func (r *classRepository) Create(class *domain.Class) error {
	query := `
	INSERT INTO
		name
		, location
		, type
		, equipment
	VALUES (?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, class.Name, class.Location, class.Type, class.Equipment)
	if err != nil {
		log.Error().Err(err).Interface("class", class).Msg("failed to create class")
		return fmt.Errorf("failed to create class: %w", err)
	}

	return nil

}

// Delete deletes a class.
func (r *classRepository) Delete(id string) error {
	query := `
	DELETE FROM
		classes
	WHERE 
		id = ?
	`

	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("failed to delete class")
		return fmt.Errorf("failed to delete class: %w", err)
	}

	return nil
}

// GetAll returns all classes.
func (r *classRepository) GetAll() ([]domain.Class, error) {
	var classes []domain.Class

	query := `
	SELECT
		id
		, name
		, location
		, type
		, equipment
	FROM 
		classes
	`

	err := r.db.Select(&classes, query)
	if err != nil {
		log.Error().Err(err).Msg("failed to get all classes")
		return nil, fmt.Errorf("failed to get all classes: %w", err)
	}

	return classes, nil
}

// GetByID implements domain.ClassRepository.
func (r *classRepository) GetByID(id string) (*domain.Class, error) {
	var class domain.Class

	query := `
	SELECT 
		id
		, name
		, location
		, type
		, equipment
	FROM 
		classes
	WHERE 
		id = ?
	`

	err := r.db.Get(&class, query, id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("failed to get class")
		return nil, fmt.Errorf("failed to get class: %w", err)
	}

	return &class, nil
}

// GetByLocation implements domain.ClassRepository.
func (r *classRepository) GetByLocation(location domain.Location) ([]domain.Class, error) {
	var classes []domain.Class

	query := `
	SELECT
		id
		, name
		, location
		, type
		, equipment
	FROM 
		classes
	WHERE 
		location = ?
	`

	err := r.db.Select(&classes, query, location)
	if err != nil {
		log.Error().Err(err).Interface("location", location).Msg("failed to get class by location")
		return nil, fmt.Errorf("failed to get class by location: %w", err)
	}

	return classes, nil
}

// GetByType implements domain.ClassRepository.
func (r *classRepository) GetByType(classType domain.ClassType) ([]domain.Class, error) {
	var classes []domain.Class

	query := `
	SELECT
		id
		, name
		, location
		, type
		, equipment
	FROM 
		classes
	WHERE
		type = ? 
	`

	err := r.db.Select(&classes, query, classType)
	if err != nil {
		log.Error().Err(err).Interface("class type", classType).Msg("failed to get class by type")
		return nil, fmt.Errorf("failed to get class by type: %w", err)
	}

	return classes, nil
}

// Update implements domain.ClassRepository.
func (r *classRepository) Update(class *domain.Class) error {
	query := `
	UPDATE 
		classes
	SET 
		name = ?
		, location = ?
		, type = ?
		, equipment = ?
	WHERE 
		id = ?
	`

	_, err := r.db.Exec(query, class.Name, class.Location, class.Type, class.Equipment, class.ID)
	if err != nil {
		log.Error().Err(err).Interface("class", class).Msg("failed to update class")
		return fmt.Errorf("failed to update class: %w", err)
	}

	return nil
}

// GetByName returns a class by name
func (r *classRepository) GetByName(name string) (*domain.Class, error) {
	var class domain.Class

	query := `
	SELECT
		id
		, name
		, location
		, type
		, equipment
	WHERE 
		name = ?
	`

	err := r.db.Get(&class, query, name)
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("failed to get class by name")
		return nil, fmt.Errorf("failed to get class by name: %w", err)
	}

	return &class, nil
}
