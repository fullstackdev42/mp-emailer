package campaign

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/fullstackdev42/mp-emailer/database"
	"github.com/fullstackdev42/mp-emailer/shared"
)

// RepositoryInterface defines the methods that a campaign repository must implement
type RepositoryInterface interface {
	Create(dto *CreateCampaignDTO) (*Campaign, error)
	GetAll() ([]Campaign, error)
	Update(dto *UpdateCampaignDTO) error
	Delete(dto DeleteCampaignDTO) error
	GetByID(dto GetCampaignDTO) (*Campaign, error)
	GetCampaign(dto GetCampaignDTO) (*Campaign, error)
}

// Ensure that Repository implements RepositoryInterface
var _ RepositoryInterface = (*Repository)(nil)

// Repository handles CRUD operations for campaigns
type Repository struct {
	db database.Interface
}

// Create creates a new campaign in the database
func (r *Repository) Create(dto *CreateCampaignDTO) (*Campaign, error) {
	query := `INSERT INTO campaigns (name, description, template, owner_id) VALUES (?, ?, ?, ?)`
	result, err := r.db.Exec(query, dto.Name, dto.Description, dto.Template, dto.OwnerID)
	if err != nil {
		return nil, fmt.Errorf("error creating campaign: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %w", err)
	}
	campaign := &Campaign{
		ID:          int(id),
		Name:        dto.Name,
		Description: dto.Description,
		Template:    dto.Template,
		OwnerID:     dto.OwnerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return campaign, nil
}

// GetAll retrieves all campaigns from the database
func (r *Repository) GetAll() ([]Campaign, error) {
	query := "SELECT id, name, description, template, owner_id, created_at, updated_at FROM campaigns"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying campaigns: %w", err)
	}
	defer rows.Close()

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

// Update updates an existing campaign in the database
func (r *Repository) Update(dto *UpdateCampaignDTO) error {
	query := "UPDATE campaigns SET name = ?, description = ?, template = ?, updated_at = NOW() WHERE id = ?"
	_, err := r.db.Exec(query, dto.Name, dto.Description, dto.Template, dto.ID)
	if err != nil {
		return fmt.Errorf("error updating campaign: %w", err)
	}
	return nil
}

// Delete deletes a campaign from the database
func (r *Repository) Delete(dto DeleteCampaignDTO) error {
	query := "DELETE FROM campaigns WHERE id = ?"
	_, err := r.db.Exec(query, dto.ID)
	if err != nil {
		return fmt.Errorf("error deleting campaign: %w", err)
	}
	return nil
}

// GetByID retrieves a campaign by its ID
func (r *Repository) GetByID(dto GetCampaignDTO) (*Campaign, error) {
	query := "SELECT id, name, description, template, owner_id, created_at, updated_at FROM campaigns WHERE id = ?"
	row := r.db.QueryRow(query, dto.ID)
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

// GetCampaign retrieves a campaign by its parameters
func (r *Repository) GetCampaign(dto GetCampaignDTO) (*Campaign, error) {
	return r.GetByID(dto)
}
