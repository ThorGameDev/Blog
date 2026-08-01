package blogpage

import (
	"blogbackend/internal/db"
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

func pagegen(w http.ResponseWriter, req *http.Request) {
	page := req.URL.Path[6:]

	var content string
	err := db.Pool.QueryRow(context.Background(), "SELECT content FROM translations, content WHERE url = $1 AND translations.translationid = content.translationid", page).Scan(&content)
	if err != nil {
		slog.Error("Could not find page in SQL", "err", err)
		http.Error(w, "Could not find page in SQL", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, content)
}

func Register() {
	http.HandleFunc("/blog/", pagegen)
}

