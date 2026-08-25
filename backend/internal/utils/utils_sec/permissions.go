package utils_sec

import (
	"blogbackend/internal/utils/db"
	"blogbackend/internal/utils/utils_err"
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

const (
	ReadOnly    = 0
	User        = 1
	EarlyAccess = 2
	Bonus       = 3
	VIP         = 4
	Trusted     = 5
	Owner       = 6
)

func PermissionLevel(sessionID string) (int, error) {
	var privilege int
	err := db.Pool.QueryRow(context.Background(), "SELECT privilege FROM users, sessions WHERE session_token = $1 AND sessions.uid = users.uid", sessionID).Scan(&privilege)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errors.New(utils_err.SessionExpired)
		} else {
			slog.Error("Failed to get privilege from user. ", "err", err)
			return 0, errors.New(utils_err.InternalError)
		}
	}

	return privilege, nil
}

func RequireLevel(sessionID string, minimum int) error {
	level, err := PermissionLevel(sessionID)
	if err != nil {
		return err
	}
	if level >= minimum {
		return nil
	}
	return errors.New(utils_err.NoPermissions)
}
