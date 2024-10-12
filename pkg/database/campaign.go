package database

import (
	"fmt"
	"time"

	"github.com/fullstackdev42/mp-emailer/pkg/models"
)

func (db *DB) CreateCampaign(campaign *models.Campaign) error {
	query := `
        INSERT INTO campaigns (name, template, owner_id)
        VALUES (?, ?, ?)
    `
	result, err := db.Exec(query, campaign.Name, campaign.Template, campaign.OwnerID)
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

// GetCampaignByID retrieves a campaign by its ID
func (db *DB) GetCampaignByID(id int) (*models.Campaign, error) {
	query := "SELECT id, name, template, owner_id, created_at, updated_at FROM campaigns WHERE id = ?"
	row := db.QueryRow(query, id)

	var createdAt, updatedAt []byte
	campaign := &models.Campaign{}

	err := row.Scan(&campaign.ID, &campaign.Name, &campaign.Template, &campaign.OwnerID, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	// Convert []byte to time.Time
	campaign.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
	if err != nil {
		return nil, fmt.Errorf("error parsing created_at: %w", err)
	}

	campaign.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", string(updatedAt))
	if err != nil {
		return nil, fmt.Errorf("error parsing updated_at: %w", err)
	}

	return campaign, nil
}
