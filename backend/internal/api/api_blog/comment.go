package api_blog

import (
	"blogbackend/internal/utils/db"
	"blogbackend/internal/utils/utils_sec"
	"context"
	"log/slog"
	"net/http"
	"strconv"
)

func postComment(translationId int, uid int, containerId *int, content string) {
	status, err := db.Pool.Exec(context.Background(),
		"INSERT INTO comments (translation_id, uid, container_id, content) VALUES ($1, $2, $3, $4)",
		translationId, uid, containerId, content)
	if err != nil {
		slog.Error("Failed to create comment! ", "err", err)
		return
	}
	if status.RowsAffected() != 1 {
		slog.Error("Created a weird number of rows! ", "rowsAffected", status.RowsAffected())
		return
	}
}

func comment(w http.ResponseWriter, req *http.Request) {
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

	translationId, err := strconv.Atoi(req.URL.Query().Get("translationId"))
	if err != nil {
		http.Error(w, "Error parsing query parameters", http.StatusBadRequest)
		return
	}
	currentTranslation := req.URL.Query().Get("lang")
	if currentTranslation == "" {
		currentTranslation = "en"
	}

	commentData := req.PostForm.Get("commentData")

	postComment(translationId, uid, nil, commentData)

	var fromURL string
	err = db.Pool.QueryRow(context.Background(),
		`SELECT CONCAT('/', lang_code, url) FROM translations WHERE translation_id = $1`,
		translationId).Scan(&fromURL)
	if err != nil {
		http.Error(w, "Critical SQL error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, req, fromURL, http.StatusSeeOther)
}

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
	slog.Info("Data", "c", commentId)

	currentTranslation := req.URL.Query().Get("lang")
	if currentTranslation == "" {
		currentTranslation = "en"
	}

	replyData := req.PostForm.Get("replyData")

	var fromURL string
	var translationId int
	err = db.Pool.QueryRow(context.Background(),
		`SELECT CONCAT('/', lang_code, url), translations.translation_id translation_id
			FROM comments, translations
			WHERE comment_id = $1
			AND comments.translation_id = translations.translation_id`,
		commentId).Scan(&fromURL, &translationId)
	if err != nil {
		slog.Error("", "err", err)
		http.Error(w, "Critical SQL error", http.StatusInternalServerError)
		return
	}

	postComment(translationId, uid, &commentId, replyData)

	http.Redirect(w, req, fromURL, http.StatusSeeOther)
}
