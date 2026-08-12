package whitelist

import (
	"blogbackend/internal/db"
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

// This section can be sped up a ton by creating a cache immediately
// Implement the observer pattern, so that a "Cache expire" feature is easy to implement

func SanitizeLangCode(rawLangCode string) string {
	if len(rawLangCode) != 2 {
		return "en"
	}

	var exists bool
	err := db.Pool.QueryRow(context.Background(),
		`SELECT EXISTS(
			SELECT 1 FROM languages
			WHERE lang_code = $1
		)`,
		rawLangCode).Scan(&exists)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			slog.Error("Failed to check the validity of the langcode. ", "err", err)
		}
		return "en"
	}
	if exists {
		return rawLangCode
	}

	return "en"
}

func SanitizeURL(rawURL string) string {
	// Make sure URL is long enough to contain a langCode
	if len(rawURL) <= 4 {
		return "/index.html"
	}

	// Split url into code and location
	langCode := rawURL[1:3]
	pageURL := rawURL[3:]

	// Check if page exists in website
	var exists bool
	err := db.Pool.QueryRow(context.Background(),
		`SELECT EXISTS(
			SELECT 1 FROM translations, languages
			WHERE languages.lang_code = translations.lang_code
			AND languages.lang_code = $1
			AND url = $2
		)`,
		langCode, pageURL).Scan(&exists)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			slog.Error("Failed to check the existence of the URL. ", "err", err)
		}
		return "/index.html"
	}
	if exists {
		return rawURL
	}

	return "/index.html"
}
