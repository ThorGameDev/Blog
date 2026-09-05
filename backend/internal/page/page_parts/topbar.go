package page_parts

import (
	"blogbackend/internal/utils/utils_url"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
)

func GenerateTopBar(uid int, pageURL string, langCode string, queryParams url.Values) string {
	var topbar strings.Builder
	topbar.WriteString("<div id=topBar>")
	topbar.WriteString("<h2>{{ Global.BlogTitle }}</h2>")
	topbar.WriteString("<nav>")

	// get list of all URLs for the topbar
	pageList, err := utils_url.GetPagesOfIndex(0, langCode, 2)
	if err != nil {
		slog.Error("Error while getting index", "err", err)
	}
	for _, val := range pageList {
		fmt.Fprintf(&topbar, `<a href="%s">%s</a>`, val.PageURL, val.PageTitle)
	}
	topbar.WriteString("</nav>")
	topbar.WriteString(generateAccountDetails(uid, pageURL, langCode))
	topbar.WriteString("</div>")

	// Work on side bar
	topbar.WriteString("<nav id=sideBar>")

	// get list of all URLs for the sidebar
	pageList, err = utils_url.GetPagesOfIndex(0, langCode, 1)
	if err != nil {
		slog.Error("Error while getting index", "err", err)
	}
	for _, val := range pageList {
		fmt.Fprintf(&topbar, `<a href="%s">%s</a>`, val.PageURL, val.PageTitle)
	}

	// Add the "Swap translation" links to the bottom
	topbar.WriteString(generateLangLinks(langCode, pageURL, queryParams))
	topbar.WriteString("</nav>")

	return topbar.String()
}

