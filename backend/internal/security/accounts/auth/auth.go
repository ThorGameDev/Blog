package auth

import (
	"blogbackend/internal/db"
	"blogbackend/internal/page/accountpage/errorcode"
	"blogbackend/internal/security/whitelist"
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
	err := db.Pool.QueryRow(context.Background(), "SELECT pash, uid FROM users WHERE username = $1", username).Scan(&pash, &uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New(errorcode.IncorrectUsername)
		} else {
			slog.Error("Failed to check the existence of the account. ", "err", err)
			return errors.New(errorcode.InternalSqlError)
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(pash), []byte(password)); err != nil {
		return errors.New(errorcode.IncorrectPassword)
	}

	// Login good, proceed to setting session
	randbytes := make([]byte, 32)
	if _, err := rand.Read(randbytes); err != nil {
		return errors.New("")
	}
	cookie := base64.URLEncoding.EncodeToString(randbytes)

	expireDate := time.Now().Add(time.Hour * 24) // Expires cookie after a day of neglect

	status, err := db.Pool.Exec(context.Background(),
		"INSERT INTO sessions (session_token, uid, expire_date) VALUES ($1, $2, $3)",
		cookie, uid, expireDate)
	if err != nil {
		slog.Error("Failed to create account! ", "err", err)
		return errors.New(errorcode.InternalSqlError)
	}
	if status.RowsAffected() != 1 {
		slog.Error("Created a weird number of rows! ", "rowsAffected", status.RowsAffected())
		return errors.New(errorcode.InternalSqlError)
	}

	sessionCookie := &http.Cookie{
		Name:     "SessionID",
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
		return errors.New(errorcode.UnmatchedPasswords)
	}

	var exists bool
	err := db.Pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
	if err != nil {
		slog.Error("Failed to check the existence of the account. ", "err", err)
		return errors.New(errorcode.InternalSqlError)
	}
	if exists {
		return errors.New(errorcode.AccountExists)
	}

	pash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Error while hashing password.", "err", err)
		return errors.New(errorcode.InternalError)
	}

	status, err := db.Pool.Exec(context.Background(),
		"INSERT INTO users (username, pash, pfp_file_id, privilege) VALUES ($1, $2, 'file', 0)",
		username, string(pash))
	if err != nil {
		slog.Error("Failed to create account! ", "err", err)
		return errors.New(errorcode.InternalSqlError)
	}

	if status.RowsAffected() != 1 {
		slog.Error("Created a weird number of rows! ", "rowsAffected", status.RowsAffected())
		return errors.New(errorcode.InternalSqlError)
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
	fromurl := whitelist.SanitizeURL(req.PostForm.Get("from"))
	hasjs := req.PostForm.Get("hasjs")

	err = loginTo(w, username, password)
	if err != nil {
		if hasjs == "1" {
			errmsg := errorcode.CodeToMessage(err.Error())
			fmt.Fprintf(w, "%s", errmsg)
		} else {
			errmsg := url.QueryEscape(err.Error())
			http.Redirect(w, req, "/login.html?from="+fromurl+"&err="+errmsg, http.StatusSeeOther)
		}
		return
	}

	http.Redirect(w, req, fromurl, http.StatusSeeOther)
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
	fromurl := whitelist.SanitizeURL(req.PostForm.Get("from"))
	confirmPass := req.PostForm.Get("confirmPass")
	hasjs := req.PostForm.Get("hasjs")

	err = createAccount(w, username, password, confirmPass)
	if err != nil {
		if hasjs == "1" {
			errmsg := errorcode.CodeToMessage(err.Error())
			fmt.Fprintf(w, "%s", errmsg)
		} else {
			errmsg := url.QueryEscape(err.Error())
			http.Redirect(w, req, "/signup.html?from="+fromurl+"&err="+errmsg, http.StatusSeeOther)
		}
		return
	}

	http.Redirect(w, req, fromurl, http.StatusSeeOther)
}

func Register() {
	http.HandleFunc("/api/security/signup", signup)
	http.HandleFunc("/api/security/login", login)
}
