INSERT INTO page_types (
    page_type_id, parent_page_type_id, type_name, substitution_types, template_url
) VALUES
(
    -1,
    NULL,
    'Global',
    '
    {
        "Global.BlogTitle": "Text",
        "Global.AGPLv3License": "Text",
        "Global.SourceCode": "Text",
        "TopBar": "TopBar",
        "BottomBar": "BottomBar",
        "PageTitle": "Title",
        "LangCode": "LangCode",
        "LangTags": "LangTags"
    }
    '::JSONB,
    NULL
),
(
    -2,
    -1,
    'HasComments',
    '
    {
        "Global.CommentSectionHeader": "Text",
        "Global.SubmitComment": "Text",
        "Global.LoginToComment": "Text",
        "Global.Replies": "Text",
        "Global.SubmitReply": "Text"
    }
    '::JSONB,
    NULL
),
(
    1,
    -2,
    'BlogPage',
    '
    {
        "Content": "Content",
        "CommentSection": "CommentSection"
    }
    '::JSONB,
    'http://nginx-frontend:8080/templates/blogpage.html'
),
(
    2,
    -1,
    'InfoPage',
    '
    {
        "Content": "Content"
    }
    '::JSONB,
    'http://nginx-frontend:8080/templates/infopage.html'
),
(
    100,
    -1,
    'LoginPage',
    '
    {
        "Errors": "Errors",
        "ReturnURL": "ReturnURL",
        "UsernamePrompt": "Text",
        "PasswordPrompt": "Text",
        "SubmitPrompt": "Text",
        "SwitchPrompt": "TemplateText"
    }
    '::JSONB,
    'http://nginx-frontend:8080/templates/login.html'
),
(
    101,
    -1,
    'SignupPage',
    '
    {
        "Errors": "Errors",
        "ReturnURL": "ReturnURL",
        "UsernamePrompt": "Text",
        "PasswordPrompt": "Text",
        "ConfirmPassPrompt": "Text",
        "SubmitPrompt": "Text",
        "SwitchPrompt": "TemplateText"
    }
    '::JSONB,
    'http://nginx-frontend:8080/templates/signup.html'
),
(
    102,
    -1,
    'UserAccountPage',
    '
    {
        "ChangePassPrompt": "Text",
        "OldPassPrompt": "Text",
        "NewPassPrompt": "Text",
        "ConfirmPassPrompt": "Text",
        "UpdatePassPrompt": "Text",
        "SelectPFPPrompt": "Text",
        "UpdatePFPPrompt": "Text",
        "AccountDetails": "User.AccountDetails"
    }
    '::JSONB,
    'http://nginx-frontend:8080/templates/user.html'
),
(
    103,
    -2,
    'Comments',
    '
    {
        "Comment": "Comment"
    }
    '::JSONB,
    'http://nginx-frontend:8080/templates/comment.html'
),
(
    200,
    -1,
    'CreatorDashboard',
    '
    {
        "ManagePage": "Text",
        "PermissionsPrompt": "Text",
        "SubmitPrompt": "Text",
        "PageTypeDropdown": "Creator.PageTypeDropdown",
        "Dashboard": "Creator.Dashboard"
    }
    '::JSONB,
    'http://nginx-frontend:8080/templates/creator/dashboard.html'
),
(
    201,
    -1,
    'PageEditor',
    '
    {
        "SaveChangesPrompt": "Text",
        "AddTranslationPrompt": "Text",
        "NewTranslationUrlPrompt": "Text",
        "NewTranslationTitlePrompt": "Text",
        "AddTestPrompt": "Text",
        "Editor": "Creator.Editor"
    }
    '::JSONB,
    'http://nginx-frontend:8080/templates/creator/editor.html'
);

INSERT INTO pages (page_type_id, required_privilege, page_index) VALUES
(100, 0, B'000'),
(101, 0, B'000'),
(102, 0, B'000'),
(103, 0, B'000'),
-- Set to 6 latter. 0 for testing (Because signing in every time is a pain)
(200, 0, B'010'), -- Shown in top bar for testing
(201, 0, B'000'),
-- Eventually needs removed. Needed for testing default pages
(1, 0, B'001'),
(2, 0, B'111'); -- about.html
