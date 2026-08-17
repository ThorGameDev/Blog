package creator

import (
	"blogbackend/internal/db"
	"blogbackend/internal/utils/requesturl"
	"context"
	"fmt"
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

	pageQueryParam := fmt.Sprintf("?page=%d", pageId)

	currentTranslation := req.URL.Query().Get("lang")
	redirectURL := requesturl.TranslateURL("/en/creator/editor.html", nil, currentTranslation)

	http.Redirect(w, req, redirectURL+pageQueryParam, http.StatusSeeOther)
}

func editTest(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := req.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	finalSubstitutions := make(map[string]string)
	for key, values := range req.PostForm {
		if len(values) == 1 {
			finalSubstitutions[key] = values[0]
		} else {
			slog.Warn("The wrong number of values was found with the key", "num", len(values), "key", key, "values", values)
		}
	}

	translationId := req.URL.Query().Get("translation")
	testId := req.URL.Query().Get("test")

	status, err := db.Pool.Exec(context.Background(),
		`UPDATE tests
			SET test_substitutions = $1
			WHERE translation_id = $2
			AND test_id = $3`,
		finalSubstitutions, translationId, testId)
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

func Register() {
	http.HandleFunc("/api/creator/newPage", newPage)
	http.HandleFunc("/api/creator/editTest", editTest)
}
