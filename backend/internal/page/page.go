package page

import (
	"blogbackend/internal/page/page_parts"
	"blogbackend/internal/utils/db"
	"blogbackend/internal/utils/utils_err"
	"blogbackend/internal/utils/utils_fs"
	"blogbackend/internal/utils/utils_sec"
	"blogbackend/internal/utils/utils_url"
	"context"
	"errors"
	"fmt"
	"html"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/valyala/fasttemplate"
)

func getLangTags(domainURL string, currentLangCode string, alternatePages []utils_url.LangURL) (string, error) {
	var langTagCluster strings.Builder
	for _, val := range alternatePages {
		toURL := domainURL + "/" + val.LangCode + val.PageURL
		fmt.Fprintf(&langTagCluster, `<link rel="alternate" hreflang="%s" href="%s">`, val.LangCode, toURL)
		if val.IsPrimary {
			fmt.Fprintf(&langTagCluster, `<link rel="alternate" hreflang="x-default" href="%s">`, toURL)
		}
	}

	var langTags string
	err := db.Pool.QueryRow(context.Background(),
		`SELECT page_tags FROM languages WHERE lang_code = $1`,
		currentLangCode).Scan(&langTags)

	langTagCluster.WriteString(langTags)

	return langTagCluster.String(), err
}

func getLangLinks(currentLangCode string, alternatePages []utils_url.LangURL) (string, error) {
	var langLinks strings.Builder
	for _, val := range alternatePages {
		if currentLangCode != val.LangCode {
			toURL := "/" + val.LangCode + val.PageURL + val.QueryParams
			fmt.Fprintf(&langLinks, `<a rel="alternate" hreflang="%s" href="%s">%s</a>`, val.LangCode, toURL, val.LangName)
		}
	}
	return langLinks.String(), nil
}

func badURLRedirect(w http.ResponseWriter, req *http.Request, fromPage string, queryParams url.Values, langCode string) {
	newURL := utils_url.TranslateURL("/##"+fromPage, queryParams, langCode)
	http.Redirect(w, req, newURL, http.StatusPermanentRedirect)
}

func createLoginLinks(fromPageUrl string, langCode string) string {
	queryParamData := url.Values{}
	queryParamData.Set("from", "/"+langCode+fromPageUrl)
	queryParams := queryParamData.Encode()
	var loginURL string
	var loginTitle string
	var signupURL string
	var signupTitle string

	const sql_query string = `SELECT url, title
		FROM translations
		WHERE page_id = (
			SELECT page_id FROM translations
				WHERE url = $1
				AND lang_code = $2
		)
		AND lang_code = $3`
	err := db.Pool.QueryRow(context.Background(), sql_query, "/login.html", "en", langCode).Scan(&loginURL, &loginTitle)
	if err != nil {
		slog.Error("Failed to get url and title of login", "err", err)
		return ""
	}
	err = db.Pool.QueryRow(context.Background(), sql_query, "/signup.html", "en", langCode).Scan(&signupURL, &signupTitle)
	if err != nil {
		slog.Error("Failed to get url and title of signup", "err", err)
		return ""
	}

	return fmt.Sprintf(`<a href="/%s%s?%s">%s</a><a href="/%s%s?%s">%s</a>`,
		langCode, loginURL, queryParams, loginTitle,
		langCode, signupURL, queryParams, signupTitle,
	)
}

func getAccountDetails(req *http.Request, pageURL string, langCode string) (string, error) {
	session_id, err := req.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			return createLoginLinks(pageURL, langCode), nil
		}
		slog.Error("Failed to get session ID cookie!", "err", err)
		return "", err
	}
	var username string
	var pfp_file_id string
	err = db.Pool.QueryRow(context.Background(),
		`SELECT username, pfp_file_id FROM users, sessions
			WHERE users.uid = sessions.uid
			AND sessions.session_token = $1`,
		session_id.Value).Scan(&username, &pfp_file_id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return createLoginLinks(pageURL, langCode), nil
		}
		slog.Error("Error while fetching user information", "err", err)
		return "", err
	}
	accountPageURL := utils_url.TranslateURL("/en/user.html", nil, langCode)
	return fmt.Sprintf(`<a href="%s"><h3>%s</h3><img src="%s"></a>`, accountPageURL, username, pfp_file_id), nil
}

