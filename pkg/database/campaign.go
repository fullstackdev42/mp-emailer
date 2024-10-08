package database

import (
	"fmt"
	"strings"

	"github.com/fullstackdev42/mp-emailer/pkg/models"
)

func (db *DB) CreateCampaign(campaign *models.Campaign) error {
	query := `
		INSERT INTO campaigns (id, name, template, owner_id, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := db.Exec(query,
		campaign.ID,
		campaign.Name,
		campaign.Template,
		campaign.OwnerID,
		campaign.CreatedAt,
		campaign.UpdatedAt,
	)
	if err != nil {
		// Check for duplicate entry error
		if isDuplicateEntryError(err) {
			return models.ErrDuplicateCampaign
		}
		return fmt.Errorf("error creating campaign: %w", err)
	}
	return nil
}

// Add this helper function to check for duplicate entry errors
func isDuplicateEntryError(err error) bool {
	return strings.Contains(err.Error(), "Duplicate entry")
}

func (db *DB) GetCampaigns() ([]models.Campaign, error) {
	query := "SELECT id, name, template FROM campaigns"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error fetching campaigns: %w", err)
	}
	defer rows.Close()

	var campaigns []models.Campaign
	for rows.Next() {
		var campaign models.Campaign
		if err := rows.Scan(&campaign.ID, &campaign.Name, &campaign.Template); err != nil {
			return nil, fmt.Errorf("error scanning campaign: %w", err)
		}
		campaigns = append(campaigns, campaign)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error with rows: %w", err)
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

func (db *DB) DeleteCampaign(id string) error {
	query := "DELETE FROM campaigns WHERE id = ?"
	_, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting campaign: %w", err)
	}
	return nil
}
