package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

func loginTo(w http.ResponseWriter, username string, password string) error {
	if username != "user" {
		return errors.New("nouser")
	}

	if password != "pass" {
		return errors.New("nopass")
	}

	// Most of these settings are temporary. A better cookie is needed latter, such as when https is required for accounts
	sessionCookie := &http.Cookie{
		Name:     "Session",
		Value:    "Logged In",
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, sessionCookie)

	return nil
}

func createAccount(w http.ResponseWriter, username string, password string, confirmPass string) error {
	if username == "old" {
		return errors.New("exists")
	}

	if password != confirmPass {
		return errors.New("badpass")
	}

	//If sql, create account logic here

	return loginTo(w, username, password)
}

func asValidURL(rawURL string) string {
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

func asValidSignupError(errid string) string {
	errmap := map[string]string{
		"exists":  "That account already exists!",
		"badpass": "The passwords do not match!",
		"nouser":  "Incorrect username!",
		"nopass":  "Incorrect password!",
	}
	if _, ok := errmap[errid]; ok {
		return errmap[errid]
	}
	return "Unknown error"
}

func login(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := req.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	username := req.PostForm.Get("username")
	password := req.PostForm.Get("password")
	fromurl := asValidURL(req.PostForm.Get("from"))
	hasjs := req.PostForm.Get("hasjs")

	err = loginTo(w, username, password)
	if err != nil {
		if hasjs == "1" {
			errmsg := asValidSignupError(err.Error())
			fmt.Fprintf(w, "%s", errmsg)
		} else {
			errmsg := url.QueryEscape(err.Error())
			http.Redirect(w, req, "/login.html?from="+fromurl+"&err="+errmsg, http.StatusSeeOther)
		}
		return
	}

	http.Redirect(w, req, fromurl, http.StatusSeeOther)
}

func signup(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := req.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	username := req.PostForm.Get("username")
	password := req.PostForm.Get("password")
	fromurl := asValidURL(req.PostForm.Get("from"))
	confirmPass := req.PostForm.Get("confirmPass")
	hasjs := req.PostForm.Get("hasjs")

	err = createAccount(w, username, password, confirmPass)
	if err != nil {
		if hasjs == "1" {
			errmsg := asValidSignupError(err.Error())
			fmt.Fprintf(w, "%s", errmsg)
		} else {
			errmsg := url.QueryEscape(err.Error())
			http.Redirect(w, req, "/signup.html?from="+fromurl+"&err="+errmsg, http.StatusSeeOther)
		}
		return
	}

	http.Redirect(w, req, fromurl, http.StatusSeeOther)
}

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

	fromurl := asValidURL(req.URL.Query().Get("from"))
	page = strings.ReplaceAll(page, "/index.html", fromurl)

	errorurl := req.URL.Query().Get("err")
	if errorurl != ""{
		errormsg := asValidSignupError(errorurl)
		page = strings.ReplaceAll(page, `<p id=errorDisplay>`, `<p class=errorDisplay>`+errormsg)
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, page)
}

func signupPage(w http.ResponseWriter, req *http.Request) {
	url := "http://nginx-frontend/signup.html"
	accountPage(w, req, url)
}

func loginPage(w http.ResponseWriter, req *http.Request) {
	url := "http://nginx-frontend/login.html"
	accountPage(w, req, url)
}

func main() {
	slog.Info("Starting backend!")
	http.HandleFunc("/api/security/signup", signup)
	http.HandleFunc("/api/security/login", login)

	http.HandleFunc("/login.html", loginPage)
	http.HandleFunc("/signup.html", signupPage)

	http.ListenAndServe(":8090", nil)
}
