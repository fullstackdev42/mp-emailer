package database

import (
	"fmt"

	"github.com/fullstackdev42/mp-emailer/pkg/models"
)

func (db *DB) CreateCampaign(campaign *models.Campaign) error {
	query := `
		INSERT INTO campaigns (name, template, owner_id, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := db.Exec(query,
		campaign.Name,
		campaign.Template,
		campaign.OwnerID,
		campaign.CreatedAt,
		campaign.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error creating campaign: %w", err)
	}

	// Get the auto-generated ID
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert ID: %w", err)
	}

	campaign.ID = int(id)
	return nil
}

func (db *DB) GetCampaigns() ([]models.Campaign, error) {
	rows, err := db.Query("SELECT id, name, template, owner_id FROM campaigns")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var campaigns []models.Campaign
	for rows.Next() {
		var campaign models.Campaign
		err := rows.Scan(&campaign.ID, &campaign.Name, &campaign.Template, &campaign.OwnerID)
		if err != nil {
			return nil, err
		}
		campaigns = append(campaigns, campaign)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return campaigns, nil
}

func (db *DB) UpdateCampaign(campaign *models.Campaign) error {
	query := "UPDATE campaigns SET name = ?, template = ? WHERE id = ?"
	_, err := db.Exec(query, campaign.Name, campaign.Template, campaign.ID)
	if err != nil {
		return fmt.Errorf("error updating campaign: %w", err)
	}
	return nil
}

func (db *DB) DeleteCampaign(id int) error {
	query := "DELETE FROM campaigns WHERE id = ?"
	_, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting campaign: %w", err)
	}
	return nil
}

// GetCampaignByID retrieves a campaign from the database by its ID
func (db *DB) GetCampaignByID(id int) (*models.Campaign, error) {
	query := "SELECT id, name, template, owner_id, created_at, updated_at FROM campaigns WHERE id = ?"
	campaign := &models.Campaign{}
	err := db.QueryRow(query, id).Scan(
		&campaign.ID,
		&campaign.Name,
		&campaign.Template,
		&campaign.OwnerID,
		&campaign.CreatedAt,
		&campaign.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return campaign, nil
}
