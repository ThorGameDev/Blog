package utils_url

import (
	"blogbackend/internal/utils/db"
	"context"

	"github.com/jackc/pgx/v5"
)

type PageDetails struct {
	PageURL   string
	PageTitle string
}

func GetPagesOfIndex(uid int, langCode string, neededIndexBits int) ([]PageDetails, error) {
	rows, err := db.Pool.Query(context.Background(),
		`SELECT url, title FROM translations, pages
			WHERE (page_index::int & $1) = $1
			AND pages.page_id = translations.page_id
			AND lang_code = $2`,
		neededIndexBits, langCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(
		rows,
		func(row pgx.CollectableRow) (PageDetails, error) {
			var rowURL string
			var rowTitle string
			err := row.Scan(&rowURL, &rowTitle)
			return PageDetails{PageURL: rowURL, PageTitle: rowTitle}, err
		},
	)
}
