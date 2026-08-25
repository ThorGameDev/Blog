package utils_fs

import (
	"fmt"
	"io"
	"net/http"
)

func RetrievePage(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Failed to find base sign up page: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to find base sign up page: %d", res.StatusCode)
	}

	text, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("Got the page just fine, but could not be read: %w", err)
	}

	return string(text), nil
}
