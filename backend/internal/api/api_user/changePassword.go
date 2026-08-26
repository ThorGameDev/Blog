package api_user

import (
	"blogbackend/internal/utils/db"
	"context"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5"
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

	session_id, err := req.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "No Session", http.StatusBadRequest)
			return
		}
		slog.Error("Failed to get session ID cookie!", "err", err)
		return
	}

	var pash string
	var uid int
	err = db.Pool.QueryRow(context.Background(),
		`SELECT users.password_hash, users.uid FROM users, sessions
			WHERE session_token = $1
			AND sessions.uid = users.uid`,
		session_id.Value).Scan(&pash, &uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Session Expired", http.StatusBadRequest)
			return
		}
		slog.Error("Failed to decode session ID cookie!", "err", err)
		http.Error(w, "Bad Session!", http.StatusBadRequest)
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
}
