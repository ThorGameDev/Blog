package page_parts

import (
	"blogbackend/internal/utils/db"
	"blogbackend/internal/utils/utils_url"
	"fmt"
	"log/slog"
	"strings"
)

func GenerateBottomBar(langCode string) string {
	var bottomBar strings.Builder

	bottomBar.WriteString("<div id=bottomBar>")

	copyright := db.GetSiteSetting("Copyright")
	fmt.Fprintf(&bottomBar, `<p>%s</p>`, copyright)

	bottomBar.WriteString("<nav>")

	licenseURL := db.GetSiteSetting("LicenseURL")
	sourceCodeURL := db.GetSiteSetting("SourceCodeURL")

	fmt.Fprintf(&bottomBar, `<a href="%s">{{ Global.AGPLv3License }}</a>`, licenseURL)
	fmt.Fprintf(&bottomBar, `<a href="%s">{{ Global.SourceCode }}</a>`, sourceCodeURL)

	// get list of all URLs for the topbar
	pageList, err := utils_url.GetPagesOfIndex(langCode, 4)
	if err != nil {
		slog.Error("Error while getting index", "err", err)
	}
	for _, val := range pageList {
		fmt.Fprintf(&bottomBar, `<a href="/%s%s">%s</a>`, langCode, val.PageURL, val.PageTitle)
	}

	bottomBar.WriteString("</nav></div>")

	return bottomBar.String()
}
