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
INSERT INTO languages VALUES
('en', 'English', '', true),
(
    'ja',
    '日本語',
    '<meta name="robots" content="noindex, follow">',
    false
);

CREATE TABLE error_codes (
    code_id VARCHAR(16) NOT NULL,
    lang_code VARCHAR(2) NOT NULL,
    content VARCHAR(256) NOT NULL,
    FOREIGN KEY (lang_code) REFERENCES languages(lang_code),
    PRIMARY KEY (code_id, lang_code)
);
INSERT INTO error_codes VALUES
('AccExs', 'en', 'That account already exists!'),
('UmtchP', 'en', 'The passwords do not match!'),
('UWrong', 'en', 'Incorrect username!'),
('PWrong', 'en', 'Incorrect password!'),
('IntErr', 'en', 'Internal Server Error'),
('NoInfo', 'en', 'No analytics avaliable yet!'),
('SesExp', 'en', 'Your session expired!'),
('NoPerm', 'en', 'You do not have enough permissions'),
('AccExs', 'ja', 'そのアカウントは既にあります！'),
('UmtchP', 'ja', 'パスワードと確認は同じではありません！'),
('UWrong', 'ja', 'ユーザー名は正しくない！'),
('PWrong', 'ja', 'パスワードは正しくない！'),
('IntErr', 'ja', 'サーバーの中には問題いっぱいある！'),
('NoInfo', 'ja', 'このページの情報はありません。'),
('SesExp', 'ja', '君のセッションは消えました。'),
('NoPerm', 'ja', '許可がありません。');

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
    'LoginPage',
    '
    {
        "LangCode": "LangCode",
        "PageTitle": "Text",
        "LangTags": "LangTags",
        "Errors": "Errors",
        "UsernamePrompt": "Text",
        "PasswordPrompt": "Text",
        "ReturnURL": "ReturnURL",
        "SubmitPrompt": "Text",
        "LangRedirects": "LangRedirects",
        "SwitchPrompt": "TemplateText"
    }
    '::JSONB,
    'http://nginx-frontend:8080/templates/login.html'
),
(
    2,
    'SignupPage',
    '
    {
        "LangCode": "LangCode",
        "PageTitle": "Text",
        "LangTags": "LangTags",
        "Errors": "Errors",
        "UsernamePrompt": "Text",
        "PasswordPrompt": "Text",
        "ConfirmPassPrompt": "Text",
        "ReturnURL": "ReturnURL",
        "SubmitPrompt": "Text",
        "LangRedirects": "LangRedirects",
        "SwitchPrompt": "TemplateText"
    }
    '::JSONB,
    'http://nginx-frontend:8080/templates/signup.html'
),
(
    3,
    'BlogPage',
    '
    {
        "AccountDetails": "AccountDetails",
        "LangCode": "LangCode",
        "PageTitle": "Text",
        "LangTags": "LangTags",
        "Content": "Text",
        "LangRedirects": "LangRedirects"
    }
    '::JSONB,
    'http://nginx-frontend:8080/templates/blogpage.html'
),
(
    4,
    'CreatorDashboard',
    '
    {
        "PageTitle": "Text",
        "ManagePage": "Text",
        "PermissionsPrompt": "Text",
        "PageTypeDropdown": "Creator.PageTypeDropdown",
        "SubmitPrompt": "Text",
        "LangCode": "LangCode",
        "LangTags": "LangTags",
        "LangRedirects": "LangRedirects",
        "AccountDetails": "AccountDetails",
        "Dashboard": "Creator.Dashboard"
    }
    '::JSONB,
    'http://nginx-frontend:8080/templates/creator/dashboard.html'
),
(
    5,
    'PageEditor',
    '
    {
        "PageTitle": "Text",
        "LangCode": "LangCode",
        "LangTags": "LangTags",
        "LangRedirects": "LangRedirects",
        "AccountDetails": "AccountDetails"
    }
    '::JSONB,
    'http://nginx-frontend:8080/templates/creator/editor.html'
);

