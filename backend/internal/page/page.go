package page

import (
	"blogbackend/internal/db"
	"blogbackend/internal/page/retrieve"
	"blogbackend/internal/utils/requesturl"
	"context"
	"fmt"
	"text/template"
	"log/slog"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
)

type langURL struct {
	langCode string
	url      string
}

func getAlternateURLs(fromPage langURL) ([]langURL, error) {
	rows, err := db.Pool.Query(context.Background(),
		`SELECT lang_code, url FROM translations
			WHERE page_id = (
				SELECT page_id FROM translations 
					WHERE lang_code = $1 AND url = $2
			)`,
		fromPage.langCode, fromPage.url,
	)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(
		rows,
		func(row pgx.CollectableRow) (langURL, error) {
			var rowLangCode string
			var rowURL string
			err := row.Scan(&rowLangCode, &rowURL)
			return langURL{rowLangCode, rowURL}, err
		},
	)
}

func nameLangCode(langCode string) (string, error) {
	// This needs to be in sql, which is kinda awful considering it never changes
	switch langCode {
	case "en":
		return "english", nil
	case "ja":
		return "日本語", nil
	}
	return "None", nil
}

func getLangTags(domainURL string, currentLangCode string, alternatePages []langURL) (string, error) {
	var langTagCluster strings.Builder
	for _, val := range alternatePages {
		fmt.Fprintf(&langTagCluster, `<link rel="alternate" hreflang="%s" href="%s/%s%s">`, currentLangCode, domainURL, val.langCode, val.url)
	}

	var langTags string
	err := db.Pool.QueryRow(context.Background(),
		`SELECT page_tags FROM languages WHERE lang_code = $1`,
		currentLangCode).Scan(&langTags)

	langTagCluster.WriteString(langTags)

	return langTagCluster.String(), err
}

func getLangLinks(currentLangCode string, alternatePages []langURL) (string, error) {
	var langLinks strings.Builder
	for _, val := range alternatePages {
		if currentLangCode != val.langCode {
			langName, err := nameLangCode(val.langCode)
			if err != nil {
				return "", err
			}
			fmt.Fprintf(&langLinks, `<a rel="alternate" hreflang="%s" href="/%s%s">%s</a>`, val.langCode, val.langCode, val.url, langName)
		}
	}
	return langLinks.String(), nil
}

func pageGen(w http.ResponseWriter, req *http.Request) {
	domainURL := requesturl.GetRequestURL(req)
	slog.Info(domainURL)

	fullPageURL := req.URL.Path
	slog.Info(fullPageURL)
	langCode := fullPageURL[1:3]
	slog.Info(langCode)
	pageURL := fullPageURL[3:]
	slog.Info(pageURL)


	// Get page typeset
	var substitutionTypes map[string]string
	var templateUrl string
	err := db.Pool.QueryRow(context.Background(),
		`SELECT substitution_types, template_url
			FROM page_type, pages, translations
			WHERE pages.page_type_id = page_type.page_type_id
			AND translations.page_id = pages.page_id 
			AND lang_code = $1
			AND url = $2`,
		langCode, pageURL).Scan(&substitutionTypes, &templateUrl)
	if err != nil {
		slog.Error("Could not find page in SQL", "err", err)
		http.Error(w, "Could not find page in SQL", http.StatusNotFound)
		return
	}

	// Get substitutions
	var baseSubstitutions map[string]string
	var testSubstitutions map[string]string
	err = db.Pool.QueryRow(context.Background(),
		`SELECT substitutions, test_substitutions
			FROM translations, tests
			WHERE translations.translation_id = tests.translation_id
			AND lang_code = $1
			AND url = $2`,
		langCode, pageURL).Scan(&baseSubstitutions, &testSubstitutions)
	if err != nil {
		slog.Error("Could not find page in SQL", "err", err)
		http.Error(w, "Could not find page in SQL", http.StatusNotFound)
		return
	}

	altURLs, err := getAlternateURLs(langURL{langCode, pageURL})

	finalSubstitutions := make(map[string]string)
	for key, val := range substitutionTypes {
		switch val {
		case "URL":
			finalSubstitutions[key] = fullPageURL
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
		case "Text":
			if substitutionValue, ok := testSubstitutions[key]; ok {
				finalSubstitutions[key] = substitutionValue
			} else if substitutionValue, ok := baseSubstitutions[key]; ok {
				finalSubstitutions[key] = substitutionValue
			} else {
				slog.Warn("Untranslated content in", "url", fullPageURL, "key", key)
				finalSubstitutions[key] = ""
			}
		}
	}
	
	// Get base page
	pageData, err := retrieve.RetrievePage(templateUrl)
	if err != nil {
		slog.Error("Error fetching page", "url", templateUrl, "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Parse base page as html template
	pageTemplate, err := template.New(pageURL).Parse(pageData)
	if err != nil {
		slog.Error("Could not parse HTML template", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Insert data
	w.Header().Set("Content-Type", "text/html")
	pageTemplate.Execute(w, finalSubstitutions)
}

// Implement a bad URL redirect
// Users might just manually rewrite the url, and it shouldn't break then
// /ja/blog/page1.html should redirect to /ja/ブログ/パージ１.html
// /en/ブログ/パージ１.html should redirect to /en/blog/page1.html
// Use http.Redirect(w, req, "to", http.StatusPermanentRedirect)

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
		http.HandleFunc("/" + retrieved + "/", pageGen)
	}
}
