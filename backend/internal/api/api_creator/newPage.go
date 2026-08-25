package api_creator

import (
	"blogbackend/internal/utils/db"
	"blogbackend/internal/utils/utils_url"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5"
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
	redirectURL := utils_url.TranslateURL("/en/creator/editor.html", nil, currentTranslation)

	http.Redirect(w, req, redirectURL+pageQueryParam, http.StatusSeeOther)
}

func nextTestId(translationId string) (string, error) {
	// get previous testId from sql
	var previousTestId string
	err := db.Pool.QueryRow(context.Background(),
		`SELECT test_id FROM tests
			WHERE translation_id = $1
			ORDER BY test_id DESC`,
		translationId).Scan(&previousTestId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "00", nil
		}
		return "", err
	}

	// Get the int representation of the current testId
	var value uint64
	for _, r := range previousTestId {
		value <<= 5 // Shift left by 5 bits (base32 is 2^5)
		if r >= '0' && r <= '9' {
			value |= uint64(r - '0')
		} else if r >= 'A' && r <= 'V' {
			value |= uint64(r - 'A' + 10)
		} else if r >= 'a' && r <= 'v' {
			value |= uint64(r - 'a' + 10)
		}
	}

	// Increment the testId
	value++

	// Ensure data fits in two characters (aka 10 bits total)
	if value >= 1024 {
		return "", errors.New("Value too long")
	}

	// Map index to letters
	alphabet := "0123456789ABCDEFGHIJKLMNOPQRSTUV"
	testId := make([]byte, 2)
	testId[1] = alphabet[value&31]      // Compare the first set of 5 bits
	testId[0] = alphabet[(value>>5)&31] // Compare the second set of 5 bits

	// Return as string
	return string(testId), nil
}

func addTest(w http.ResponseWriter, req *http.Request) {
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

	translationId := req.URL.Query().Get("translation")
	testId, err := nextTestId(translationId)
	if err != nil {
		slog.Error("Error getting testId", "err", err)
		http.Error(w, "Error getting testId", http.StatusBadRequest)
		return
	}

	status, err := db.Pool.Exec(context.Background(),
		`INSERT INTO tests (test_id, translation_id, test_substitutions) VALUES
			($1, $2, '{}'::JSONB)`,
		testId, translationId)
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

func editTest(w http.ResponseWriter, req *http.Request) {
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
	http.HandleFunc("/api/creator/addTest", addTest)
	http.HandleFunc("/api/creator/addTranslation", addTranslation)
}