CREATE TABLE pages (
    page_id SERIAL,
    page_type_id INT NOT NULL,
    required_privilege INT NOT NULL,
    FOREIGN KEY (page_type_id) REFERENCES page_type(page_type_id),
    PRIMARY KEY (page_id)
);
INSERT INTO pages (page_type_id, required_privilege) VALUES
(1, 0),
(2, 0),
(3, 0),
-- Set to 6 latter. 0 for testing (Because signing in every time is a pain)
(4, 0),
(5, 0);


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
(1, 'en', '{ "PageTitle": "Login" }'::JSONB, '/login.html'),
(1, 'ja', '{ "PageTitle": "ログイン" }'::JSONB, '/ログイン.html'),
(2, 'en', '{ "PageTitle": "Signup" }'::JSONB, '/signup.html'),
(2, 'ja', '{ "PageTitle": "登録" }'::JSONB, '/とうろく.html'),
(3, 'en', '{ "PageTitle": "Page 1" }'::JSONB, '/blog/page1.html'),
(3, 'ja', '{ "PageTitle": "ページ１" }'::JSONB, '/ブログ/ページ１.html'),
(
    4,
    'en',
    '{ "PageTitle": "Creator Dashboard" }'::JSONB,
    '/creator/dashboard.html'
),
(
    4,
    'ja',
    '{ "PageTitle": "クリエイターダッシュボード" }'::JSONB,
    '/クリエイター/ダッシュボード.html'
),
(
    5,
    'en',
    '{ "PageTitle": "Editor" }'::JSONB,
    '/creator/editor.html'
),
(
    5,
    'ja',
    '{ "PageTitle": "エディタ" }'::JSONB,
    '/クリエイター/エディタ.html'
);

CREATE TABLE tests (
    test_id VARCHAR(2) NOT NULL,
    translation_id INT NOT NULL,
    test_substitutions JSONB NOT NULL,
    FOREIGN KEY (translation_id) REFERENCES translations(translation_id),
    PRIMARY KEY (test_id, translation_id)
);

CREATE INDEX idx_test_translationid ON tests (translation_id);

INSERT INTO tests (test_id, translation_id, test_substitutions) VALUES
(
    'aa',
    1,
    '
    {
        "UsernamePrompt": "Username",
        "PasswordPrompt": "Password",
        "SubmitPrompt": "Submit",
        "SwitchPrompt": "Don''t have an account? <a href=\"/en/signup.html?from={{ ReturnURL }}\">Sign Up</a> instead."
    }
    '::JSONB
),
(
    'aa',
    2,
    '
    {
        "UsernamePrompt": "ユーザー名",
        "PasswordPrompt": "パスワード",
        "SubmitPrompt": "ログイン",
        "SwitchPrompt": "アカウントをありませんならば、<a href=\"/ja/とうろく.html?from={{ ReturnURL }}\">登録</a>しませんか？"
    }
    '::JSONB
),
(
    'aa',
    3,
    '
    {
        "UsernamePrompt": "Username",
        "PasswordPrompt": "Password",
        "ConfirmPassPrompt": "Confirm Password",
        "SubmitPrompt": "Submit",
        "SwitchPrompt": "Already have an account? <a href=\"/en/login.html?from={{ ReturnURL }}\">Login</a> instead."
    }
    '::JSONB
),
(
    'aa',
    4,
    '
    {
        "UsernamePrompt": "ユーザー名",
        "PasswordPrompt": "パスワード",
        "ConfirmPassPrompt": "パスワード確認",
        "SubmitPrompt": "登録",
        "SwitchPrompt": "既にアカウントをありますか？<a href=\"/ja/ログイン.html?from={{ ReturnURL }}\">ログイン</a>してください。"
    }
    '::JSONB
),
('aa', 5, '{ "Content": "Hello world! Welcome to page 1!" }'::JSONB),
('ab', 5, '{ "Content": "Hello World! Welcome to Page 1!" }'::JSONB),
('ba', 5, '{ "Content": "Hello world!!! Welcome to page 1!!!" }'::JSONB),
('bb', 5, '{ "Content": "Hello World!!! Welcome to Page 1!!!" }'::JSONB),
('aa', 6, '{ "Content": "ハローワールド！ページ１えようこそ！" }'::JSONB),
('ab', 6, '{ "Content": "こんいちわ世界！ページ１えようこそ！" }'::JSONB),
(
    'aa',
    7,
    '
    {
        "ManagePage": "Manage Page",
        "PermissionsPrompt": "Permissions",
        "SubmitPrompt": "Create New Page"
    }
    '::JSONB
),
(
    'aa',
    8,
    '
    {
        "ManagePage": "ページ管理",
        "PermissionsPrompt": "許可",
        "SubmitPrompt": "新しいページを作る"
    }
    '::JSONB
),
('aa', 9, '{ }'::JSONB),
('aa', 10, '{ }'::JSONB);
