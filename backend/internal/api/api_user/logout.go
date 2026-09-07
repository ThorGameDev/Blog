package api_user

import (
	"blogbackend/internal/utils/utils_sec"
	"blogbackend/internal/utils/utils_url"
	"net/http"
)

func logout(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	utils_sec.EndSession(w, req)

	fromURL := utils_url.SanitizeURL(req.Header.Get("Referer"))
	http.Redirect(w, req, fromURL, http.StatusSeeOther)
}

