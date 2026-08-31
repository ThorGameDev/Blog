package api_security

import (
	"blogbackend/internal/utils/db"
	"blogbackend/internal/utils/utils_err"
	"blogbackend/internal/utils/utils_sec"
	"blogbackend/internal/utils/utils_url"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func loginTo(w http.ResponseWriter, username string, password string) error {
	var pash string
	var uid int
	err := db.Pool.QueryRow(context.Background(), "SELECT password_hash, uid FROM users WHERE username = $1", username).Scan(&pash, &uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New(utils_err.IncorrectUsername)
		} else {
			slog.Error("Failed to check the existence of the account. ", "err", err)
			return errors.New(utils_err.InternalError)
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(pash), []byte(password)); err != nil {
		return errors.New(utils_err.IncorrectPassword)
	}

	// Login complete, add the login cookie
	utils_sec.RegisterSession(w, uid)

	return nil
}

func login(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := req.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	username := req.PostForm.Get("username")
	password := req.PostForm.Get("password")
	fromURL := utils_url.SanitizeURL(req.PostForm.Get("from"))
	hasJS := req.PostForm.Get("hasJS")
	currentTranslation := req.URL.Query().Get("lang")
	if currentTranslation == "" {
		currentTranslation = "en"
	}

	err = loginTo(w, username, password)
	if err != nil {
		if hasJS == "1" {
			errMsg := utils_err.CodeToMessage(err.Error(), currentTranslation)
			fmt.Fprintf(w, "%s", errMsg)
		} else {
			params := url.Values{}
			params.Add("err", err.Error())
			params.Add("from", fromURL)
			translatedURL := utils_url.TranslateURL("/en/login.html", params, currentTranslation)
			http.Redirect(w, req, translatedURL, http.StatusSeeOther)
		}
		return
	}

	http.Redirect(w, req, fromURL, http.StatusSeeOther)
}
