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

func AsValidSignupError(errid string) string {
	errmap := map[string]string{
		"exists":   "That account already exists!",
		"badpass":  "The passwords do not match!",
		"nouser":   "Incorrect username!",
		"nopass":   "Incorrect password!",
		"sqlfail":  "Internal Server Error: Running SQL query failed!",
		"hashfail": "Internal Server Error: Failure hashing password",
	}
	if _, ok := errmap[errid]; ok {
		return errmap[errid]
	}
	return "Unknown error"
}
