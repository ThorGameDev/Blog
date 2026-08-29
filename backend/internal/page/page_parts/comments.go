package page_parts

import (
	"blogbackend/internal/utils/db"
	"context"
	"fmt"
	"log/slog"
	"strings"
)

func GenerateComments(translationId int, langCode string) (string, error) {
	// Get all root comments for the current page, with oldest first
	rows, err := db.Pool.Query(context.Background(),
		`SELECT comment_id, content, users.username, profile_pictures.url
			FROM comments, users, profile_pictures
			WHERE translation_id = $1
			AND comments.container_id IS NULL
			AND comments.uid = users.uid
			AND users.pfp_id = profile_pictures.pfp_id
			ORDER BY comment_id ASC`,
		translationId)
	if err != nil {
		slog.Error("Failed to get comment data!", "err", err)
		return "", err
	}
	defer rows.Close()

	var commentsSection strings.Builder
	for rows.Next() {
		var commentId int
		var content string
		var username string
		var pfpURL string
		if err := rows.Scan(&commentId, &content, &username, &pfpURL); err != nil {
			slog.Error("Failed to read comment", "err", err)
			return "", err
		}

		fmt.Fprintf(&commentsSection, `<h3>%s</h3><img src="%s"><p>%s</p>`, username, pfpURL, content)
	}
	return commentsSection.String(), nil
}
