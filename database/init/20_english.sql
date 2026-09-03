INSERT INTO languages (lang_code, lang_name, page_tags, is_primary) VALUES
('en', 'English', '', true);

INSERT INTO error_codes (code_id, lang_code, content) VALUES
('AccExs', 'en', 'That account already exists!'),
('UmtchP', 'en', 'The passwords do not match!'),
('UWrong', 'en', 'Incorrect username!'),
('PWrong', 'en', 'Incorrect password!'),
('IntErr', 'en', 'Internal Server Error'),
('NoInfo', 'en', 'No analytics avaliable yet!'),
('SesExp', 'en', 'Your session expired!'),
('NoPerm', 'en', 'You do not have enough permissions');

INSERT INTO sitewide_tests (test_id, lang_code, substitutions) VALUES
(
    '00000001',
    'en',
    '
    {
        "Global.Replies": " Replies"
    }
    '::JSONB
);

INSERT INTO translations (page_id, lang_code, title, url) VALUES
(1, 'en', 'Login', '/login.html'),
(2, 'en', 'Signup', '/signup.html'),
(3, 'en', 'User', '/user.html'),
(4, 'en', 'Comment', '/comment.html'),
(
    5,
    'en',
    'Creator Dashboard',
    '/creator/dashboard.html'
),
(
    6,
    'en',
    'Editor',
    '/creator/editor.html'
),
(7, 'en', 'Page 1', '/blog/page1.html');

INSERT INTO tests (test_id, translation_id, substitutions) VALUES
(
    '01',
    1,
    '
    {
        "UsernamePrompt": "Username",
        "PasswordPrompt": "Password",
        "SubmitPrompt": "Login",
        "SwitchPrompt": "Don''t have an account? <a href=\"/en/signup.html?from={{ ReturnURL }}\">Sign Up</a> instead."
    }
    '::JSONB
),
(
    '01',
    2,
    '
    {
        "UsernamePrompt": "Username",
        "PasswordPrompt": "Password",
        "ConfirmPassPrompt": "Confirm Password",
        "SubmitPrompt": "Sign Up",
        "SwitchPrompt": "Already have an account? <a href=\"/en/login.html?from={{ ReturnURL }}\">Login</a> instead."
    }
    '::JSONB
),
(
    '01',
    3,
    '
    {
        "ChangePassPrompt": "Change Password:",
        "OldPassPrompt": "Old Password",
        "NewPassPrompt": "New Password",
        "ConfirmPassPrompt": "Confirm New Password",
        "UpdatePassPrompt": "Update",
        "SelectPFPPrompt": "Change Profile Picture",
        "UpdatePFPPrompt": "Update"
    }
    '::JSONB
),
('01', 4, '{ }'::JSONB),
(
    '01',
    5,
    '
    {
        "ManagePage": "Manage Page",
        "PermissionsPrompt": "Permissions",
        "SubmitPrompt": "Create New Page"
    }
    '::JSONB
),
(
    '01',
    6,
    '
    {
        "SaveChangesPrompt": "Save Changes",
        "AddTranslationPrompt": "Add Translation",
        "NewTranslationUrlPrompt": "URL",
        "NewTranslationTitlePrompt": "Title",
        "AddTestPrompt": "Add Test"
    }
    '::JSONB
),
('01', 7, '{ "Content": "Hello world! Welcome to page 1!" }'::JSONB),
('02', 7, '{ "Content": "Hello World! Welcome to Page 1!" }'::JSONB),
('03', 7, '{ "Content": "Hello world!!! Welcome to page 1!!!" }'::JSONB),
('04', 7, '{ "Content": "Hello World!!! Welcome to Page 1!!!" }'::JSONB);

-- Meaningless temporary test data

-- Create users (That can not be logged into, due to an impossible hash)
INSERT INTO users (username, password_hash, pfp_id, privilege) VALUES
('TestUser', 'Password Hash', 1, 1),
('FakePerson', 'Password Hash', 1, 1);

INSERT INTO comments (translation_id, uid, container_id, content) VALUES
(7, 1, NULL, 'First!'),
(7, 2, NULL, 'first'),
(7, 1, 2, 'Haha, Sucker!'),
(7, 2, 3, 'Shutup.'),
(7, 1, NULL, 'This is a verry good comment!'),
(7, 2, 5, 'It''s not'),
(7, 2, 5, 'It''s really not.'),
(7, 2, 5, 'Not at all.'),
(7, 1, 8, 'No need to say it so many times'),
(7, 2, 9, 'But I do though'),
(7, 1, 10, 'Why?'),
(7, 2, 11, 'Because its how it should be'),
(7, 1, 12, 'You should be quiet.');
