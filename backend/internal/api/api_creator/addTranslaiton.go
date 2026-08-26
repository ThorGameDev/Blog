package api_creator

import (
	"blogbackend/internal/utils/db"
	"context"
	"log/slog"
	"net/http"
)

func addTranslation(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := req.ParseForm()
	if err != nil {
		slog.Error("Error parsing form", "err", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	pageId := req.URL.Query().Get("pageId")
	//langCode := req.URL.Query().Get("lang")

	translationLangCode := req.PostForm.Get("language")
	translationUrl := req.PostForm.Get("url")
	translationSubstitutions := map[string]interface{}{
		"PageTitle": req.PostForm.Get("pageTitle"),
	}

	status, err := db.Pool.Exec(context.Background(),
		`INSERT INTO translations (page_id, lang_code, substitutions, url) VALUES
		($1, $2, $3, $4)`,
		pageId, translationLangCode, translationSubstitutions, translationUrl)
	if err != nil {
		slog.Error("Failed to modify page!", "err", err)
		http.Error(w, "Failed to modify page!", http.StatusInternalServerError)
		return
	}
	if status.RowsAffected() != 1 {
		slog.Error("Created a weird number of rows!", "rowsAffected", status.RowsAffected())
		http.Error(w, "Created a weird number of rows!", http.StatusInternalServerError)
		return
	}
}
