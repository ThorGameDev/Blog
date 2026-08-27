INSERT INTO page_type (
    page_type_id, type_name, substitution_types, template_url
) VALUES
(
    1,
    'BlogPage',
    '
    {
        "PageTitle": "Title",
        "LangCode": "LangCode",
        "LangTags": "LangTags",
        "LangRedirects": "LangRedirects",
        "AccountDetails": "AccountDetails",
        "Content": "Content"
    }
    '::JSONB,
    'http://nginx-frontend:8080/templates/blogpage.html'
),
(
    100,
    'LoginPage',
    '
    {
        "PageTitle": "Title",
        "LangCode": "LangCode",
        "LangTags": "LangTags",
        "LangRedirects": "LangRedirects",
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
    'SignupPage',
    '
    {
        "PageTitle": "Title",
        "LangCode": "LangCode",
        "LangTags": "LangTags",
        "LangRedirects": "LangRedirects",
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
    'UserAccountPage',
    '
    {
        "PageTitle": "Title",
        "LangCode": "LangCode",
        "LangTags": "LangTags",
        "LangRedirects": "LangRedirects",
        "ChangePassPrompt": "Text",
        "OldPassPrompt": "Text",
        "NewPassPrompt": "Text",
        "ConfirmPassPrompt": "Text",
        "UpdatePassPrompt": "Text",
        "AccountDetails": "User.AccountDetails"
    }
    '::JSONB,
    'http://nginx-frontend:8080/templates/user.html'
),
(
    200,
    'CreatorDashboard',
    '
    {
        "PageTitle": "Title",
        "LangCode": "LangCode",
        "LangTags": "LangTags",
        "LangRedirects": "LangRedirects",
        "AccountDetails": "AccountDetails",
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
    'PageEditor',
    '
    {
        "PageTitle": "Title",
        "LangCode": "LangCode",
        "LangTags": "LangTags",
        "LangRedirects": "LangRedirects",
        "AccountDetails": "AccountDetails",
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
-- Set to 6 latter. 0 for testing (Because signing in every time is a pain)
(200, 0),
(201, 0),
-- Eventually needs removed. Needed for testing default pages
(1, 0);
