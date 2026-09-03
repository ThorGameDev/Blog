package page_parts

import (
	"blogbackend/internal/utils/db"
	"blogbackend/internal/utils/utils_url"
	"context"
	"fmt"
	"log/slog"
	"strings"
)

func GenerateLangTags(currentLangCode string, alternatePages []utils_url.LangURL) string {
	var domainURL string
	err := db.Pool.QueryRow(context.Background(),
		`SELECT val FROM site_settings WHERE key = 'URL'`).Scan(&domainURL)
	if err != nil {
		slog.Error("Could not get url from Site_settings!", "err", err)
		return ""
	}

	var langTagCluster strings.Builder
	for _, val := range alternatePages {
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

func GenerateLangLinks(currentLangCode string, alternatePages []utils_url.LangURL) string {
	var langLinks strings.Builder
	for _, val := range alternatePages {
		if currentLangCode != val.LangCode {
			toURL := "/" + val.LangCode + val.PageURL + val.QueryParams
			fmt.Fprintf(&langLinks, `<a rel="alternate" hreflang="%s" href="%s">%s</a>`, val.LangCode, toURL, val.LangName)
		}
	}
	return langLinks.String()
}
