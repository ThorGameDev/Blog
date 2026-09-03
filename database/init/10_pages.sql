INSERT INTO page_types (
    page_type_id, parent_page_type_id, type_name, substitution_types, template_url
) VALUES
(
    -1,
    NULL,
    'Global',
    '
    {
        "PageTitle": "Title",
        "LangCode": "LangCode",
        "LangTags": "LangTags",
        "AccountDetails": "AccountDetails",
        "LangRedirects": "LangRedirects"
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

INSERT INTO pages (page_type_id, required_privilege) VALUES
(100, 0),
(101, 0),
(102, 0),
(103, 0),
-- Set to 6 latter. 0 for testing (Because signing in every time is a pain)
(200, 0),
(201, 0),
-- Eventually needs removed. Needed for testing default pages
(1, 0);