func pageGen(w http.ResponseWriter, req *http.Request) {
	domainURL := utils_url.GetRequestURL(req)

	fullPageURL := req.URL.Path
	langCode := fullPageURL[1:3]
	pageURL := fullPageURL[3:]
	queryParams := req.URL.Query()

	// Get page typeset
	var substitutionTypes map[string]string
	var templateUrl string
	var requiredPrivilege int
	var title string
	var translationId int
	err := db.Pool.QueryRow(context.Background(),
		`SELECT substitution_types, template_url, required_privilege, title, translation_id
			FROM page_type, pages, translations
			WHERE pages.page_type_id = page_type.page_type_id
			AND translations.page_id = pages.page_id
			AND lang_code = $1
			AND url = $2`,
		langCode, pageURL).Scan(&substitutionTypes, &templateUrl, &requiredPrivilege, &title, &translationId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			badURLRedirect(w, req, pageURL, queryParams, langCode)
		} else {
			slog.Error("Critical SQL error while searching for page", "err", err)
			http.Error(w, "Could not find page in SQL", http.StatusNotFound)
		}
		return
	}

	// Check if user has the permissions required to view the page
	if requiredPrivilege > 0 {
		sessionId, err := req.Cookie("session_id")
		if err != nil {
			if err != http.ErrNoCookie {
				slog.Error("Unknown cookie error", "err", err)
			}
			http.Error(w, "No permissions to access page", http.StatusForbidden)
			return
		}
		err = utils_sec.RequireLevel(sessionId.Value, requiredPrivilege)
		if err != nil {
			http.Error(w, utils_err.CodeToMessage(err.Error(), langCode), http.StatusForbidden)
			return
		}
	}

	// Get substitutions
	var substitutions map[string]string
	err = db.Pool.QueryRow(context.Background(),
		`SELECT test_substitutions FROM tests
			WHERE translation_id = $1
			ORDER BY (test_id = $2) DESC`,
		translationId, queryParams.Get("test")).Scan(&substitutions)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			slog.Error("Critical SQL error while getting page content", "err", err)
		}
		http.Error(w, "Could not find page in SQL", http.StatusNotFound)
		return
	}

	altURLs, err := utils_url.GetAlternateURLs(pageURL, langCode, queryParams)
	if err != nil {
		slog.Error("Failure getting alternate URLs!", "err", err)
		http.Error(w, "Failure getting alternate URLs", http.StatusNotFound)
		return
	}

	finalSubstitutions := make(map[string]interface{})
	for key, val := range substitutionTypes {
		switch val {
		case "URL":
			finalSubstitutions[key] = url.PathEscape(fullPageURL)
		case "Title":
			finalSubstitutions[key] = title
		case "LangCode":
			finalSubstitutions[key] = langCode
		case "LangTags":
			langTags, err := getLangTags(domainURL, langCode, altURLs)
			if err != nil {
				slog.Error("Error getting language tags", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			finalSubstitutions[key] = langTags
		case "LangRedirects":
			langLinks, err := getLangLinks(langCode, altURLs)
			if err != nil {
				slog.Error("Error getting language links", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			finalSubstitutions[key] = langLinks
		case "Errors":
			errorURL := req.URL.Query().Get("err")
			if errorURL != "" {
				errorMessage := utils_err.CodeToMessage(errorURL, langCode)
				finalSubstitutions[key] = html.EscapeString(errorMessage)
			} else {
				finalSubstitutions[key] = ""
			}
		case "AccountDetails":
			details, err := getAccountDetails(req, pageURL, langCode)
			if err != nil {
				slog.Error("Error getting Account Details", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			finalSubstitutions[key] = details
		case "ReturnURL":
			returnURL := utils_url.SanitizeURL(req.URL.Query().Get("from"))
			htmlEscaped := html.EscapeString(returnURL)
			finalSubstitutions[key] = htmlEscaped
		case "Text", "TemplateText", "Content":
			if substitutionValue, ok := substitutions[key]; ok {
				finalSubstitutions[key] = substitutionValue
			} else {
				slog.Warn("Untranslated content in", "url", fullPageURL, "key", key)
				finalSubstitutions[key] = ""
			}
		case "Creator.Dashboard":
			dashboardData, err := page_parts.GenerateCreatorDashboard(langCode)
			if err != nil {
				slog.Error("Error while creating creator dashboard", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			finalSubstitutions[key] = dashboardData
		case "Creator.PageTypeDropdown":
			finalSubstitutions[key] = page_parts.GeneratePageTypeDropdown()
		case "Creator.Editor":
			finalSubstitutions[key] = page_parts.GenerateEditor(langCode, req.URL.Query().Get("page"))
		case "User.AccountDetails":
			details, err := page_parts.GenerateUserPage(req, pageURL, langCode)
			if err != nil {
				slog.Error("Error getting Account Details", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			finalSubstitutions[key] = details
		}
	}

	// Resolve template texts
	for key, val := range substitutionTypes {
		switch val {
		case "TemplateText", "Creator.Dashboard", "Creator.Editor":
			subTemplate := fasttemplate.New(finalSubstitutions[key].(string), "{{ ", " }}")
			finalSubstitutions[key] = subTemplate.ExecuteString(finalSubstitutions)
		}
	}

	// Get base page
	pageData, err := utils_fs.RetrievePage(templateUrl)
	if err != nil {
		slog.Error("Error fetching page", "url", templateUrl, "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Parse base page as html template
	//pageTemplate, err := template.New(pageURL).Parse(pageData)
	pageTemplate := fasttemplate.New(pageData, "{{ ", " }}")
	htmlData := pageTemplate.ExecuteString(finalSubstitutions)

	// Return
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, htmlData)
}

func Register() {
	// Get a list of URL entrypoints
	rows, err := db.Pool.Query(context.Background(), "SELECT lang_code FROM languages;")
	if err != nil {
		slog.Error("Error while querying the database", "err", err)
	}
	defer rows.Close()

	// Register an api call for each. ie /en/blog/ and /ja/ブログ/
	var retrieved string
	for rows.Next() {
		if err := rows.Scan(&retrieved); err != nil {
			slog.Error("Critical Error!", "err", err)
		}
		http.HandleFunc("/"+retrieved+"/", pageGen)
	}
}
