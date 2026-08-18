INSERT INTO page_type (
    page_type_id, type_name, substitution_types, template_url
) VALUES
(
    1,
    'BlogPage',
    '
    {
        "PageTitle": "Text",
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
        "PageTitle": "Text",
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
        "PageTitle": "Text",
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
    200,
    'CreatorDashboard',
    '
    {
        "PageTitle": "Text",
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
        "PageTitle": "Text",
        "LangCode": "LangCode",
        "LangTags": "LangTags",
        "LangRedirects": "LangRedirects",
        "AccountDetails": "AccountDetails",
        "SaveChangesPrompt": "Text",
        "AddTestPrompt": "Text",
        "Editor": "Creator.Editor"
    }
    '::JSONB,
    'http://nginx-frontend:8080/templates/creator/editor.html'
);

INSERT INTO pages (page_type_id, required_privilege) VALUES
(100, 0),
(101, 0),
-- Set to 6 latter. 0 for testing (Because signing in every time is a pain)
(200, 0),
(201, 0),
-- Eventually needs removed. Needed for testing default pages
(1, 0);
