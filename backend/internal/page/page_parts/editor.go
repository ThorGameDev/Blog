package page_parts

import (
	"blogbackend/internal/utils/db"
	"context"
	"fmt"
	"html"
	"log/slog"
	"sort"
	"strings"
)

var esc = html.EscapeString

func GenerateEditor(langCode string, pageId string) string {
	var editorData strings.Builder

	missingLanguages, err := db.Pool.Query(context.Background(),
		`SELECT lang_code, lang_name FROM languages
			WHERE lang_code NOT IN (
				SELECT lang_code FROM translations
					WHERE page_id = $1
			)`, pageId)
	if err != nil {
		slog.Error("Error while querying the database", "err", err)
		return ""
	}
	defer missingLanguages.Close()

	var missingLangOptions strings.Builder
	missingLanguage := false

	for missingLanguages.Next() {
		missingLanguage = true
		var missingLangCode string
		var missingLangName string
		if err := missingLanguages.Scan(&missingLangCode, &missingLangName); err != nil {
			slog.Error("Critical Error!", "err", err)
			return ""
		}
		fmt.Fprintf(&missingLangOptions, `<option value="%s">%s</option>`, missingLangCode, missingLangName)
	}

	if missingLanguage {
		fmt.Fprintf(&editorData, `<form id=addTranslationForm action="/api/creator/addTranslation?lang=%s&pageId=%s" method=post>`, langCode, pageId)
		fmt.Fprintf(&editorData, `<select name=language>%s</select>`, missingLangOptions.String())
		editorData.WriteString(`<input name=url type=text placeholder={{ NewTranslationUrlPrompt }}>`)
		editorData.WriteString(`<input name=pageTitle type=text placeholder={{ NewTranslationTitlePrompt }}>`)
		editorData.WriteString(`<button type=submit>{{ AddTranslationPrompt }}</button>`)
		editorData.WriteString(`</form>`)
	}

	var substitutionTypes map[string]string
	var typeName string
	err = db.Pool.QueryRow(context.Background(),
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

	// Sort the substitution types, so that iterating is done consistently. For now, it's just alphabetical order
	subKeys := make([]string, 0, len(substitutionTypes))
	for key := range substitutionTypes {
		subKeys = append(subKeys, key)
	}
	sort.Strings(subKeys)

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
		fmt.Fprintf(&editorData, `<form action="/api/creator/addTest?translation=%d" method="post">`, translationId)
		editorData.WriteString(`<button type="submit">{{ AddTestPrompt }}</button></form></div>`)
		editorData.WriteString(`</form>`)

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

			// TODO: secure the iframe
			// Sandbox is probably not needed, on account of the fact that this page is all site controlled
			// But chances are, removing it will allow for an xss privilege escalation.
			// Its one thing to hack the users, but It's another to have access to the page generator.
			// But I'm removing the sandbox to allow for font loading. Currently, there are CORS errors
			fmt.Fprintf(&comparisonData, `<div class="iframeHold"><iframe src="/%s%s?test=%s"></iframe></div>`, rowLangCode, rowURL, testId)

			fmt.Fprintf(&comparisonData, `<form action="/api/creator/editTest?translation=%d&test=%s" method="post">`, translationId, testId)
			// Create the editor
			for _, key := range subKeys {
				subType := substitutionTypes[key]
				escKey := esc(key)
				if subType == "Text" || subType == "TemplateText" || subType == "Content" {
					inputFieldId := fmt.Sprintf("i-%s-%s", escKey, versionCode)
					fmt.Fprintf(&comparisonData, `<label for="%s">%s</label>`, inputFieldId, escKey)
					translationVal, ok := rowTranslationSubstitutions[key]
					testVal := esc(testSubstitutions[key])
					var additionalAttributes string
					if ok {
						// Maybe it's silly, but in keeping with the desire to minify everything, the separating space should be here instead of the other string
						additionalAttributes = " disabled"
						testVal = esc(translationVal)
					}

					switch subType {
					case "Text":
						fmt.Fprintf(&comparisonData, `<input type="text" name="%s" id="%s" value="%s"%s>`, escKey, inputFieldId, testVal, additionalAttributes)
					case "TemplateText":
						// Will be different eventually
						fmt.Fprintf(&comparisonData, `<input type="text" name="%s" id="%s" value="%s"%s>`, escKey, inputFieldId, testVal, additionalAttributes)
					case "Content":
						fmt.Fprintf(&comparisonData, `<textarea name="%s" id="%s"%s>%s</textarea>`, escKey, inputFieldId, additionalAttributes, esc(testSubstitutions[key]))
					}
					comparisonData.WriteString(`</br>`)
				}
			}
			// Close the editor form and div, while including a submit button
			comparisonData.WriteString(`<button type="submit">{{ SaveChangesPrompt }}</button></form></div>`)
		}
	}
	editorData.WriteString(`<div class="sideBySide">`)
	editorData.WriteString(comparisonData.String())
	editorData.WriteString(`</div>`)

	return editorData.String()
}
