INSERT INTO languages VALUES
('en', 'English', '', true);

INSERT INTO error_codes VALUES
('AccExs', 'en', 'That account already exists!'),
('UmtchP', 'en', 'The passwords do not match!'),
('UWrong', 'en', 'Incorrect username!'),
('PWrong', 'en', 'Incorrect password!'),
('IntErr', 'en', 'Internal Server Error'),
('NoInfo', 'en', 'No analytics avaliable yet!'),
('SesExp', 'en', 'Your session expired!'),
('NoPerm', 'en', 'You do not have enough permissions');

INSERT INTO translations (page_id, lang_code, title, url) VALUES
(1, 'en', 'Login', '/login.html'),
(2, 'en', 'Signup', '/signup.html'),
(3, 'en', 'User', '/user.html'),
(
    4,
    'en',
    'Creator Dashboard',
    '/creator/dashboard.html'
),
(
    5,
    'en',
    'Editor',
    '/creator/editor.html'
),
(6, 'en', 'Page 1', '/blog/page1.html');

INSERT INTO tests (test_id, translation_id, substitutions) VALUES
(
    '00',
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
    '00',
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
    '00',
    3,
    '
    {
        "ChangePassPrompt": "Change Password:",
        "OldPassPrompt": "Old Password",
        "NewPassPrompt": "New Password",
        "ConfirmPassPrompt": "Confirm New Password",
        "UpdatePassPrompt": "Update"
    }
    '::JSONB
),
(
    '00',
    4,
    '
    {
        "ManagePage": "Manage Page",
        "PermissionsPrompt": "Permissions",
        "SubmitPrompt": "Create New Page"
    }
    '::JSONB
),
(
    '00',
    5,
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
('00', 6, '{ "Content": "Hello world! Welcome to page 1!" }'::JSONB),
('01', 6, '{ "Content": "Hello World! Welcome to Page 1!" }'::JSONB),
('02', 6, '{ "Content": "Hello world!!! Welcome to page 1!!!" }'::JSONB),
('03', 6, '{ "Content": "Hello World!!! Welcome to Page 1!!!" }'::JSONB);
