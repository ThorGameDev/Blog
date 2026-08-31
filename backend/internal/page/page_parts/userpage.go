package page_parts

import (
	"blogbackend/internal/utils/db"
	"context"
	"fmt"
	"log/slog"
	"strings"
)

func GenerateUserPage(uid int, pageURL string, langCode string) (string, error) {
	if uid == -1 {
		return "<h1>Not logged in!</h1>", nil
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
		slog.Error("Error while fetching user information", "err", err)
		return "", err
	}
	var pageData strings.Builder
	fmt.Fprintf(&pageData, `<h1>%s</h1>`, username)
	fmt.Fprintf(&pageData, `<img src="%s">`, pfp_url)
	return pageData.String(), nil
}
