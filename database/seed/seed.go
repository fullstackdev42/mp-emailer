package seed

import (
	"database/sql"
)

func Data(db *sql.DB) error {
	_, err := db.Exec(`
        INSERT INTO users (id, username, email, password_hash, created_at, updated_at) 
        VALUES (UUID(), 'foobar', 'jonesrussell42@gmail.com', 
        '$2a$10$7U0oMJZ0qtKcrJPI0otrXOTczXRfHdYD64JZ6oB2QTluNMSF9zmE6', 
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

        SET @user_id = (SELECT id FROM users WHERE username = 'foobar');

        INSERT INTO campaigns (id, name, description, template, owner_id, created_at, updated_at) 
        VALUES
        (UUID(), 'Unmarked Burials', 'Urge for increased funding...', 
        '... template content ...', @user_id, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
        (UUID(), 'Climate Action Now', 'Advocate for stronger climate policies...', 
        '... template content ...', @user_id, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
    `)
	return err
}
