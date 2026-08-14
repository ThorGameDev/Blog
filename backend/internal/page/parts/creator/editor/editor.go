package editor

import (
	"blogbackend/internal/db"
	"context"
	"fmt"
	"log/slog"
	"strings"
)

func GenerateEditor(langCode string, pageId string) string {
	var substitution_types map[string]string
	var type_name string
	err := db.Pool.QueryRow(context.Background(),
		`SELECT substitution_types, type_name
			FROM pages, page_type
			WHERE page_id = $1
			AND pages.page_type_id = page_type.page_type_id`,
		pageId).Scan(&substitution_types, &type_name)
	if err != nil {
		slog.Error("Error while getting substitution types", "err", err)
		return ""
	}

	var editorData strings.Builder
	var comparisonData strings.Builder
	// Get data on each translation in block
	translationRows, err := db.Pool.Query(context.Background(),
		`SELECT translations.translation_id,
				translations.url,
				translations.substitutions->>'PageTitle',
				translations.lang_code
			FROM pages, translations
			WHERE translations.page_id = pages.page_id
			AND pages.page_id = $1
			ORDER BY (lang_code = $2) DESC`,
		pageId, langCode)
	if err != nil {
		slog.Error("Error while querying the database", "err", err)
		return ""
	}
	defer translationRows.Close()

	for translationRows.Next() {
		var translationId int
		var rowURL string
		var rowPageTitle string
		var rowLangCode string
		if err := translationRows.Scan(&translationId, &rowURL, &rowPageTitle, &rowLangCode); err != nil {
			slog.Error("Critical Error!", "err", err)
			return ""
		}

		fmt.Fprintf(&editorData, `<h3>%s</h3>`, rowPageTitle)
		fmt.Fprintf(&editorData, `<a hreflang="%s" href="/%s%s">%s</a>`, rowLangCode, rowLangCode, rowURL, rowURL)

		testRows, err := db.Pool.Query(context.Background(),
			`SELECT test_id FROM tests
				WHERE translation_id = $1`,
			translationId)
		if err != nil {
			slog.Error("Error while querying the database", "err", err)
			return ""
		}
		defer testRows.Close()

		for testRows.Next() {
			var testId string
			if err := testRows.Scan(&testId); err != nil {
				slog.Error("Critical Error!", "err", err)
				return ""
			}

			fmt.Fprintf(&editorData, ` <a hreflang="%s" href="/%s%s?test=%s">%s</a>`, rowLangCode, rowLangCode, rowURL, testId, testId)
			fmt.Fprintf(&editorData, ` <button type="button" class="testVisibility on" id="Toggle_e-%s-%s">On</button>`, rowLangCode, testId)

			// Create page comparison display
			fmt.Fprintf(&comparisonData, `<div id="e-%s-%s" class="directEditor">`, rowLangCode, testId)
			fmt.Fprintf(&comparisonData, `<iframe src="/%s%s?test=%s" sandbox></iframe>`, rowLangCode, rowURL, testId)
			comparisonData.WriteString(`</div>`)
		}
	}
	editorData.WriteString(`<div class="sideBySide">`)
	editorData.WriteString(comparisonData.String())
	editorData.WriteString(`</div>`)

	return editorData.String()
}
