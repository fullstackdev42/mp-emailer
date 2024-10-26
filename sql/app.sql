-- Verify new user registration
SELECT * FROM users WHERE username = 'newuser';

-- Check newly created campaign
SELECT * FROM campaigns WHERE name = 'New Campaign';

-- Verify campaign update
SELECT * FROM campaigns WHERE id = 1 AND name = 'Updated Campaign';

-- Check if campaign was deleted (should return no results if deleted)
SELECT * FROM campaigns WHERE id = 1;

-- List all campaigns
SELECT * FROM campaigns;

-- Get user details
SELECT * FROM users WHERE username = 'username';

-- Verify campaign template update
SELECT id, name, template FROM campaigns WHERE id = 1;

-- Check all campaigns for a specific user (assuming owner_id is stored)
SELECT * FROM campaigns WHERE owner_id = (SELECT id FROM users WHERE username = 'existinguser');

-- Create a new campaign
INSERT INTO campaigns (name, description, template, owner_id) 
VALUES ('New Campaign', 'This is a test campaign', 'Hello {{name}}, this is your campaign content.', '9aeaef88-3b8a-4df7-a400-d657ad3097a9');

-- Verify the newly created campaign
SELECT * FROM campaigns WHERE name = 'New Campaign';
