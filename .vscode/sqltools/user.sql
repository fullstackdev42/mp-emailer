-- @block Get all users
-- @conn MPEmailerDB
-- @label Get all users
-- @group Users
-- @name getAllUsers
SELECT * FROM users;

-- @block Verify new user registration
-- @conn MPEmailerDB
-- @label Check new user
-- @group Users
-- @name verifyNewUser
SELECT * FROM users WHERE username = 'newuser';

-- @block Get user details
-- @conn MPEmailerDB
-- @label Fetch user info
-- @group Users
-- @name getUserDetails
SELECT * FROM users WHERE username = 'username';

-- @block Create new user
-- @conn MPEmailerDB
-- @label Insert new user
-- @group Users
-- @name createUser
INSERT INTO users (id, username, email, password_hash) 
VALUES (UUID(), 'newuser', 'newuser@example.com', 'hashed_password_here');

-- @block Verify new user creation
-- @conn MPEmailerDB
-- @label Check new user
-- @group Users
-- @name verifyNewUser
SELECT * FROM users WHERE username = 'newuser';

