package api_user

import (
	"blogbackend/internal/utils/utils_sec"
	"blogbackend/internal/utils/utils_url"
	"net/http"
)

func deleteAccount(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uid := utils_sec.GetUID(req)
	if uid == -1 {
		http.Error(w, "Not logged in", http.StatusMethodNotAllowed)
		return
	}

	utils_sec.DeleteAccount(w, uid)


	fromURL := utils_url.SanitizeURL(req.Header.Get("Referer"))
	http.Redirect(w, req, fromURL, http.StatusSeeOther)
}

