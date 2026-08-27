package api_user

import "net/http"

func Register() {
	http.HandleFunc("/api/user/changePassword", changePassword)
	http.HandleFunc("/api/user/changePFP", changePFP)
}
