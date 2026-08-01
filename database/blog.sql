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
    expire_date TIMESTAMP NOT NULL,
    FOREIGN KEY (uid) REFERENCES users(uid),
    PRIMARY KEY (session_token)
);

CREATE TABLE pages (
    pageid SERIAL,
    PRIMARY KEY (pageid)
);
INSERT INTO pages (pageid) VALUES (1);

CREATE TABLE languages (
    langcode VARCHAR(2) NOT NULL,
    langname VARCHAR(32) NOT NULL,
    pagetags TEXT NOT NULL,
    urlbase VARCHAR(32) NOT NULL,
    PRIMARY KEY (langcode)
);
INSERT INTO languages VALUES
('en', 'English', '', '/en/blog/'),
(
    'ja',
    '日本語',
    '<meta name="robots" content="noindex, follow">',
    '/ja/ブログ/'
);

CREATE TABLE translations (
    translationid SERIAL,
    pageid INT NOT NULL,
    langcode VARCHAR(2) NOT NULL,
    title VARCHAR(64) NOT NULL,
    url VARCHAR(32) NOT NULL,
    UNIQUE (pageid, langcode),
    UNIQUE (url),
    FOREIGN KEY (pageid) REFERENCES pages(pageid),
    FOREIGN KEY (langcode) REFERENCES languages(langcode),
    PRIMARY KEY (translationid)
);
INSERT INTO translations (pageid, langcode, title, url) VALUES
(1, 'en', 'Page 1', 'page1.html'),
(1, 'ja', 'ページ１', 'パージ１.html');

CREATE TABLE content (
    testid VARCHAR(2) NOT NULL,
    translationid INT NOT NULL,
    content TEXT NOT NULL,
    FOREIGN KEY (translationid) REFERENCES translations(translationid),
    PRIMARY KEY (testid, translationid)
);

CREATE INDEX idx_content_translationid ON content (translationid);

INSERT INTO content (testid, translationid, content) VALUES
('aa', 1, 'Hello world! Welcome to page 1!'),
('ab', 1, 'Hello World! Welcome to Page 1!'),
('ba', 1, 'Hello world!!! Welcome to page 1!!!'),
('bb', 1, 'Hello World!!! Welcome to Page 1!!!'),
('aa', 2, 'ハローワールド！ページ１えようこそ！'),
('ab', 2, 'こんいちわ世界！ページ１えようこそ！');
