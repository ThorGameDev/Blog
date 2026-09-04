package api_blog

import (
	"blogbackend/internal/utils/db"
	"blogbackend/internal/utils/utils_sec"
	"blogbackend/internal/utils/utils_url"
	"context"
	"log/slog"
	"net/http"
	"strconv"
)

func reply(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := req.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	uid := utils_sec.GetUID(req)
	if uid == -1 {
		http.Error(w, "Not logged in", http.StatusBadRequest)
		return
	}

	commentId, err := strconv.Atoi(req.URL.Query().Get("commentId"))
	if err != nil {
		http.Error(w, "Error parsing query parameters", http.StatusBadRequest)
		return
	}

	currentTranslation := req.URL.Query().Get("lang")
	if currentTranslation == "" {
		currentTranslation = "en"
	}

	replyData := req.PostForm.Get("replyData")

	var translationId int
	err = db.Pool.QueryRow(context.Background(),
		`SELECT translations.translation_id translation_id
			FROM comments, translations
			WHERE comment_id = $1
			AND comments.translation_id = translations.translation_id`,
		commentId).Scan(&translationId)
	if err != nil {
		slog.Error("", "err", err)
		http.Error(w, "Critical SQL error", http.StatusInternalServerError)
		return
	}

	postComment(translationId, uid, &commentId, replyData)

	fromURL := utils_url.SanitizeURL(req.Header.Get("Referer"))
	http.Redirect(w, req, fromURL, http.StatusSeeOther)
}
