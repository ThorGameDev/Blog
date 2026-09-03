package page_parts

import (
	"blogbackend/internal/utils/db"
	"blogbackend/internal/utils/utils_url"
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/jackc/pgx/v5"
)

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
func GenerateAccountDetails(uid int, pageURL string, langCode string) string {
	if uid == -1 {
		return createLoginLinks(pageURL, langCode)
	}
	var username string
	var pfp_url string
	err := db.Pool.QueryRow(context.Background(),
		`SELECT username, profile_pictures.url
			FROM users, profile_pictures
			WHERE uid = $1
			AND profile_pictures.pfp_id = users.pfp_id`,
		uid).Scan(&username, &pfp_url)
	if err != nil {
		if err == pgx.ErrNoRows {
			return createLoginLinks(pageURL, langCode)
		}
		slog.Error("Error while fetching user information", "err", err)
		return ""
	}
	accountPageURL := utils_url.TranslateURL("/en/user.html", nil, langCode)
	return fmt.Sprintf(`<a href="%s"><h3>%s</h3><img src="%s"></a>`, accountPageURL, username, pfp_url)
}
