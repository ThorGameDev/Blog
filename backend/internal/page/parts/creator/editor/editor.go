package editor

import (
	"blogbackend/internal/db"
	"context"
	"fmt"
	"html"
	"log/slog"
	"sort"
	"strings"
)

var esc = html.EscapeString

func GenerateEditor(langCode string, pageId string) string {
	var substitutionTypes map[string]string
	var typeName string
	err := db.Pool.QueryRow(context.Background(),
		`SELECT substitution_types, type_name
			FROM pages, page_type
			WHERE page_id = $1
			AND pages.page_type_id = page_type.page_type_id`,
		pageId).Scan(&substitutionTypes, &typeName)
	if err != nil {
		slog.Error("Error while getting substitution types", "err", err)
		return ""
	}
	typeName = esc(typeName)

	// Sort the keys, so that iterating is done consistently. For now, it's just alphabetical order
	subKeys := make([]string, 0, len(substitutionTypes))
	for key := range substitutionTypes {
		subKeys = append(subKeys, key)
	}
	sort.Strings(subKeys)

	var editorData strings.Builder
	var comparisonData strings.Builder
	// Get data on each translation in block
	translationRows, err := db.Pool.Query(context.Background(),
		`SELECT translations.translation_id,
				translations.url,
				translations.substitutions,
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
		var rowTranslationSubstitutions map[string]string
		var rowLangCode string
		if err := translationRows.Scan(&translationId, &rowURL, &rowTranslationSubstitutions, &rowLangCode); err != nil {
			slog.Error("Critical Error!", "err", err)
			return ""
		}
		rowURL = esc(rowURL)
		rowLangCode = esc(rowLangCode)

		fmt.Fprintf(&editorData, `<h3>%s</h3>`, rowTranslationSubstitutions["PageTitle"])
		fmt.Fprintf(&editorData, `<a hreflang="%s" href="/%s%s">%s</a>`, rowLangCode, rowLangCode, rowURL, rowURL)

		testRows, err := db.Pool.Query(context.Background(),
			`SELECT test_id, test_substitutions FROM tests
				WHERE translation_id = $1`,
			translationId)
		if err != nil {
			slog.Error("Error while querying the database", "err", err)
			return ""
		}
		defer testRows.Close()

		for testRows.Next() {
			var testId string
			var testSubstitutions map[string]string
			if err := testRows.Scan(&testId, &testSubstitutions); err != nil {
				slog.Error("Critical Error!", "err", err)
				return ""
			}
			testId = esc(testId)

			versionCode := esc(fmt.Sprintf("%s-%s", rowLangCode, testId))

			fmt.Fprintf(&editorData, ` <a hreflang="%s" href="/%s%s?test=%s">%s</a>`, rowLangCode, rowLangCode, rowURL, testId, testId)
			fmt.Fprintf(&editorData, ` <button type="button" class="testVisibility on" id="Toggle_e-%s">On</button>`, versionCode)

			// Create page comparison display
			fmt.Fprintf(&comparisonData, `<div id="e-%s" class="directEditor">`, versionCode)
			fmt.Fprintf(&comparisonData, `<div class="iframeHold"><iframe src="/%s%s?test=%s" sandbox></iframe></div>`, rowLangCode, rowURL, testId)

			fmt.Fprintf(&comparisonData, `<form action="/api/creator/editTest?translation=%d&test=%s" method="post">`, translationId, testId)
			// Create the editor
			for _, key := range subKeys {
				escKey := esc(key)
				if substitutionTypes[key] == "Text" || substitutionTypes[key] == "TemplateText" {
					inputFieldId := fmt.Sprintf("i-%s-%s", escKey, versionCode)
					fmt.Fprintf(&comparisonData, `<label for="%s">%s</label>`, inputFieldId, escKey)
					translationVal, ok := rowTranslationSubstitutions[key]
					if ok {
						fmt.Fprintf(&comparisonData, `<input type="text" name="%s" id="%s" value="%s" disabled>`, escKey, inputFieldId, esc(translationVal))
					} else {
						fmt.Fprintf(&comparisonData, `<input type="text" name="%s" id="%s" value="%s">`, escKey, inputFieldId, esc(testSubstitutions[key]))
					}
					comparisonData.WriteString(`</br>`)
				}
			}
			comparisonData.WriteString(`</form></div>`)
		}
	}
	editorData.WriteString(`<div class="sideBySide">`)
	editorData.WriteString(comparisonData.String())
	editorData.WriteString(`</div>`)

	return editorData.String()
}
