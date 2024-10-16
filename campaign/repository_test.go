package campaign

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fullstackdev42/mp-emailer/internal/database"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, RepositoryInterface) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockDB := &database.DB{SQL: db}
	repo := NewRepository(mockDB)

	return db, mock, repo
}

func TestRepository_Create(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	campaign := &Campaign{
		Name:     "Test Campaign",
		Template: "Test Template",
		OwnerID:  1,
	}

	mock.ExpectExec("INSERT INTO campaigns").
		WithArgs(campaign.Name, campaign.Template, campaign.OwnerID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(campaign)

	assert.NoError(t, err)
	assert.Equal(t, 1, campaign.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetAll(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "template", "owner_id", "created_at", "updated_at"}).
		AddRow(1, "Campaign 1", "Template 1", 1, time.Now(), time.Now()).
		AddRow(2, "Campaign 2", "Template 2", 2, time.Now(), time.Now())

	mock.ExpectQuery("SELECT (.+) FROM campaigns").WillReturnRows(rows)

	campaigns, err := repo.GetAll()

	assert.NoError(t, err)
	assert.Len(t, campaigns, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Update(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	campaign := &Campaign{
		ID:       1,
		Name:     "Updated Campaign",
		Template: "Updated Template",
	}

	mock.ExpectExec("UPDATE campaigns SET").
		WithArgs(campaign.Name, campaign.Template, campaign.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(campaign)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Delete(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	id := 1

	mock.ExpectExec("DELETE FROM campaigns WHERE").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Delete(id)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByID(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	id := 1
	rows := sqlmock.NewRows([]string{"id", "name", "template", "owner_id", "created_at", "updated_at"}).
		AddRow(id, "Campaign 1", "Template 1", 1, time.Now(), time.Now())

	mock.ExpectQuery("SELECT (.+) FROM campaigns WHERE").
		WithArgs(id).
		WillReturnRows(rows)

	campaign, err := repo.GetByID(id)

	assert.NoError(t, err)
	assert.NotNil(t, campaign)
	assert.Equal(t, id, campaign.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetCampaign(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	id := 1
	rows := sqlmock.NewRows([]string{"id", "name", "template", "owner_id", "created_at", "updated_at"}).
		AddRow(id, "Campaign 1", "Template 1", 1, time.Now(), time.Now())

	mock.ExpectQuery("SELECT (.+) FROM campaigns WHERE").
		WithArgs(id).
		WillReturnRows(rows)

	campaign, err := repo.GetCampaign(id)

	assert.NoError(t, err)
	assert.NotNil(t, campaign)
	assert.Equal(t, id, campaign.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
