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

func badURLRedirect(w http.ResponseWriter, req *http.Request, fromPage string, queryParams url.Values, langCode string) {
	newURL := utils_url.TranslateURL("/##"+fromPage, queryParams, langCode)
	http.Redirect(w, req, newURL, http.StatusPermanentRedirect)
}

func getRandomTest(translationId int) string {
	var randTest string
	err := db.Pool.QueryRow(context.Background(),
		`SELECT test_id FROM tests
			WHERE translation_id = $1
			ORDER BY RANDOM()
			LIMIT 1`,
		translationId).Scan(&randTest)
	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Error("No active tests for translation", "translationId", translationId)
		} else {
			slog.Error("SQL error while getting random test", "translationId", translationId, "err", err)
		}
		return "00"
	}
	return randTest
}

func isTestActive(translationId int, testId string) bool {
	// exit early if no tests were decided
	if testId == "00" {
		return false
	}

	var exists bool
	err := db.Pool.QueryRow(context.Background(),
		`SELECT EXISTS (
			SELECT 1 FROM tests
				WHERE translation_id = $1
				AND test_id = $2
		)`,
		translationId, testId).Scan(&exists)
	if err != nil {
		slog.Error("Unknown error while checking if test exists!", "err", err)
		return false
	}

	return exists
}

const PAGES_PER_COOKIE = 64

