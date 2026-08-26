package api_creator

import "net/http"

func Register() {
	http.HandleFunc("/api/creator/newPage", newPage)
	http.HandleFunc("/api/creator/editTest", editTest)
	http.HandleFunc("/api/creator/addTest", addTest)
	http.HandleFunc("/api/creator/addTranslation", addTranslation)
}
