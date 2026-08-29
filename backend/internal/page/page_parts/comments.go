package page_parts

import (
	"blogbackend/internal/utils/db"
	"blogbackend/internal/utils/utils_url"
	"context"
	"fmt"
	"log/slog"
	"strings"
)

func getComments(translationId int, langCode string, containerId *int) (string, error) {
	// Get all root comments for the current page, with oldest first
	rows, err := db.Pool.Query(context.Background(),
		`SELECT comment_id,
				content,
				users.username,
				profile_pictures.url,
				(
					SELECT COUNT(*) FROM comments sub
					WHERE sub.container_id = comments.comment_id
				) AS num_children
			FROM comments, users, profile_pictures
			WHERE translation_id = $1
			AND comments.container_id IS NOT DISTINCT FROM $2
			AND comments.uid = users.uid
			AND users.pfp_id = profile_pictures.pfp_id
			ORDER BY comment_id ASC`,
		translationId, containerId)
	if err != nil {
		slog.Error("Failed to get comment data!", "err", err)
		return "", err
	}
	defer rows.Close()

	newURL := utils_url.TranslateURL("/en/comment.html", nil, langCode)
	var commentsSection strings.Builder
	for rows.Next() {
		var commentId int
		var content string
		var username string
		var pfpURL string
		var numChildren int
		if err := rows.Scan(&commentId, &content, &username, &pfpURL, &numChildren); err != nil {
			slog.Error("Failed to read comment", "err", err)
			return "", err
		}

		fmt.Fprintf(&commentsSection, `<div class=comment><h3>%s</h3><img src="%s"><p>%s</p>`, username, pfpURL, content)
		if numChildren >= 1 {
			fmt.Fprintf(&commentsSection, `<a href="%s?c=%d" class=commentExpand>%d Replies</a>`, newURL, commentId, numChildren)
		}
		commentsSection.WriteString(`</div>`)
	}
	return commentsSection.String(), nil
}

func GenerateCommentSection(translationId int, langCode string) (string, error) {
	return getComments(translationId, langCode, nil)
}

func GenerateCommentInfo(langCode string, commentId int) (string, error) {
	var content string
	var username string
	var pfpURL string
	var translationId int
	err := db.Pool.QueryRow(context.Background(),
		`SELECT content, users.username, profile_pictures.url, translation_id
			FROM comments, users, profile_pictures
			WHERE comment_id = $1
			AND comments.uid = users.uid
			AND users.pfp_id = profile_pictures.pfp_id`,
		commentId).Scan(&content, &username, &pfpURL, &translationId)
	if err != nil {
		return "", err
	}

	var commentInfo strings.Builder
	fmt.Fprintf(&commentInfo, `<div class=comment><h3>%s</h3><img src="%s"><p>%s</p>`, username, pfpURL, content)
	commentInfo.WriteString("<div class=replies>")

	replies, err := getComments(translationId, langCode, &commentId)
	if err != nil {
		return "", err
	}
	commentInfo.WriteString(replies)
	commentInfo.WriteString("</div></div>")

	return commentInfo.String(), nil
}
