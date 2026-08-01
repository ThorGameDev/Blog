package main

import (
	"net/http"

	"blogbackend/internal/db"
	"blogbackend/internal/page/accountpage"
	"blogbackend/internal/page/blogpage"
	"blogbackend/internal/security/accounts/auth"
)

func main() {
	if err := db.Init(); err != nil {
		panic(err)
	}
	defer db.Close()

	auth.Register()
	accountpage.Register()
	blogpage.Register()

	http.ListenAndServe(":8090", nil)
}
