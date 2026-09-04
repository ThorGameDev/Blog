package api_creator

import (
	"blogbackend/internal/utils/db"
	"blogbackend/internal/utils/utils_url"
	"context"
	"log/slog"
	"net/http"
)

func newPage(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := req.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	pageType := req.PostForm.Get("pageType")
	minPermissions := req.PostForm.Get("minPermissions")

	var pageId int
	err = db.Pool.QueryRow(context.Background(),
		"INSERT INTO pages (page_type_id, required_privilege) VALUES ($1, $2) RETURNING page_id",
		pageType, minPermissions).Scan(&pageId)
	if err != nil {
		slog.Error("Failed to create new page! ", "err", err)
		http.Error(w, "Failed to create new page! ", http.StatusInternalServerError)
		return
	}

	fromURL := utils_url.SanitizeURL(req.Header.Get("Referer"))
	http.Redirect(w, req, fromURL, http.StatusSeeOther)
}
