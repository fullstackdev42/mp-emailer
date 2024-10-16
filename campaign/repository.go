package campaign

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/fullstackdev42/mp-emailer/internal/database"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(campaign *Campaign) error {
	query := `INSERT INTO campaigns (name, template, owner_id)VALUES (?, ?, ?)`
	result, err := r.db.SQL.Exec(query, campaign.Name, campaign.Template, campaign.OwnerID)
	if err != nil {
		return fmt.Errorf("error creating campaign: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert ID: %w", err)
	}
	campaign.ID = int(id)
	return nil
}

func (r *Repository) GetAll() ([]Campaign, error) {
	query := "SELECT id, name, template, owner_id, created_at, updated_at FROM campaigns"
	rows, err := r.db.SQL.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying campaigns: %w", err)
	}
	defer rows.Close()

	var campaigns []Campaign
	for rows.Next() {
		var c Campaign
		var createdAt, updatedAt sql.NullString
		err := rows.Scan(&c.ID, &c.Name, &c.Template, &c.OwnerID, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning campaign: %w", err)
		}
		c.CreatedAt, _ = parseDateTime(createdAt.String)
		c.UpdatedAt, _ = parseDateTime(updatedAt.String)
		campaigns = append(campaigns, c)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating campaigns: %w", err)
	}

	return campaigns, nil
}

func (r *Repository) Update(campaign *Campaign) error {
	query := "UPDATE campaigns SET name = ?, template = ?, updated_at = NOW() WHERE id = ?"
	_, err := r.db.SQL.Exec(query, campaign.Name, campaign.Template, campaign.ID)
	if err != nil {
		return fmt.Errorf("error updating campaign: %w", err)
	}
	return nil
}

func (r *Repository) Delete(id int) error {
	query := "DELETE FROM campaigns WHERE id = ?"
	_, err := r.db.SQL.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting campaign: %w", err)
	}
	return nil
}

func (r *Repository) GetByID(id int) (*Campaign, error) {
	query := "SELECT id, name, template, owner_id, created_at, updated_at FROM campaigns WHERE id = ?"
	row := r.db.SQL.QueryRow(query, id)

	var c Campaign
	var createdAt, updatedAt sql.NullString
	err := row.Scan(&c.ID, &c.Name, &c.Template, &c.OwnerID, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("error scanning campaign: %w", err)
	}

	c.CreatedAt, _ = parseDateTime(createdAt.String)
	c.UpdatedAt, _ = parseDateTime(updatedAt.String)

	return &c, nil
}

func (r *Repository) GetCampaign(id int) (*Campaign, error) {
	return r.GetByID(id)
}

func parseDateTime(dateStr string) (time.Time, error) {
	if dateStr == "0000-00-00 00:00:00" {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02 15:04:05", dateStr)
}
