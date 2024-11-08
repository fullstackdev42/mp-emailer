-- Users table
CREATE TABLE IF NOT EXISTS users (
    id CHAR(36) PRIMARY KEY DEFAULT UUID(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_users_deleted_at ON users(deleted_at); 

-- Campaigns table
CREATE TABLE IF NOT EXISTS campaigns (
    id CHAR(36) PRIMARY KEY DEFAULT UUID(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    template TEXT NOT NULL,
    owner_id CHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE INDEX idx_campaigns_deleted_at ON campaigns(deleted_at);
