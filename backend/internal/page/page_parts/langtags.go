package page_parts

import (
	"blogbackend/internal/utils/db"
	"blogbackend/internal/utils/utils_url"
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
)

func GenerateLangTags(currentLangCode string, pageURL string, queryParams url.Values) string {
	domainURL := db.GetSiteSetting("URL")

	altURLs, err := utils_url.GetAlternateURLs(pageURL, currentLangCode, nil)
	if err != nil {
		slog.Error("Failure getting alternate URLs!", "err", err)
		return ""
	}

	var langTagCluster strings.Builder
	for _, val := range altURLs {
		toURL := domainURL + "/" + val.LangCode + val.PageURL
		fmt.Fprintf(&langTagCluster, `<link rel="alternate" hreflang="%s" href="%s">`, val.LangCode, toURL)
		if val.IsPrimary {
			fmt.Fprintf(&langTagCluster, `<link rel="alternate" hreflang="x-default" href="%s">`, toURL)
		}
	}

	var langTags string
	err = db.Pool.QueryRow(context.Background(),
		`SELECT page_tags FROM languages WHERE lang_code = $1`,
		currentLangCode).Scan(&langTags)
	if err != nil {
		slog.Error("Failed to get language specific tags from SQL", "err", err)
		return ""
	}

	langTagCluster.WriteString(langTags)

	return langTagCluster.String()
}

func generateLangLinks(currentLangCode string, pageURL string, queryParams url.Values) string {
	altURLs, err := utils_url.GetAlternateURLs(pageURL, currentLangCode, queryParams)
	if err != nil {
		slog.Error("Failure getting alternate URLs!", "err", err)
		return ""
	}

	var langLinks strings.Builder
	for _, val := range altURLs {
		if currentLangCode != val.LangCode {
			toURL := "/" + val.LangCode + val.PageURL + val.QueryParams
			fmt.Fprintf(&langLinks, `<a rel="alternate" hreflang="%s" href="%s">%s</a>`, val.LangCode, toURL, val.LangName)
		}
	}
	return langLinks.String()
}
