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

INSERT INTO translations (page_id, lang_code, substitutions, url) VALUES
(1, 'ja', '{ "PageTitle": "ログイン" }'::JSONB, '/ログイン.html'),
(2, 'ja', '{ "PageTitle": "登録" }'::JSONB, '/とうろく.html'),
(3, 'ja', '{ "PageTitle": "ユーザー" }'::JSONB, '/ユーザー.html'),
(
    4,
    'ja',
    '{ "PageTitle": "クリエイターダッシュボード" }'::JSONB,
    '/クリエイター/ダッシュボード.html'
),
(
    5,
    'ja',
    '{ "PageTitle": "エディタ" }'::JSONB,
    '/クリエイター/エディタ.html'
),
(6, 'ja', '{ "PageTitle": "ページ１" }'::JSONB, '/ブログ/ページ１.html');

INSERT INTO tests (test_id, translation_id, test_substitutions) VALUES
(
    '00',
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
    '00',
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
    '00',
    103,
    '
    {
        "ChangePassPrompt": "パスワードを替える：",
        "OldPassPrompt": "古いパスワード",
        "NewPassPrompt": "新しいパスワード",
        "ConfirmPassPrompt": "パスワード確認",
        "UpdatePassPrompt": "更新"
    }
    '::JSONB
),
(
    '00',
    104,
    '
    {
        "ManagePage": "ページ管理",
        "PermissionsPrompt": "許可",
        "SubmitPrompt": "新しいページを作る"
    }
    '::JSONB
),
(
    '00',
    105,
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
('00', 106, '{ "Content": "ハローワールド！ページ１えようこそ！" }'::JSONB),
('01', 106, '{ "Content": "こんいちわ世界！ページ１えようこそ！" }'::JSONB);
