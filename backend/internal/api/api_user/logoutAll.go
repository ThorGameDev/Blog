package api_user

import (
	"blogbackend/internal/utils/utils_sec"
	"blogbackend/internal/utils/utils_url"
	"net/http"
)

func logoutAll(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uid := utils_sec.GetUID(req)
	if uid == -1 {
		http.Error(w, "Not logged in", http.StatusBadRequest)
		return
	}

	utils_sec.EndAllSession(w, uid)

	fromURL := utils_url.SanitizeURL(req.Header.Get("Referer"))
	http.Redirect(w, req, fromURL, http.StatusSeeOther)
}

