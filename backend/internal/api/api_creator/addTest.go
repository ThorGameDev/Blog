package api_creator

import (
	"blogbackend/internal/utils/db"
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5"
)

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
