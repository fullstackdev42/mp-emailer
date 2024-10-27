-- @block Verify new user registration
-- @conn MP Emailer DB
-- @label Check new user
-- @group Users
-- @name verifyNewUser
SELECT * FROM users WHERE username = 'newuser';

-- @block Check newly created campaign
-- @conn MP Emailer DB
-- @label Verify new campaign
-- @group Campaigns
-- @name checkNewCampaign
SELECT * FROM campaigns WHERE name = 'New Campaign';

-- @block Verify campaign update
-- @conn MP Emailer DB
-- @label Check updated campaign
-- @group Campaigns
-- @name verifyCampaignUpdate
SELECT * FROM campaigns WHERE id = 1 AND name = 'Updated Campaign';

-- @block Check if campaign was deleted
-- @conn MP Emailer DB
-- @label Verify campaign deletion
-- @group Campaigns
-- @name checkCampaignDeletion
SELECT * FROM campaigns WHERE id = 1;

-- @block List all campaigns
-- @conn MP Emailer DB
-- @label Get all campaigns
-- @group Campaigns
-- @name getAllCampaigns
SELECT * FROM campaigns;

-- @block Get user details
-- @conn MP Emailer DB
-- @label Fetch user info
-- @group Users
-- @name getUserDetails
SELECT * FROM users WHERE username = 'username';

-- @block Verify campaign template update
-- @conn MP Emailer DB
-- @label Check campaign template
-- @group Campaigns
-- @name verifyCampaignTemplate
SELECT id, name, template FROM campaigns WHERE id = 1;

-- @block Check all campaigns for a specific user
-- @conn MP Emailer DB
-- @label User's campaigns
-- @group Campaigns
-- @name getUserCampaigns
SELECT * FROM campaigns WHERE owner_id = (SELECT id FROM users WHERE username = 'existinguser');

-- @block Create new campaign
-- @conn MP Emailer DB
-- @label Insert new campaign
-- @group Campaigns
-- @name createCampaign
INSERT INTO campaigns (name, description, template, owner_id) 
VALUES ('New Campaign', 'This is a test campaign', 'Hello {{name}}, this is your campaign content.', '9aeaef88-3b8a-4df7-a400-d657ad3097a9');

-- @block Verify the newly created campaign
-- @conn MP Emailer DB
-- @label Check new campaign
-- @group Campaigns
-- @name verifyNewCampaign
SELECT * FROM campaigns WHERE name = 'New Campaign';

-- @block Update an existing campaign
-- @conn MP Emailer DB
-- @label Update campaign
-- @group Campaigns
-- @name updateCampaign
UPDATE campaigns
SET name = 'Updated Campaign', 
    description = 'This campaign has been updated', 
    template = 'Hello {{name}}, this is your updated campaign content.'
WHERE id = 1;

-- @block Delete a campaign
-- @conn MP Emailer DB
-- @label Remove campaign
-- @group Campaigns
-- @name deleteCampaign
DELETE FROM campaigns
WHERE id = 2;

-- @block Get a specific campaign by ID
-- @conn MP Emailer DB
-- @label Fetch campaign by ID
-- @group Campaigns
-- @name getCampaignById
SELECT id, name, description, template, owner_id, created_at, updated_at
FROM campaigns
WHERE id = 3;
