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
	queryParams := url.Values{}
	queryParams.Set("from", "/"+langCode+fromPageUrl)

	loginURL := utils_url.TranslateURL("/en/login.html", queryParams, langCode)

	// TODO: Use a "Not logged in" profile picture, and add a "Not logged in" tool tip
	return fmt.Sprintf(`<a id=accountIcon href="%s"><img src=/res/default_pfp.png></a>`, loginURL)
}

func generateAccountDetails(uid int, pageURL string, langCode string) string {
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
	return fmt.Sprintf(`<a id=accountIcon href="%s" title="%s"><img src="%s"></a>`, accountPageURL, username, pfp_url)
}
