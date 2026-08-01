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

	// Example of getting rows
	// rows, err := db.Pool.Query(context.Background(), "SELECT username FROM users;")
	// if err != nil {
	// 	slog.Error("Error while querying the database", "err", err)
	// }
	// defer rows.Close()
	//
	// var retrieved string
	// for rows.Next() {
	// 	if err := rows.Scan(&retrieved); err != nil {
	// 		slog.Error("Critical Error!", "err", err)
	// 	}
	// 	slog.Info(retrieved)
	// }

	auth.Register()
	accountpage.Register()
	blogpage.Register()

	http.ListenAndServe(":8090", nil)
}
