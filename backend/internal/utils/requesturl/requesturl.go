package requesturl

import (
	"fmt"
	"net/http"
)

// Eventually, it should hardcode the site's specific URLs. For now, no urls exist
func GetRequestURL(req *http.Request) string {
	urlscheme := "http"
	if req.TLS != nil || req.Header.Get("X-Forwarded-Proto") == "https" {
		urlscheme = "https"
	}
	return fmt.Sprintf("%s://%s", urlscheme, req.Host)
}

