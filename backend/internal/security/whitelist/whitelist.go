package whitelist

func SanitizeURL(rawURL string) string {
	validURLs := map[string]bool{
		"/index.html":   true,
		"/admin.html":   true,
		"/account.html": true,
	}
	if validURLs[rawURL] {
		return rawURL
	}
	return "/index.html"
}
