INSERT INTO languages VALUES
(
    'ja',
    '日本語',
    '<meta name="robots" content="noindex, follow">',
    false
);

INSERT INTO error_codes VALUES
('AccExs', 'ja', 'そのアカウントは既にあります！'),
('UmtchP', 'ja', 'パスワードと確認は同じではありません！'),
('UWrong', 'ja', 'ユーザー名は正しくない！'),
('PWrong', 'ja', 'パスワードは正しくない！'),
('IntErr', 'ja', 'サーバーの中には問題いっぱいある！'),
('NoInfo', 'ja', 'このページの情報はありません。'),
('SesExp', 'ja', '君のセッションは消えました。'),
('NoPerm', 'ja', '許可がありません。');

INSERT INTO translations (page_id, lang_code, substitutions, url) VALUES
(1, 'ja', '{ "PageTitle": "ログイン" }'::JSONB, '/ログイン.html'),
(2, 'ja', '{ "PageTitle": "登録" }'::JSONB, '/とうろく.html'),
(
    3,
    'ja',
    '{ "PageTitle": "クリエイターダッシュボード" }'::JSONB,
    '/クリエイター/ダッシュボード.html'
),
(
    4,
    'ja',
    '{ "PageTitle": "エディタ" }'::JSONB,
    '/クリエイター/エディタ.html'
),
(5, 'ja', '{ "PageTitle": "ページ１" }'::JSONB, '/ブログ/ページ１.html');

INSERT INTO tests (test_id, translation_id, test_substitutions) VALUES
(
    'aa',
    6,
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
    7,
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
('aa', 10, '{ "Content": "ハローワールド！ページ１えようこそ！" }'::JSONB),
('ab', 10, '{ "Content": "こんいちわ世界！ページ１えようこそ！" }'::JSONB);
