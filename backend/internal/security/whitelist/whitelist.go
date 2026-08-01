package whitelist

import (
	"blogbackend/internal/db"
	"context"
	"log/slog"
	"path"
)

func SanitizeURL(rawURL string) string {
	validURLs := map[string]bool{
		"/index.html":   true,
		"/admin.html":   true,
		"/account.html": true,
	}
	if validURLs[rawURL] {
		return rawURL
	}
	pagepath, pageurl := path.Split(rawURL)

	var exists bool
	err := db.Pool.QueryRow(context.Background(),
		`SELECT EXISTS(
			SELECT 1 FROM translations, languages
			WHERE languages.langcode = translations.langcode
			AND urlbase = $1
			AND url = $2
		)`,
		pagepath, pageurl).Scan(&exists)
	if err != nil {
		slog.Error("Failed to check the existence of the URL. ", "err", err)
		return "/index.html"
	}
	if exists {
		return rawURL
	}

	return "/index.html"
}
