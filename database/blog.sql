CREATE TABLE users (
    username VARCHAR(32) NOT NULL,
    pash VARCHAR(60) NOT NULL,
    pfp_file_id VARCHAR(32) NOT NULL,
    privilege SMALLINT NOT NULL,
    PRIMARY KEY (username)
);
INSERT INTO users VALUES ('username', 'password', 'file', 0);
