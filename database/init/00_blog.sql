CREATE TABLE site_settings (
    key TEXT NOT NULL,
    val TEXT NOT NULL,
    PRIMARY KEY (key)
);
INSERT INTO site_settings (key, val) VALUES
('Version', '00.00.00'), -- Will be critical for adding version migration scripts eventually.
('Name', 'Blog'),
('URL', 'http://localhost:8080');

CREATE TABLE profile_pictures (
    pfp_id SERIAL,
    user_uploaded BOOLEAN NOT NULL,
    url VARCHAR(256) NOT NULL,
    PRIMARY KEY (pfp_id)
);

CREATE TABLE users (
    uid SERIAL,
    username VARCHAR(32) NOT NULL,
    password_hash VARCHAR(60) NOT NULL,
    pfp_id INT NOT NULL,
    privilege SMALLINT NOT NULL,
    FOREIGN KEY (pfp_id) REFERENCES profile_pictures(pfp_id),
    PRIMARY KEY (uid)
);

CREATE TABLE sessions (
    session_token VARCHAR(44) NOT NULL,
    uid INT NOT NULL,
    expire_date TIMESTAMP NOT NULL,
    FOREIGN KEY (uid) REFERENCES users(uid),
    PRIMARY KEY (session_token)
);

CREATE TABLE languages (
    lang_code VARCHAR(2) NOT NULL,
    lang_name VARCHAR(32) NOT NULL,
    page_tags TEXT NOT NULL,
    is_primary BOOL NOT NULL,
    PRIMARY KEY (lang_code)
);

CREATE TABLE error_codes (
    code_id VARCHAR(16) NOT NULL,
    lang_code VARCHAR(2) NOT NULL,
    content VARCHAR(256) NOT NULL,
    FOREIGN KEY (lang_code) REFERENCES languages(lang_code),
    PRIMARY KEY (code_id, lang_code)
);

CREATE TABLE page_types (
    page_type_id SERIAL,
    type_name VARCHAR(64) NOT NULL,
    substitution_types JSONB NOT NULL,
    template_url VARCHAR(256) NOT NULL,
    PRIMARY KEY (page_type_id)
);

CREATE TABLE pages (
    page_id SERIAL,
    page_type_id INT NOT NULL,
    required_privilege INT NOT NULL,
    FOREIGN KEY (page_type_id) REFERENCES page_types(page_type_id),
    PRIMARY KEY (page_id)
);

CREATE TABLE translations (
    translation_id SERIAL,
    page_id INT NOT NULL,
    lang_code VARCHAR(2) NOT NULL,
    title VARCHAR(32) NOT NULL,
    url VARCHAR(256) NOT NULL,
    UNIQUE (page_id, lang_code),
    UNIQUE (url),
    FOREIGN KEY (page_id) REFERENCES pages(page_id),
    FOREIGN KEY (lang_code) REFERENCES languages(lang_code),
    PRIMARY KEY (translation_id)
);

CREATE TABLE tests (
    test_id VARCHAR(2) NOT NULL,
    translation_id INT NOT NULL,
    substitutions JSONB NOT NULL,
    FOREIGN KEY (translation_id) REFERENCES translations(translation_id),
    PRIMARY KEY (test_id, translation_id)
);

CREATE INDEX idx_test_translationid ON tests (translation_id);

CREATE TABLE comments (
    comment_id SERIAL,
    translation_id INT NOT NULL,
    uid INT NOT NULL,
    container_id INT,
    content TEXT,
    FOREIGN KEY (translation_id) REFERENCES translations(translation_id),
    FOREIGN KEY (uid) REFERENCES users(uid),
    FOREIGN KEY (container_id) REFERENCES comments(comment_id),
    PRIMARY KEY (comment_id)
)
