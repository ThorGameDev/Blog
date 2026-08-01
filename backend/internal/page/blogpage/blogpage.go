package blogpage

import (
	"blogbackend/internal/db"
	"blogbackend/internal/page/retrieve"
	"blogbackend/internal/utils/requesturl"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"strings"
	"text/template"
)

func pagegen(w http.ResponseWriter, req *http.Request) {
	domainURL := requesturl.GetRequestURL(req)

	fullpageurl := req.URL.Path
	pagepath, pageurl := path.Split(fullpageurl)

	// Get language of page
	var langcode string
	var langtags string
	err := db.Pool.QueryRow(context.Background(),
		`SELECT langcode, pagetags FROM languages WHERE urlbase = $1`,
		pagepath).Scan(&langcode, &langtags)
	if err != nil {
		slog.Error("Could not find corresponding URL base", "err", err)
		http.Error(w, "Could not find corresponding URL base", http.StatusNotFound)
		return
	}

	// Get specific page
	var content string
	var title string
	var pageid int
	err = db.Pool.QueryRow(context.Background(),
		`SELECT content, title, pageid
			FROM translations, content 
			WHERE langcode = $1
			AND url = $2 
			AND translations.translationid = content.translationid`,
		langcode, pageurl).Scan(&content, &title, &pageid)
	if err != nil {
		slog.Error("Could not find page in SQL", "err", err)
		http.Error(w, "Could not find page in SQL", http.StatusNotFound)
		return
	}

	// Generate alternate languages refs
	rows, err := db.Pool.Query(context.Background(),
		`SELECT languages.langcode, CONCAT(urlbase, url), langname FROM translations, languages 
			WHERE languages.langcode = translations.langcode
			AND pageid = $1;`,
		pageid)
	if err != nil {
		slog.Error("Error while querying the database", "err", err)
	}
	defer rows.Close()

	var additionalTags strings.Builder
	var LangRedirects strings.Builder

	var rowlangcode string
	var rowurl string
	var rowlangname string
	for rows.Next() {
		if err := rows.Scan(&rowlangcode, &rowurl, &rowlangname); err != nil {
			slog.Error("Critical Error!", "err", err)
		}
		fmt.Fprintf(&additionalTags, `<link rel="alternate" hreflang="%s" href="%s%s">`, rowlangcode, domainURL, rowurl)
		if langcode != rowlangcode {
			fmt.Fprintf(&LangRedirects, `<a rel="alternate" hreflang="%s" href="%s">%s</a>`, rowlangcode, rowurl, rowlangname)
		}
	}

	// Append the language specific tags to the alternate_tags
	additionalTags.WriteString(langtags)

	// Get base page
	url := "http://nginx-frontend:8080/templates/blogpage.html"
	pageData, err := retrieve.RetrievePage(url)
	if err != nil {
		slog.Error("Error fetching page", "url", url, "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Parse base page as html template
	pageTemplate, err := template.New(pageurl).Parse(pageData)
	if err != nil {
		slog.Error("Could not parse HTML template", pageData, url, "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Insert data
	pageSubstitutions := struct {
		PageURL        string
		LangCode       string
		PageTitle      string
		AdditionalTags string
		Content        string
		LangRedirects  string
	}{
		PageURL:        fullpageurl,
		LangCode:       langcode,
		PageTitle:      title,
		AdditionalTags: additionalTags.String(),
		Content:        content,
		LangRedirects:  LangRedirects.String(),
	}

	w.Header().Set("Content-Type", "text/html")
	pageTemplate.Execute(w, pageSubstitutions)
}

// Implement a bad URL redirector
// Users might just manually rewrite the url, and it shouldn't break then
// /ja/blog/page1.html should redirect to /ja/ブログ/パージ１.html
// /en/ブログ/パージ１.html should redirect to /en/blog/page1.html
// Use http.Redirect(w, req, "to", http.StatusPermanentRedirect)

func Register() {
	// Get a list of URL entrypoints
	rows, err := db.Pool.Query(context.Background(), "SELECT urlbase FROM languages;")
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
		http.HandleFunc(retrieved, pagegen)
	}
}
