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
    PRIMARY KEY (lang_code)
);
INSERT INTO languages VALUES
('en', 'English', ''),
(
    'ja',
    '日本語',
    '<meta name="robots" content="noindex, follow">'
);

CREATE TABLE page_type (
    page_type_id SERIAL,
    type_name VARCHAR(64) NOT NULL,
    substitution_types JSONB NOT NULL,
    template_url VARCHAR(64) NOT NULL,
    PRIMARY KEY (page_type_id)
);
INSERT INTO page_type (
    page_type_id, type_name, substitution_types, template_url
) VALUES
(
    1,
    'Blogpage',
    '
    {
        "PageURL": "URL", 
        "LangCode": "LangCode", 
        "PageTitle": "Text", 
        "LangTags": "LangTags",
        "Content": "Text",
        "LangRedirects": "LangRedirects"
    }
    '::JSONB,
    'http://nginx-frontend:8080/templates/blogpage.html'
);

CREATE TABLE pages (
    page_id SERIAL,
    page_type_id INT NOT NULL,
    FOREIGN KEY (page_type_id) REFERENCES page_type(page_type_id),
    PRIMARY KEY (page_id)
);
INSERT INTO pages (page_id, page_type_id) VALUES (1, 1);


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
INSERT INTO translations (page_id, lang_code, substitutions, url) VALUES
(1, 'en', '{ "PageTitle": "Page 1" }'::JSONB, '/blog/page1.html'),
(1, 'ja', '{ "PageTitle": "ページ１" }'::JSONB, '/ブログ/パージ１.html');

CREATE TABLE tests (
    test_id VARCHAR(2) NOT NULL,
    translation_id INT NOT NULL,
    test_substitutions JSONB NOT NULL,
    FOREIGN KEY (translation_id) REFERENCES translations(translation_id),
    PRIMARY KEY (test_id, translation_id)
);

CREATE INDEX idx_test_translationid ON tests (translation_id);

INSERT INTO tests (test_id, translation_id, test_substitutions) VALUES
('aa', 1, '{ "Content": "Hello world! Welcome to page 1!" }'::JSONB),
('ab', 1, '{ "Content": "Hello World! Welcome to Page 1!" }'::JSONB),
('ba', 1, '{ "Content": "Hello world!!! Welcome to page 1!!!" }'::JSONB),
('bb', 1, '{ "Content": "Hello World!!! Welcome to Page 1!!!" }'::JSONB),
('aa', 2, '{ "Content": "ハローワールド！ページ１えようこそ！" }'::JSONB),
('ab', 2, '{ "Content": "こんいちわ世界！ページ１えようこそ！" }'::JSONB);
