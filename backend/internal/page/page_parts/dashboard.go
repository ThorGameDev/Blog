package page_parts

import (
	"blogbackend/internal/utils/db"
	"blogbackend/internal/utils/utils_err"
	"blogbackend/internal/utils/utils_url"
	"context"
	"fmt"
	"log/slog"
	"strings"
)

func createAnalytics(pageid int, langCode string) string {
	return fmt.Sprintf(`<div class="analytics"><p>%s</p></div>`, utils_err.CodeToMessage(utils_err.NoAnalytics, langCode))
}

func GeneratePageTypeDropdown() string {
	var dropdown strings.Builder
	dropdown.WriteString(`<select name="pageType" id="pageType">`)

	rows, err := db.Pool.Query(context.Background(),
		`SELECT page_type_id, type_name from page_type`)
	if err != nil {
		dropdown.WriteString(`</select>`)
		return dropdown.String()
	}
	defer rows.Close()

	for rows.Next() {
		var row_type_id int
		var row_type_name string
		if err := rows.Scan(&row_type_id, &row_type_name); err != nil {
			slog.Error("Critical Error!", "err", err)
			dropdown.WriteString(`</select>`)
			return dropdown.String()
		}
		fmt.Fprintf(&dropdown, `<option value="%d">%s</option>`, row_type_id, row_type_name)
	}

	dropdown.WriteString(`</select>`)
	return dropdown.String()
}

func GenerateCreatorDashboard(langCode string) (string, error) {
	editorLink := utils_url.TranslateURL("/en/creator/editor.html", nil, langCode)

	pageRows, err := db.Pool.Query(context.Background(),
		`SELECT page_id, type_name
			FROM pages, page_type
			WHERE pages.page_type_id = page_type.page_type_id
			ORDER BY pages.page_id DESC`)
	if err != nil {
		slog.Error("Error while querying the database", "err", err)
		return "", err
	}
	defer pageRows.Close()

	var dashboardData strings.Builder

	// Create a main block that will sort the contents in a grid
	dashboardData.WriteString(`<div class="mainBlock">`)

	for pageRows.Next() {
		var pageId int
		var typeName string
		if err := pageRows.Scan(&pageId, &typeName); err != nil {
			slog.Error("Critical Error!A", "err", err)
			return "", err
		}

		// Start the block
		fmt.Fprintf(&dashboardData, `<div class="pageBlock %s"><div class="info"><div class="translations">`, typeName)

		// Get data on each translation in block
		translationRows, err := db.Pool.Query(context.Background(),
			`SELECT translations.url, title, translations.lang_code
				FROM pages, translations
				WHERE translations.page_id = pages.page_id
				AND pages.page_id = $1
				ORDER BY (lang_code = $2) DESC`,
			pageId, langCode)
		if err != nil {
			slog.Error("Error while querying the database", "err", err)
			return "", err
		}
		defer translationRows.Close()

		// Add translation data
		for translationRows.Next() {
			var rowURL string
			var rowPageTitle string
			var rowLangCode string
			if err := translationRows.Scan(&rowURL, &rowPageTitle, &rowLangCode); err != nil {
				slog.Error("Critical Error!", "err", err)
				return "", err
			}
			fmt.Fprintf(&dashboardData, `<h3>%s</h3>`, rowPageTitle)
			fmt.Fprintf(&dashboardData, `<a hreflang="%s" href="/%s%s">%s</a>`, rowLangCode, rowLangCode, rowURL, rowURL)
		}

		// Close translations block
		dashboardData.WriteString(`</div>`)

		// Add analytic block
		dashboardData.WriteString(createAnalytics(pageId, langCode))

		// Close info block, add Manage Page link, and close pageBlock
		fmt.Fprintf(&dashboardData, `</div><a href="%s?page=%d">{{ ManagePage }}</a></div>`, editorLink, pageId)

	}

	// Close Main Block
	dashboardData.WriteString("</div>")

	return dashboardData.String(), nil
}
