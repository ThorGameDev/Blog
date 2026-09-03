package utils_sec

import (
	"blogbackend/internal/utils/db"
	"blogbackend/internal/utils/utils_err"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

// TODO: Change Session storage method
// there really is no real reason for sessions to be stored in SQL.
// A queue would be a good implementation, because that would make deleting expired sessions easy
// It just needs to be thread safe

func expireSessions() {
	slog.Error("Running expiration")
	db.Pool.Exec(context.Background(),
		`DELETE FROM sessions WHERE expire_date < $1`,
		time.Now())
}

func RegisterSession(w http.ResponseWriter, uid int) error {
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

func GetUID(req *http.Request) int {
	sessionId, err := req.Cookie("session_id")
	if err != nil {
		if err != http.ErrNoCookie {
			slog.Error("Unknown cookie error", "err", err)
		}
		return -1
	}
	var uid int
	var expireDate time.Time
	err = db.Pool.QueryRow(context.Background(),
		`SELECT uid, expire_date FROM sessions WHERE session_token = $1`,
		sessionId.Value).Scan(&uid, &expireDate)
	if err != nil {
		if err != pgx.ErrNoRows {
			slog.Error("SQL Error while getting uid", "err", err)
		}
		slog.Error("SQL error", "err", err)
		return -1
	}

	// The session exists, but has expired
	if time.Now().After(expireDate) {
		expireSessions()
		return -1
	}

	return uid
}
