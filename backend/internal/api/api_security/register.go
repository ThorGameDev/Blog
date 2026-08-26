package api_security

import "net/http"

func Register() {
	http.HandleFunc("/api/security/signup", signup)
	http.HandleFunc("/api/security/login", login)
}
