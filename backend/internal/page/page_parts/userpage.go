package page_parts

import (
	"blogbackend/internal/utils/db"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
)

func GenerateUserPage(req *http.Request, pageURL string, langCode string) (string, error) {
	session_id, err := req.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			return "<h1>Not logged in!</h1>", nil
		}
		slog.Error("Failed to get session ID cookie!", "err", err)
		return "", err
	}
	var uid int
	var username string
	var pfp_url string
	err = db.Pool.QueryRow(context.Background(),
		`SELECT users.uid, username, profile_pictures.url
			FROM users, sessions, profile_pictures
			WHERE users.uid = sessions.uid
			AND sessions.session_token = $1
			AND profile_pictures.pfp_id = users.pfp_id`,
		session_id.Value).Scan(&uid, &username, &pfp_url)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "<h1>Session Expired!</h1>", nil
		}
		slog.Error("Error while fetching user information", "err", err)
		return "", err
	}
	var pageData strings.Builder
	fmt.Fprintf(&pageData, `<h1>%s</h1>`, username)
	fmt.Fprintf(&pageData, `<img src="%s">`, pfp_url)
	return pageData.String(), nil
}
