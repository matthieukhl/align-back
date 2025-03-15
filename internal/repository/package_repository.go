package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/matthieukhl/align-back/internal/domain"
	"github.com/rs/zerolog/log"
)

type packageRepository struct {
	db *sqlx.DB
}

// NewPackageRepository creates a new package repository
func NewPackageRepository(db *sqlx.DB) domain.PackageRepository {
	return &packageRepository{
		db: db,
	}
}

// Create creates a new Package.
func (r *packageRepository) Create(pkg *domain.Package) error {
	query := `
	INSERT INTO
		packages (
			name
			, number_of_sessions
			, type
			, price
		)
	VALUES (?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, pkg.Name, pkg.NumberOfSessions, pkg.Type, pkg.Price)
	if err != nil {
		log.Error().Err(err).Interface("package", pkg).Msg("failed to create package")
		return fmt.Errorf("failed to create package: %w", err)
	}

	return nil
}

// Delete implements domain.PackageRepository.
func (r *packageRepository) Delete(id string) error {
	query := `
	DELETE FROM
		packages
	WHERE
		id = ?
	`

	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("failed to delete package")
		return fmt.Errorf("failed to delete package: %w", err)
	}

	return nil
}

// GetAll implements domain.PackageRepository.
func (r *packageRepository) GetAll() ([]domain.Package, error) {
	var packages []domain.Package

	query := `
	SELECT * FROM packages
	`

	err := r.db.Select(&packages, query)
	if err != nil {
		log.Error().Err(err).Msg("failed to get all packages")
		return nil, fmt.Errorf("failed to get all packages: %w", err)
	}

	return packages, nil
}

// GetByID implements domain.PackageRepository.
func (r *packageRepository) GetByID(id string) (*domain.Package, error) {
	var pkg domain.Package

	query := `
	SELECT 
		*
	FROM 
		packages
	WHERE 
		id = ?
	`

	err := r.db.Get(&pkg, query, id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("failed to get package by ID")
		return nil, fmt.Errorf("failed to get package by ID: %w", err)
	}

	return &pkg, nil
}

// GetByType implements domain.PackageRepository.
func (r *packageRepository) GetByType(pkgType domain.PackageType) ([]domain.Package, error) {
	var packages []domain.Package

	query := `
	SELECT
		id
		, name
		, number_of_sessions
		, type
		, price
	FROM 
		packages
	WHERE 
		type = ?
	`

	err := r.db.Select(&packages, query, pkgType)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error().Err(err).Interface("pacage type", pkgType).Msg("failed to get package by package type")
		return nil, fmt.Errorf("failed to get package by package type")
	}

	return packages, nil
}

// Update implements domain.PackageRepository.
func (r *packageRepository) Update(pkg *domain.Package) error {
	query := `
	UPDATE
		packages
	SET 
		name = ?
		, number_of_sessions = ?
		, type = ?
		, price = ?
	WHERE
		id = ?
	`

	_, err := r.db.Exec(query, pkg.Name, pkg.NumberOfSessions, pkg.Type, pkg.Price, pkg.ID)
	if err != nil {
		log.Error().Err(err).Interface("package", pkg).Msg("failed to update package")
		return fmt.Errorf("failed to update client: %w", err)
	}

	return nil
}
