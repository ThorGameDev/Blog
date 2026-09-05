package db

import (
	"context"
	"log/slog"
)

func GetSiteSetting(key string) string {
	var val string
	err := Pool.QueryRow(context.Background(),
		`SELECT val FROM site_settings WHERE key = $1`, key).Scan(&val)
	if err != nil {
		slog.Error("Could not get url from Site_settings!", "err", err)
		return ""
	}
	return val
}
