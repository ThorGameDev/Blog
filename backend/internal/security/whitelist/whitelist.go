package whitelist

import (
	"blogbackend/internal/db"
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

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
