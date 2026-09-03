package api_blog

import "net/http"

func Register() {
	http.HandleFunc("/api/blog/comment", comment)
	http.HandleFunc("/api/blog/reply", reply)
}
