package page_parts

import (
	"blogbackend/internal/utils/db"
	"blogbackend/internal/utils/utils_url"
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

func getComments(translationId int, langCode string, containerId *int) string {
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
		if err != pgx.ErrNoRows {
			slog.Error("Failed to get comment data!", "err", err)
		}
		return ""
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
			return ""
		}

		fmt.Fprintf(&commentsSection, `<div class=comment><h3>%s</h3><img src="%s"><p>%s</p>`, username, pfpURL, content)
		if numChildren >= 1 {
			fmt.Fprintf(&commentsSection, `<a href="%s?commentId=%d" class=commentExpand>%d{{ Global.Replies }}</a>`, newURL, commentId, numChildren)
		}
		fmt.Fprintf(&commentsSection, `<a href="%s?commentId=%d" class=reply>{{ Global.SubmitReply }}</a>`, newURL, commentId)
		commentsSection.WriteString(`</div>`)
	}
	return commentsSection.String()
}

func GenerateCommentSection(uid int, translationId int, langCode string) string {
	var commentSection strings.Builder
	commentSection.WriteString(`<h2>{{ Global.CommentSectionHeader }}</h2>`)
	commentSection.WriteString(`<div id=commentSection>`)
	// Add the "Reply to" form if logged in
	if uid != -1 {
		fmt.Fprintf(&commentSection, `<form action="/api/blog/comment?translationId=%d&lang=%s" method=post>`, translationId, langCode)
		commentSection.WriteString(`<textarea name=commentData></textarea>`)
		commentSection.WriteString(`<button type=submit>{{ Global.SubmitComment }}</button>`)
		commentSection.WriteString(`</form>`)
	} else {
		commentSection.WriteString(`<p>{{ Global.LoginToComment }}</p>`)
	}
	commentSection.WriteString(`<div class=comments>`)
	commentSection.WriteString(getComments(translationId, langCode, nil))
	commentSection.WriteString(`</div></div>`)

	return commentSection.String()
}

func GenerateCommentInfo(uid int, langCode string, queryParams url.Values) string {
	commentId, err := strconv.Atoi(queryParams.Get("commentId"))
	if err != nil {
		return ""
	}

	var content string
	var username string
	var pfpURL string
	var translationId int
	err = db.Pool.QueryRow(context.Background(),
		`SELECT content, users.username, profile_pictures.url, translation_id
			FROM comments, users, profile_pictures
			WHERE comment_id = $1
			AND comments.uid = users.uid
			AND users.pfp_id = profile_pictures.pfp_id`,
		commentId).Scan(&content, &username, &pfpURL, &translationId)
	if err != nil {
		if err != pgx.ErrNoRows {
			slog.Error("SQL error while getting comment info", "err", err)
		}
		return ""
	}

	var commentInfo strings.Builder
	fmt.Fprintf(&commentInfo, `<div class=comment><h3>%s</h3><img src="%s"><p>%s</p>`, username, pfpURL, content)

	// Add the "Reply to" form if logged in
	if uid != -1 {
		fmt.Fprintf(&commentInfo, `<form action="/api/blog/reply?commentId=%d&lang=%s" method=post>`, commentId, langCode)
		commentInfo.WriteString(`<textarea name=replyData></textarea>`)
		commentInfo.WriteString(`<button type=submit>{{ Global.SubmitReply }}</button>`)
		commentInfo.WriteString(`</form>`)
	}

	// Add replies
	commentInfo.WriteString("<div class=replies>")

	replies := getComments(translationId, langCode, &commentId)
	commentInfo.WriteString(replies)
	commentInfo.WriteString("</div></div>")

	return commentInfo.String()
}
