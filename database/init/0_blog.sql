CREATE TABLE users (
    uid SERIAL,
    username VARCHAR(32) NOT NULL,
    password_hash VARCHAR(60) NOT NULL,
    pfp_file_id VARCHAR(32) NOT NULL,
    privilege SMALLINT NOT NULL,
    PRIMARY KEY (uid)
);

-- This is the sort of thing to move to redis eventually
-- There is currently no expire logic
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

CREATE TABLE page_type (
    page_type_id SERIAL,
    type_name VARCHAR(64) NOT NULL,
    substitution_types JSONB NOT NULL,
    template_url VARCHAR(64) NOT NULL,
    PRIMARY KEY (page_type_id)
);

CREATE TABLE pages (
    page_id SERIAL,
    page_type_id INT NOT NULL,
    required_privilege INT NOT NULL,
    FOREIGN KEY (page_type_id) REFERENCES page_type(page_type_id),
    PRIMARY KEY (page_id)
);

CREATE TABLE translations (
    translation_id SERIAL,
    page_id INT NOT NULL,
    lang_code VARCHAR(2) NOT NULL,
    substitutions JSONB NOT NULL,
    url VARCHAR(32) NOT NULL,
    UNIQUE (page_id, lang_code),
    UNIQUE (url),
    FOREIGN KEY (page_id) REFERENCES pages(page_id),
    FOREIGN KEY (lang_code) REFERENCES languages(lang_code),
    PRIMARY KEY (translation_id)
);

CREATE TABLE tests (
    test_id VARCHAR(2) NOT NULL,
    translation_id INT NOT NULL,
    test_substitutions JSONB NOT NULL,
    FOREIGN KEY (translation_id) REFERENCES translations(translation_id),
    PRIMARY KEY (test_id, translation_id)
);

CREATE INDEX idx_test_translationid ON tests (translation_id);
