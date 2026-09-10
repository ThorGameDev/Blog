package api_user

import "net/http"

func Register() {
	http.HandleFunc("/api/user/logout", logout)
	http.HandleFunc("/api/user/logoutAll", logoutAll)
	http.HandleFunc("/api/user/changePassword", changePassword)
	http.HandleFunc("/api/user/changePFP", changePFP)
	http.HandleFunc("/api/user/deleteAccount", deleteAccount)
}
