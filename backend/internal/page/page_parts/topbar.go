package page_parts

import (
	"blogbackend/internal/utils/utils_url"
	"fmt"
	"net/url"
	"strings"
)

func GenerateTopBar(uid int, pageURL string, langCode string, queryParams url.Values) string {
	var topbar strings.Builder
	topbar.WriteString("<div id=topBar>")
	topbar.WriteString("<h2>{{ Global.BlogTitle }}</h2>")
	topbar.WriteString("<nav>")

	// Just a few placeholder links, until I figure out what actually goes in a top bar
	homelink := utils_url.TranslateURL("/en/blog/page1.html", nil, langCode)
	fmt.Fprintf(&topbar, `<a href="%s">Page1</a>`, homelink)
	creatorDashboard := utils_url.TranslateURL("/en/creator/dashboard.html", nil, langCode)
	fmt.Fprintf(&topbar, `<a href="%s">Creator Dashboard</a>`, creatorDashboard)

	topbar.WriteString("</nav>")
	topbar.WriteString(generateAccountDetails(uid, pageURL, langCode))
	topbar.WriteString("</div>")

	// Work on side bar
	topbar.WriteString("<nav id=sideBar>")
	topbar.WriteString(generateLangLinks(langCode, pageURL, queryParams))
	topbar.WriteString("</nav>")

	return topbar.String()
}

