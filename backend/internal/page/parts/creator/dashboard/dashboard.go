package dashboard

import (
	"blogbackend/internal/db"
	"context"
	"fmt"
	"log/slog"
	"strings"
)

func createAnalytics(pageid int) string {
	return `<div class="analytics"><p>No analytics avaliable yet</p></div>`
}

func GenerateCreatorDashboard() (string, error) {
	rows, err := db.Pool.Query(context.Background(),
		`SELECT translations.page_id,
				page_type.type_name,
				translations.url,
				translations.substitutions->>'PageTitle',
				translations.lang_code 
			FROM pages, translations, page_type 
			WHERE translations.page_id = pages.page_id 
			AND pages.page_type_id = page_type.page_type_id
			ORDER BY translations.page_id ASC`)
	if err != nil {
		slog.Error("Error while querying the database", "err", err)
	}
	defer rows.Close()

	var dashboardData strings.Builder

	// Create a main block that will sort the contents in a grid
	dashboardData.WriteString(`<div class="mainBlock">`)

	current_page_id := 0
	for rows.Next() {
		var page_id int
		var type_name string
		var url string
		var page_title string
		var lang_code string
		if err := rows.Scan(&page_id, &type_name, &url, &page_title, &lang_code); err != nil {
			slog.Error("Critical Error!", "err", err)
		}

		// If this is a new page, create the header and footer
		if page_id != current_page_id {
			// End a previous section
			if current_page_id != 0 {
				// Close translations
				dashboardData.WriteString(`</div>`)

				// Add analytic block
				dashboardData.WriteString(createAnalytics(current_page_id))

				// Close info block, add Manage Page link, and close pageBlock
				fmt.Fprintf(&dashboardData, `</div><a href="/en/creator/edit.html?page=%d">Manage Page</a></div>`, current_page_id)
			}

			// begin the current section
			fmt.Fprintf(&dashboardData, `<div class="pageBlock %s"><div class="info"><div class="translations">`, type_name)
			current_page_id = page_id
		}

		// Insert the translation data into the current section
		fmt.Fprintf(&dashboardData, `<h3>%s</h3>`, page_title)
		fmt.Fprintf(&dashboardData, `<a hreflang="%s" href="%s">%s</a>`, lang_code, url, url)
	}

	// Finish up the final pageBlock

	// Close translations
	dashboardData.WriteString(`</div>`)

	// Add analytic block
	dashboardData.WriteString(createAnalytics(current_page_id))

	// Close info block, add Manage Page link, and close pageBlock
	fmt.Fprintf(&dashboardData, `</div><a href="/en/creator/edit.html?page=%d">Manage Page</a></div>`, current_page_id)


	// Close the mainBlock
	dashboardData.WriteString("</div>")

	return dashboardData.String(), nil
}
