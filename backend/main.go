package main

import (
	"fmt"
	"net/http"
)

func helloWorld(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello World! From Go!\n")
}

func main() {
	fmt.Println("Starting backend!")
	http.HandleFunc("/api/helloWorld", helloWorld)

	http.ListenAndServe(":8090", nil)
}
