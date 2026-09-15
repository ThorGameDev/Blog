package api_blog

import (
	"blogbackend/internal/page/page_parts"
	"blogbackend/internal/utils/db"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
)

func getReplies(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	queryParams := req.URL.Query()

	langCode := queryParams.Get("lang")
	commentId, err := strconv.Atoi(queryParams.Get("commentId"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var translationId int
	err = db.Pool.QueryRow(context.Background(),
		`SELECT translation_id FROM comments WHERE comment_id = $1`,
		commentId).Scan(&translationId)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		slog.Error("Failed to get translationId attached to comment", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	commentHtml := page_parts.GetComments(translationId, langCode, &commentId)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, commentHtml)
}
