package campaign

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/fullstackdev42/mp-emailer/internal/database"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/jonesrussell/loggo"
)

// RepositoryInterface defines the methods that a campaign repository must implement
type RepositoryInterface interface {
	Create(dto *CreateCampaignDTO) error
	GetAll() ([]Campaign, error)
	Update(dto *UpdateCampaignDTO) error
	Delete(dto DeleteCampaignDTO) error
	GetByID(dto GetCampaignDTO) (*Campaign, error)
	GetCampaign(dto GetCampaignDTO) (*Campaign, error)
}

// Ensure that Repository implements RepositoryInterface
var _ RepositoryInterface = (*Repository)(nil)

type Repository struct {
	db     *database.DB
	logger loggo.LoggerInterface
}

func (r *Repository) Create(dto *CreateCampaignDTO) error {
	query := `INSERT INTO campaigns (name, description, template, owner_id) VALUES (?, ?, ?, ?)`
	result, err := r.db.SQL.Exec(query, dto.Name, dto.Description, dto.PostalCode, dto.Template, dto.OwnerID)
	if err != nil {
		return fmt.Errorf("error creating campaign: %w", err)
	}
	_, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert ID: %w", err)
	}
	return nil
}

func (r *Repository) GetAll() ([]Campaign, error) {
	query := "SELECT id, name, description, template, owner_id, created_at, updated_at FROM campaigns"
	rows, err := r.db.SQL.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying campaigns: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			r.logger.Error("Error closing rows", err)
		}
	}(rows)

	var campaigns []Campaign
	for rows.Next() {
		var c Campaign
		var createdAt, updatedAt sql.NullString
		err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.Template, &c.OwnerID, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning campaign: %w", err)
		}
		c.CreatedAt, _ = shared.ParseDateTime(createdAt.String)
		c.UpdatedAt, _ = shared.ParseDateTime(updatedAt.String)
		campaigns = append(campaigns, c)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating campaigns: %w", err)
	}

	return campaigns, nil
}

func (r *Repository) Update(dto *UpdateCampaignDTO) error {
	query := "UPDATE campaigns SET name = ?, description = ?, postal_code = ?, template = ?, updated_at = NOW() WHERE id = ?"
	_, err := r.db.SQL.Exec(query, dto.Name, dto.Description, dto.PostalCode, dto.Template, dto.ID)
	if err != nil {
		return fmt.Errorf("error updating campaign: %w", err)
	}
	return nil
}

func (r *Repository) Delete(dto DeleteCampaignDTO) error {
	query := "DELETE FROM campaigns WHERE id = ?"
	_, err := r.db.SQL.Exec(query, dto.ID)
	if err != nil {
		return fmt.Errorf("error deleting campaign: %w", err)
	}
	return nil
}

func (r *Repository) GetByID(dto GetCampaignDTO) (*Campaign, error) {
	query := "SELECT id, name, description, template, owner_id, created_at, updated_at FROM campaigns WHERE id = ?"
	row := r.db.SQL.QueryRow(query, dto.ID)

	var c Campaign
	var createdAt, updatedAt sql.NullString
	err := row.Scan(&c.ID, &c.Name, &c.Description, &c.Template, &c.OwnerID, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCampaignNotFound
		}
		return nil, fmt.Errorf("error scanning campaign: %w", err)
	}

	c.CreatedAt, _ = shared.ParseDateTime(createdAt.String)
	c.UpdatedAt, _ = shared.ParseDateTime(updatedAt.String)

	return &c, nil
}

func (r *Repository) GetCampaign(dto GetCampaignDTO) (*Campaign, error) {
	return r.GetByID(dto)
}
