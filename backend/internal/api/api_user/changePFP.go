package api_user

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

func changePFP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	pfpFileData, header, err := req.FormFile("newPFP")
	if err != nil {
		slog.Error("Error with file", "err", err)
		http.Error(w, "Error with file", http.StatusBadRequest)
		return
	}
	defer pfpFileData.Close()

	// TODO ensure filename is uuid
	filename := "/uploads/"+header.Filename

	out, err := os.Create(filename)
	if err != nil {
		slog.Error("Unable to create the file for writing", "err", err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, pfpFileData)
	if err != nil {
		slog.Error("Something went wrong while writing the file", "err", err)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: ")
	fmt.Fprintf(w, filename)
}
