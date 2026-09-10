-- TODO: Create a good set of defaults, and add a way to add/change them
INSERT INTO profile_pictures (user_uploaded, url) VALUES
(FALSE, '/res/default_pfp.png');

INSERT INTO users (username, password_hash, pfp_id, privilege) VALUES
('ANONYMOUS', 'Locked - This is not a hash', 1, 0),
('DELETED', 'Locked - This is not a hash', 1, 0);
