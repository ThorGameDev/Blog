package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func helloWorld(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello World! From Go!\n")
}

func createAccount(w *http.ResponseWriter, username string, password string, confirmPass string) error{
	if username == "old" {
		return errors.New("That account already exists!")
	}

	if password != confirmPass {
		return errors.New("The passwords do not match!")
	}

	// Most of these settings are temporary. A better cookie is needed latter, such as when https is required for accounts
	sessionCookie := &http.Cookie {
		Name: "Session",
		Value: "Logged In",
		Path: "/",
		MaxAge: 3600,
		HttpOnly: true,
		//Secure: true,
		Secure: false,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(*w, sessionCookie)

	return nil
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
	from := req.PostForm.Get("from")
	if from == "" {
		from =  "/index.html"
	}
	fromurl := url.QueryEscape(from)
	confirmPass := req.PostForm.Get("confirmPass")
	hasjs := req.PostForm.Get("hasjs")
	
	err = createAccount(&w, username, password, confirmPass)
	if err != nil {
		if hasjs == "1" {
			fmt.Fprintf(w, err.Error())
		} else {
			errmsg := url.QueryEscape(err.Error())
			http.Redirect(w, req, "/signup.html?from=" + fromurl + "&err=" + errmsg, http.StatusSeeOther)
		}
		return
	}
	
	// Success! Return to previous page
	if hasjs == "1" {
		fmt.Fprintf(w, "Success!")
	} else {
		http.Redirect(w, req, fromurl, http.StatusSeeOther)
	}
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

func signupPage(w http.ResponseWriter, req *http.Request) {
	url := "http://nginx-frontend/signup.html"
	page, err := retrievePage(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error fetching signup page: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fromurl := req.URL.Query().Get("from")
	if fromurl == "" {
		fromurl =  "/index.html"
	}
	errormsg := req.URL.Query().Get("err")

	replaced := strings.ReplaceAll(page, "/index.html", fromurl)
    errorDisp := strings.ReplaceAll(replaced, `<p id=errorDisplay>`, `<p class=errorDisplay>` + errormsg)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, errorDisp)
}

func loginPage(w http.ResponseWriter, req *http.Request) {
	url := "http://nginx-frontend/login.html"
	page, err := retrievePage(url)
	if err != nil {
		log.Printf("Error fetching login page: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fromurl := req.URL.Query().Get("from")
	if fromurl == "" {
		fromurl =  "/index.html"
	}
	errormsg := req.URL.Query().Get("err")

	replaced := strings.ReplaceAll(page, "/index.html", fromurl)
    errorDisp := strings.ReplaceAll(replaced, `<p id=errorDisplay>`, `<p class=errorDisplay>` + errormsg)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, errorDisp)
}

func main() {
	fmt.Println("Starting backend!")
	http.HandleFunc("/api/helloWorld", helloWorld)
	http.HandleFunc("/api/security/signup", signup)

	http.HandleFunc("/login.html", loginPage)
	http.HandleFunc("/signup.html", signupPage)

	http.ListenAndServe(":8090", nil)
}
