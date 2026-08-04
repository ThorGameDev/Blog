package main

import (
	"net/http"

	"blogbackend/internal/db"
	"blogbackend/internal/page"
	"blogbackend/internal/security/accounts/auth"
)

func main() {
	if err := db.Init(); err != nil {
		panic(err)
	}
	defer db.Close()

	auth.Register()
	page.Register()

	http.ListenAndServe(":8090", nil)
}
