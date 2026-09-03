package main

import (
	"net/http"

	"blogbackend/internal/api/api_blog"
	"blogbackend/internal/api/api_creator"
	"blogbackend/internal/api/api_security"
	"blogbackend/internal/api/api_user"
	"blogbackend/internal/page"
	"blogbackend/internal/utils/db"
)

func main() {
	if err := db.Init(); err != nil {
		panic(err)
	}
	defer db.Close()

	page.Register()

	api_security.Register()
	api_creator.Register()
	api_user.Register()
	api_blog.Register()

	http.ListenAndServe(":8090", nil)
}