func getABver(w http.ResponseWriter, req *http.Request, translationId int, langCode string) (string, string) {
	// If there is a test query parameter, display that test simply
	if req.URL.Query().Has("test") {
		return req.URL.Query().Get("test"), "00000001"
	}

	cookieId := (translationId - 1) / PAGES_PER_COOKIE            // Which segment of cookie data to use
	cookieName := fmt.Sprintf("testId-%s-%d", langCode, cookieId) // Which cookie to use
	cookieIndex := ((translationId - 1) % PAGES_PER_COOKIE) * 2   // Which character in the cookie to use

	// Extract correct cookie
	var cookieStr string
	testId := "00"
	cookieData, err := req.Cookie(cookieName)
	if err != nil {
		// Do return a random test on error
		if err != http.ErrNoCookie {
			slog.Error("Unknown cookie error", "err", err)
			return getRandomTest(translationId), "00000000"
		}
		// If there is no cookie, create an empty one
		cookieStr = strings.Repeat("0", PAGES_PER_COOKIE*2)
	} else {
		// If the found cookie is not a valid size, make a new one. Otherwise, extract data from it
		if len(cookieData.Value) < cookieIndex+2 {
			cookieStr = strings.Repeat("0", PAGES_PER_COOKIE*2)
		} else {
			cookieStr = cookieData.Value
			testId = cookieStr[cookieIndex : cookieIndex+2]
		}
	}

	// If the cookie is not set, set it
	if !isTestActive(translationId, testId) {
		testId = getRandomTest(translationId)
		cookieStr = cookieStr[:cookieIndex] + testId + cookieStr[cookieIndex+2:]

		sessionCookie := &http.Cookie{
			Name:     cookieName,
			Value:    cookieStr,
			Path:     "/" + langCode + "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, sessionCookie)
	}

	return testId, "00000000"
}

func pageGen(w http.ResponseWriter, req *http.Request) {
	fullPageURL := req.URL.Path
	langCode := fullPageURL[1:3]
	pageURL := fullPageURL[3:]
	queryParams := req.URL.Query()

	// Get basic details on the page
	var pageTypeId int
	var templateUrl string
	var requiredPrivilege int
	var title string
	var translationId int
	err := db.Pool.QueryRow(context.Background(),
		`SELECT page_types.page_type_id, template_url, required_privilege, title, translation_id
			FROM page_types, pages, translations
			WHERE pages.page_type_id = page_types.page_type_id
			AND translations.page_id = pages.page_id
			AND lang_code = $1
			AND url = $2`,
		langCode, pageURL).Scan(&pageTypeId, &templateUrl, &requiredPrivilege, &title, &translationId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			badURLRedirect(w, req, pageURL, queryParams, langCode)
		} else {
			slog.Error("Critical SQL error while searching for page", "err", err)
			http.Error(w, "Could not find page in SQL", http.StatusNotFound)
		}
		return
	}

	// Get page typeset
	var substitutionTypes map[string]string
	err = db.Pool.QueryRow(context.Background(),
		`WITH RECURSIVE chain AS (
			SELECT page_type_id, parent_page_type_id, substitution_types as merged
				FROM page_types
				WHERE page_type_id = $1

			UNION ALL

			SELECT p.page_type_id, p.parent_page_type_id, p.substitution_types || c.merged
				FROM page_types p
				JOIN chain c ON p.page_type_id = c.parent_page_type_id
		)
		SELECT merged
			FROM chain
			WHERE parent_page_type_id IS NULL`,
		pageTypeId).Scan(&substitutionTypes)

	if err != nil {
		slog.Error("Critical SQL error while getting typeset from page", "err", err)
		http.Error(w, "Could not find page in SQL", http.StatusNotFound)
		return
	}

	uid := utils_sec.GetUID(req)
	// Check if user has the permissions required to view the page
	if requiredPrivilege > 0 {
		err = utils_sec.RequireLevel(uid, requiredPrivilege)
		if err != nil {
			http.Error(w, utils_err.CodeToMessage(err.Error(), langCode), http.StatusForbidden)
			return
		}
	}

	// Get test id
	testId, siteTestId := getABver(w, req, translationId, langCode)

	// Get global substitutions
	var globalSubstitutions map[string]string
	err = db.Pool.QueryRow(context.Background(),
		`SELECT substitutions FROM sitewide_tests
			WHERE lang_code = $1
			ORDER BY (test_id = $2) DESC
			LIMIT 1`,
		langCode, siteTestId).Scan(&globalSubstitutions)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			slog.Error("Critical SQL error while getting global substitutions", "err", err)
		}
		http.Error(w, "Could not find page in SQL", http.StatusNotFound)
		return
	}

	// Get substitutions
	var substitutions map[string]string
	err = db.Pool.QueryRow(context.Background(),
		`SELECT substitutions FROM tests
			WHERE translation_id = $1
			ORDER BY (test_id = $2) DESC
			LIMIT 1`,
		translationId, testId).Scan(&substitutions)
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
			finalSubstitutions[key] = url.PathEscape("/" + langCode + pageURL)
		case "Title":
			finalSubstitutions[key] = title
		case "LangCode":
			finalSubstitutions[key] = langCode
		case "LangTags":
			finalSubstitutions[key] = page_parts.GenerateLangTags(langCode, altURLs)
		case "LangRedirects":
			finalSubstitutions[key] = page_parts.GenerateLangLinks(langCode, altURLs)
		case "Errors":
			errorURL := queryParams.Get("err")
			if errorURL != "" {
				errorMessage := utils_err.CodeToMessage(errorURL, langCode)
				finalSubstitutions[key] = html.EscapeString(errorMessage)
			} else {
				finalSubstitutions[key] = ""
			}
		case "AccountDetails":
			finalSubstitutions[key] = page_parts.GenerateAccountDetails(uid, pageURL, langCode)
		case "ReturnURL":
			returnURL := utils_url.SanitizeURL(queryParams.Get("from"))
			htmlEscaped := html.EscapeString(returnURL)
			finalSubstitutions[key] = htmlEscaped
		case "CommentSection":
			finalSubstitutions[key] = page_parts.GenerateCommentSection(uid, translationId, langCode)
		case "Comment":
			finalSubstitutions[key] = page_parts.GenerateCommentInfo(uid, langCode, queryParams)
		case "Text", "TemplateText", "Content":
			if substitutionValue, ok := substitutions[key]; ok {
				finalSubstitutions[key] = substitutionValue
			} else if globalSubstitutionValue, ok := globalSubstitutions[key]; ok {
				finalSubstitutions[key] = globalSubstitutionValue
			} else {
				slog.Warn("Untranslated content in", "url", "/"+langCode+pageURL, "key", key)
				finalSubstitutions[key] = ""
			}
		case "Creator.Dashboard":
			finalSubstitutions[key] = page_parts.GenerateCreatorDashboard(langCode)
		case "Creator.PageTypeDropdown":
			finalSubstitutions[key] = page_parts.GeneratePageTypeDropdown()
		case "Creator.Editor":
			finalSubstitutions[key] = page_parts.GenerateEditor(langCode, queryParams.Get("page"))
		case "User.AccountDetails":
			finalSubstitutions[key] = page_parts.GenerateUserPage(uid, pageURL, langCode)
		}
	}

	// Resolve template texts
	for key, val := range substitutionTypes {
		switch val {
		case "TemplateText", "Creator.Dashboard", "Creator.Editor", "CommentSection", "Comment":
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
