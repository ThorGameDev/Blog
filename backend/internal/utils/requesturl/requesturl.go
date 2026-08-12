package requesturl

import (
	"blogbackend/internal/db"
	"blogbackend/internal/security/whitelist"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/jackc/pgx/v5"
)

type LangURL struct {
	LangCode    string
	PageURL     string
	QueryParams string
	LangName    string
	IsPrimary   bool
}

func GetAlternateURLs(fromPage string, fromLangCode string, queryParams url.Values) ([]LangURL, error) {
	rows, err := db.Pool.Query(context.Background(),
		`SELECT translations.lang_code, url, lang_name, is_primary FROM translations, languages
			WHERE languages.lang_code = translations.lang_code
			AND page_id = (
				SELECT page_id FROM translations 
					WHERE lang_code = $1 
					AND url = $2
			)`,
		fromLangCode, fromPage,
	)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(
		rows,
		func(row pgx.CollectableRow) (LangURL, error) {
			var rowLangCode string
			var rowURL string
			var rowLangName string
			var rowIsPrimary bool
			err := row.Scan(&rowLangCode, &rowURL, &rowLangName, &rowIsPrimary)
			newQueryParams := TranslateQueryParams(queryParams, rowLangCode)
			return LangURL{LangCode: rowLangCode, PageURL: rowURL, QueryParams: newQueryParams, LangName: rowLangName, IsPrimary: rowIsPrimary}, err
		},
	)
}

// Eventually, it should hardcode the site's specific URLs. For now, no urls exist
func GetRequestURL(req *http.Request) string {
	urlscheme := "http"
	if req.TLS != nil || req.Header.Get("X-Forwarded-Proto") == "https" {
		urlscheme = "https"
	}
	return fmt.Sprintf("%s://%s", urlscheme, req.Host)
}

func TranslateQueryParams(queryParams url.Values, newLangCode string) string {
	if len(queryParams) == 0 {
		return ""
	}
	fromQuery := queryParams.Get("from")
	translatedFrom := TranslateURL(fromQuery, nil, newLangCode)
	if translatedFrom == "" {
		return "?" + queryParams.Encode()
	}
	queryParams.Set("from", translatedFrom)
	encoded := queryParams.Encode()
	return "?" + encoded
}

func TranslateURL(decodedURL string, queryParams url.Values, newLangCode string) string {
	newLangCode = whitelist.SanitizeLangCode(newLangCode)
	var suffix string
	if len(queryParams) != 0 {
		suffix = TranslateQueryParams(queryParams, newLangCode)
	}

	var prefix string
	if len(decodedURL) >= 4 {
		if decodedURL[0] == '/' && decodedURL[3] == '/' {
			decodedURL = decodedURL[3:]
			prefix = "/" + newLangCode
		} else {
			prefix = ""
		}
	}

	var newURL string
	err := db.Pool.QueryRow(context.Background(),
		`SELECT url FROM translations
			WHERE page_id = (
				SELECT page_id FROM translations 
					WHERE url = $1
			) AND lang_code = $2`,
		decodedURL, newLangCode).Scan(&newURL)
	if err != nil {
		slog.Error("SQL Error", "err", err)
		return ""
	}
	return prefix + newURL + suffix
}
