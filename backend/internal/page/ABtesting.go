package page

import (
	"blogbackend/internal/utils/db"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
)

func getRandomTest(translationId int) string {
	var randTest string
	err := db.Pool.QueryRow(context.Background(),
		`SELECT test_id FROM tests
			WHERE translation_id = $1
			ORDER BY RANDOM()
			LIMIT 1`,
		translationId).Scan(&randTest)
	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Error("No active tests for translation", "translationId", translationId)
		} else {
			slog.Error("SQL error while getting random test", "translationId", translationId, "err", err)
		}
		return "00"
	}
	return randTest
}

func isTestActive(translationId int, testId string) bool {
	// exit early if no tests were decided
	if testId == "00" {
		return false
	}

	var exists bool
	err := db.Pool.QueryRow(context.Background(),
		`SELECT EXISTS (
			SELECT 1 FROM tests
				WHERE translation_id = $1
				AND test_id = $2
		)`,
		translationId, testId).Scan(&exists)
	if err != nil {
		slog.Error("Unknown error while checking if test exists!", "err", err)
		return false
	}

	return exists
}

const PAGES_PER_COOKIE = 64

func getABver(w http.ResponseWriter, req *http.Request, translationId int, langCode string) (string, string) {
	// If there is a test query parameter, display that test simply
	if req.URL.Query().Has("test") {
		return req.URL.Query().Get("test"), "00000001"
	}

	cookieId := (translationId - 1) / PAGES_PER_COOKIE            // Which segment of cookie data to use
	cookieName := fmt.Sprintf("testId-%s-%d", langCode, cookieId) // Which cookie to use
	cookieIndex := ((translationId - 1) % PAGES_PER_COOKIE) * 2   // Which character in the cookie to use

	// Extract correct cookie
	var cookieStr string
	testId := "00"
	cookieData, err := req.Cookie(cookieName)
	if err != nil {
		// Do return a random test on error
		if err != http.ErrNoCookie {
			slog.Error("Unknown cookie error", "err", err)
			return getRandomTest(translationId), "00000000"
		}
		// If there is no cookie, create an empty one
		cookieStr = strings.Repeat("0", PAGES_PER_COOKIE*2)
	} else {
		// If the found cookie is not a valid size, make a new one. Otherwise, extract data from it
		if len(cookieData.Value) < cookieIndex+2 {
			cookieStr = strings.Repeat("0", PAGES_PER_COOKIE*2)
		} else {
			cookieStr = cookieData.Value
			testId = cookieStr[cookieIndex : cookieIndex+2]
		}
	}

	// If the cookie is not set, set it
	if !isTestActive(translationId, testId) {
		testId = getRandomTest(translationId)
		cookieStr = cookieStr[:cookieIndex] + testId + cookieStr[cookieIndex+2:]

		sessionCookie := &http.Cookie{
			Name:     cookieName,
			Value:    cookieStr,
			Path:     "/" + langCode + "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, sessionCookie)
	}

	return testId, "00000000"
}
