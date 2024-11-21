-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS campaigns (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    template TEXT NOT NULL,
    owner_id CHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (owner_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER before_campaigns_insert 
BEFORE INSERT ON campaigns
FOR EACH ROW
SET NEW.id = UUID();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_campaigns_deleted_at ON campaigns(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS campaigns;
-- +goose StatementEnd
