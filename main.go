package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/", handleRoot)

	port := ":8080"
	fmt.Printf("Server listening at port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server %s", err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile("./index.html")
	if err != nil {
		slog.Error("Failed to read index file")
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(content)
	if err != nil {
		slog.Error("Failed to write response")
		return
	}
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		slog.Error("Failed to parse MultipartForm")
		http.Error(w, "Failed to parse Form", http.StatusBadRequest)
		return
	}

	// get uploaded file
	file, _, err := r.FormFile("file")
	if err != nil {
		slog.Error("Failed to get file from form")
		http.Error(w, "Failed to get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// create new file
	filename := fmt.Sprintf("uploaded_file_%d.json", time.Now().Unix())
	out, err := os.Create(filename)
	if err != nil {
		slog.Error("Failed to create new file for upload")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		slog.Error("Failed to copy uploaded file to server file")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	slog.Info(fmt.Sprintf("File uploaded successfully %s", filename))
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}
