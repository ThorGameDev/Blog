package api_user

import (
	"blogbackend/internal/utils/db"
	"blogbackend/internal/utils/utils_err"
	"blogbackend/internal/utils/utils_sec"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/google/uuid"
)

func copyPFPToWorkdir(pfpFileData multipart.File, filename string) error {
	outFile := "/workdir/" + filename

	out, err := os.Create(outFile)
	if err != nil {
		slog.Error("Unable to create the file for writing", "err", err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, pfpFileData)
	if err != nil {
		slog.Error("Something went wrong while writing the file", "err", err)
		return err
	}

	return nil
}

func workerProcessPFP(filename string) (string, error) {
	response, err := http.Get("http://image-sanitizer:8100/request?file=" + filename)
	if err != nil {
		slog.Error("Something went wrong with sanitizer API", "err", err)
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		slog.Error("Wrong status code", "status Code", response.StatusCode)
		return "", errors.New(utils_err.InternalError)
	}

	text, err := io.ReadAll(response.Body)
	if err != nil {
		slog.Error("Failed to get the resulting filepath", "err", err)
		return "", err
	}

	return string(text), nil
}

func moveFile(indir string, outdir string) error {
	inputFile, err := os.Open(indir)
	if err != nil {
		slog.Error("Failed to open re-encoded image file", "err", err)
		return err
	}
	defer inputFile.Close()

	if err := os.MkdirAll(filepath.Dir(outdir), 0755); err != nil {
		slog.Error("Failed to ensure directory exists", "err", err)
		return err
	}

	outputFile, err := os.Create(outdir)
	if err != nil {
		slog.Error("Unable to create the file for writing", "err", err)
		return err
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		slog.Error("Something went wrong while writing the file", "err", err)
		return err
	}

	err = os.Remove(indir)
	if err != nil {
		slog.Error("failure while removing original file", "err", err)
		return err
	}

	return nil
}

func deleteUnusedPFPs() {
	// Delete all unused profile pictures from postgreSQL
	rows, err := db.Pool.Query(context.Background(),
		`DELETE FROM profile_pictures
			WHERE user_uploaded = TRUE
			AND pfp_id NOT IN (SELECT pfp_id FROM users)
			RETURNING url`,
	)
	if err != nil {
		slog.Error("Failed to delete the unused profile pictures from SQL", "err", err)
		return
	}
	defer rows.Close()

	// Remove the associated files as well
	for rows.Next() {
		var deleteURL string
		if err := rows.Scan(&deleteURL); err != nil {
			slog.Error("Critical Error while trying to read deleted pfp's url!", "err", err)
		}
		os.Remove(deleteURL)
	}
}

func changePFP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uid := utils_sec.GetUID(req)
	if uid == -1 {
		http.Error(w, "Not logged in", http.StatusBadRequest)
		return
	}

	// Get the profile picture from the form
	pfpFileData, header, err := req.FormFile("newPFP")
	if err != nil {
		slog.Error("Error with file", "err", err)
		http.Error(w, "Error with file", http.StatusBadRequest)
		return
	}
	defer pfpFileData.Close()
	filename := uuid.New().String() + path.Ext(header.Filename)

	// Save file to the workdir
	err = copyPFPToWorkdir(pfpFileData, filename)
	if err != nil {
		return
	}

	// Inform worker of the new image to sanitize
	processedFile, err := workerProcessPFP(filename)
	if err != nil {
		return
	}

	// Copy the sanitized image to the server
	uploadFile := "/uploads/pfp/" + processedFile
	err = moveFile("/workdir/out/"+processedFile, uploadFile)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: ")
	fmt.Fprintf(w, uploadFile)

	// Link new profile picture in place of old one
	status, err := db.Pool.Exec(context.Background(),
		`WITH new_pfp AS (
			INSERT INTO profile_pictures (user_uploaded, url) VALUES
				(TRUE, $1)
			RETURNING pfp_id
		)
		UPDATE users
			SET pfp_id = new_pfp.pfp_id
			FROM new_pfp
			WHERE uid = $2`,
		uploadFile, uid)
	if err != nil {
		slog.Error("Failed to update user password! ", "err", err)
		return
	}
	if status.RowsAffected() != 1 {
		slog.Error("Created a weird number of rows! ", "rowsAffected", status.RowsAffected())
		return
	}

	// Delete the old profile picture from disk
	deleteUnusedPFPs()
}
