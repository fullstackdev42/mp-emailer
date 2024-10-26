-- @block Bookmarked query
-- @group Campaigns
-- @name Campaign create

-- Create a new campaign
INSERT INTO campaigns (name, description, template, owner_id) 
VALUES ('New Campaign', 'This is a test campaign', 'Hello {{name}}, this is your campaign content.', '9aeaef88-3b8a-4df7-a400-d657ad3097a9');
