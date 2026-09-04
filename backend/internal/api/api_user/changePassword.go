package api_user

import (
	"blogbackend/internal/utils/db"
	"blogbackend/internal/utils/utils_sec"
	"blogbackend/internal/utils/utils_url"
	"context"
	"log/slog"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func changePassword(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := req.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	//currentTranslation := req.URL.Query().Get("lang")

	newPass := req.PostForm.Get("newPass")
	confirmPass := req.PostForm.Get("confirmPass")
	if newPass != confirmPass {
		http.Error(w, "Passwords do not match!", http.StatusBadRequest)
		return
	}

	uid := utils_sec.GetUID(req)
	if uid == -1 {
		http.Error(w, "Not logged in", http.StatusBadRequest)
		return
	}

	var pash string
	err = db.Pool.QueryRow(context.Background(),
		`SELECT users.password_hash FROM users
			WHERE uid = $1`,
		uid).Scan(&pash)
	if err != nil {
		slog.Error("Failed to decode session ID cookie!", "err", err)
		http.Error(w, "Internal Server Error", http.StatusBadRequest)
		return
	}

	oldPass := req.PostForm.Get("oldPass")
	if err := bcrypt.CompareHashAndPassword([]byte(pash), []byte(oldPass)); err != nil {
		http.Error(w, "Incorrect Password", http.StatusBadRequest)
		return
	}

	newPash, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Error while hashing password.", "err", err)
		return
	}

	status, err := db.Pool.Exec(context.Background(),
		`UPDATE users
			SET password_hash = $1
			WHERE uid = $2`,
		string(newPash), uid)
	if err != nil {
		slog.Error("Failed to update user password! ", "err", err)
		return
	}
	if status.RowsAffected() != 1 {
		slog.Error("Created a weird number of rows! ", "rowsAffected", status.RowsAffected())
		return
	}

	fromURL := utils_url.SanitizeURL(req.Header.Get("Referer"))
	http.Redirect(w, req, fromURL, http.StatusSeeOther)
}
