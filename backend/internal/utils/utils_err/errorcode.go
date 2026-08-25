package utils_err

import (
	"blogbackend/internal/utils/db"
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

const (
	AccountExists      = "AccExs"
	UnmatchedPasswords = "UmtchP"
	IncorrectUsername  = "UWrong"
	IncorrectPassword  = "PWrong"
	InternalError      = "IntErr"
	NoAnalytics        = "NoInfo"
	SessionExpired     = "SesExp"
	NoPermissions      = "NoPerm"
)

func CodeToMessage(errID string, langCode string) string {
	var errMsg string
	err := db.Pool.QueryRow(context.Background(),
		`SELECT content FROM error_codes
			WHERE code_id = $1
			AND lang_code = $2`,
		errID, langCode).Scan(&errMsg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) && langCode != "en" {
			// If there is no translation for that error, return the untranslated error translation
			return CodeToMessage(errID, "en")
		} else if errors.Is(err, pgx.ErrNoRows) {
			// If the error doesn't exist, even in English, complain
			return "Unknown error"
		} else {
			// If there is a weird sql error, log it
			slog.Error("Problem while finding translation for error code", "err", err)
			return "Unknown error"
		}
	}
	return errMsg
}
