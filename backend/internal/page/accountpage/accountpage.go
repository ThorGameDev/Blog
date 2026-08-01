package accountpage

import (
	"fmt"
	"html"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"blogbackend/internal/page/accountpage/errorcode"
	"blogbackend/internal/security/whitelist"
)

func retrievePage(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Failed to find base sign up page: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to find base sign up page: %d", res.StatusCode)
	}

	text, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("Got the page just fine, but could not be read: %w", err)
	}

	return string(text), nil
}

func accountPage(w http.ResponseWriter, req *http.Request, url string) {
	page, err := retrievePage(url)
	if err != nil {
		slog.Error("Error fetching page", page, url, "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fromurl := whitelist.SanitizeURL(req.URL.Query().Get("from"))
	page = strings.ReplaceAll(page, "/index.html", html.EscapeString(fromurl))

	errorurl := req.URL.Query().Get("err")
	if errorurl != "" {
		errormsg := errorcode.CodeToMessage(errorurl)
		page = strings.ReplaceAll(page, `<p id=errorDisplay>`, `<p class=errorDisplay>`+html.EscapeString(errormsg))
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, page)
}

func signupPage(w http.ResponseWriter, req *http.Request) {
	url := "http://nginx-frontend:8080/signup.html"
	accountPage(w, req, url)
}

func loginPage(w http.ResponseWriter, req *http.Request) {
	url := "http://nginx-frontend:8080/login.html"
	accountPage(w, req, url)
}

func Register() {
	http.HandleFunc("/login.html", loginPage)
	http.HandleFunc("/signup.html", signupPage)
}
