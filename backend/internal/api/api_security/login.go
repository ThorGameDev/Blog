package api_security

import (
	"blogbackend/internal/utils/db"
	"blogbackend/internal/utils/utils_err"
	"blogbackend/internal/utils/utils_url"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

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

	// Login good, proceed to setting session
	randBytes := make([]byte, 32)
	if _, err := rand.Read(randBytes); err != nil {
		return errors.New("")
	}
	cookie := base64.URLEncoding.EncodeToString(randBytes)

	expireDate := time.Now().Add(time.Hour * 24) // Expires cookie after a day of neglect

	status, err := db.Pool.Exec(context.Background(),
		"INSERT INTO sessions (session_token, uid, expire_date) VALUES ($1, $2, $3)",
		cookie, uid, expireDate)
	if err != nil {
		slog.Error("Failed to create account! ", "err", err)
		return errors.New(utils_err.InternalError)
	}
	if status.RowsAffected() != 1 {
		slog.Error("Created a weird number of rows! ", "rowsAffected", status.RowsAffected())
		return errors.New(utils_err.InternalError)
	}

	sessionCookie := &http.Cookie{
		Name:     "session_id",
		Value:    cookie,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, sessionCookie)

	return nil
}

func createAccount(w http.ResponseWriter, username string, password string, confirmPass string) error {
	if password != confirmPass {
		return errors.New(utils_err.UnmatchedPasswords)
	}

	var exists bool
	err := db.Pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
	if err != nil {
		slog.Error("Failed to check the existence of the account. ", "err", err)
		return errors.New(utils_err.InternalError)
	}
	if exists {
		return errors.New(utils_err.AccountExists)
	}

	pash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Error while hashing password.", "err", err)
		return errors.New(utils_err.InternalError)
	}

	status, err := db.Pool.Exec(context.Background(),
		"INSERT INTO users (username, password_hash, pfp_file_id, privilege) VALUES ($1, $2, 'file', 1)",
		username, string(pash))
	if err != nil {
		slog.Error("Failed to create account! ", "err", err)
		return errors.New(utils_err.InternalError)
	}

	if status.RowsAffected() != 1 {
		slog.Error("Created a weird number of rows! ", "rowsAffected", status.RowsAffected())
		return errors.New(utils_err.InternalError)
	}

	return loginTo(w, username, password)
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

func signup(w http.ResponseWriter, req *http.Request) {
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
	confirmPass := req.PostForm.Get("confirmPass")
	hasJS := req.PostForm.Get("hasJS")
	currentTranslation := req.URL.Query().Get("lang")
	if currentTranslation == "" {
		currentTranslation = "en"
	}

	err = createAccount(w, username, password, confirmPass)
	if err != nil {
		if hasJS == "1" {
			errMsg := utils_err.CodeToMessage(err.Error(), currentTranslation)
			fmt.Fprintf(w, "%s", errMsg)
		} else {
			params := url.Values{}
			params.Add("err", err.Error())
			params.Add("from", fromURL)
			translatedURL := utils_url.TranslateURL("/en/signup.html", params, currentTranslation)
			http.Redirect(w, req, translatedURL, http.StatusSeeOther)
		}
		return
	}

	http.Redirect(w, req, fromURL, http.StatusSeeOther)
}

func Register() {
	http.HandleFunc("/api/security/signup", signup)
	http.HandleFunc("/api/security/login", login)
}
