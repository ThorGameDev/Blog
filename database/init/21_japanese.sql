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

-- Ensure that the ordering of the translations is independent of the English order for initialization
ALTER SEQUENCE translations_translation_id_seq RESTART WITH 101;

INSERT INTO translations (page_id, lang_code, title, url) VALUES
(1, 'ja', 'ログイン', '/ログイン.html'),
(2, 'ja', '登録', '/とうろく.html'),
(3, 'ja', 'ユーザー', '/ユーザー.html'),
(4, 'ja', 'コメント', '/コメント.html'),
(
    5,
    'ja',
    'クリエイターダッシュボード',
    '/クリエイター/ダッシュボード.html'
),
(
    6,
    'ja',
    'エディタ',
    '/クリエイター/エディタ.html'
),
(7, 'ja', 'ページ１', '/ブログ/ページ１.html');

INSERT INTO tests (test_id, translation_id, substitutions) VALUES
(
    '01',
    101,
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
    '01',
    102,
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
    '01',
    103,
    '
    {
        "ChangePassPrompt": "パスワードを替える：",
        "OldPassPrompt": "古いパスワード",
        "NewPassPrompt": "新しいパスワード",
        "ConfirmPassPrompt": "パスワード確認",
        "UpdatePassPrompt": "更新",
        "SelectPFPPrompt": "写真を替える",
        "UpdatePFPPrompt": "更新"
    }
    '::JSONB
),
('01', 104, '{ }'::JSONB),
(
    '01',
    105,
    '
    {
        "ManagePage": "ページ管理",
        "PermissionsPrompt": "許可",
        "SubmitPrompt": "新しいページを作る"
    }
    '::JSONB
),
(
    '01',
    106,
    '
    {
        "SaveChangesPrompt": "保存",
        "AddTranslationPrompt": "新しい翻訳",
        "NewTranslationUrlPrompt": "URL",
        "NewTranslationTitlePrompt": "ページの名",
        "AddTestPrompt": "新しいテスト"
    }
    '::JSONB
),
('01', 107, '{ "Content": "ハローワールド！ページ１えようこそ！" }'::JSONB),
('02', 107, '{ "Content": "こんいちわ世界！ページ１えようこそ！" }'::JSONB);
