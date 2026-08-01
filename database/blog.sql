CREATE TABLE users (
    uid SERIAL,
    username VARCHAR(32) NOT NULL,
    pash VARCHAR(60) NOT NULL,
    pfp_file_id VARCHAR(32) NOT NULL,
    privilege SMALLINT NOT NULL,
    PRIMARY KEY (uid)
);

-- This is the sort of thing to move to redis eventually
-- There is currently no expire logic
CREATE TABLE sessions (
    session_token VARCHAR(44) NOT NULL,
    uid INT NOT NULL,
    expire_date TIMESTAMP,
    FOREIGN KEY (uid) REFERENCES users(uid),
    PRIMARY KEY (session_token)
)
