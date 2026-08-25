package main

import (
	"net/http"

	"blogbackend/internal/creator"
	"blogbackend/internal/db"
	"blogbackend/internal/page"
	"blogbackend/internal/security/accounts/auth"
	userapi "blogbackend/internal/user"
)

func main() {
	if err := db.Init(); err != nil {
		panic(err)
	}
	defer db.Close()

	auth.Register()
	page.Register()
	creator.Register()
	userapi.Register()

	http.ListenAndServe(":8090", nil)
}
