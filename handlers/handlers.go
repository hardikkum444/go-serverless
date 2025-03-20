package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/hardikkum444/go-serverless/docker"
	"github.com/hardikkum444/go-serverless/utils"
)

var imageStore = make(map[string]string)

// LoggingMiddleware is a middleware that logs the request method and path
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next(w, r)
	}
}

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // putting 10 MB as the maximum file size
	if err != nil {
		http.Error(w, "failed to parse multipart form", http.StatusBadRequest)
		return
	}

	// getting the zip file from the request
	file, header, err := r.FormFile("code")
	if err != nil {
		http.Error(w, "failed to retrieve zip file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	tempDir, err := utils.CreateTempDir()
	if err != nil {
		http.Error(w, "failed to create temp directory", http.StatusInternalServerError)
		return
	}
	defer utils.CleanupTempDir(tempDir)

	zipPath := utils.SaveZipFile(tempDir, header.Filename, file)
	if zipPath == "" {
		http.Error(w, "failed to save zip file", http.StatusInternalServerError)
	}

	extractDir, err := utils.ExtractZip(zipPath, tempDir)
	if err != nil {
		http.Error(w, "failed to extract zip file", http.StatusInternalServerError)
		return
	}

	handlerFile, language, err := utils.DetectHandlerFile(extractDir)
	if err != nil {
		http.Error(w, "failed to detect handler file", http.StatusBadRequest)
		return
	}

	imageID, err := docker.BuildDockerImage(extractDir, language, handlerFile)
	if err != nil {
		http.Error(w, "failed to build docker image"+err.Error(), http.StatusInternalServerError)
		return
	}

	functionID := uuid.New().String()
	imageStore[functionID] = imageID

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Docker image built successfully. Image ID: %s\n Function ID %s\n", imageID, functionID)))
}

func ExecuteHandler(w http.ResponseWriter, r *http.Request) {
	functionID := r.URL.Query().Get("functionID")
	if functionID == "" {
		http.Error(w, "Function ID (fn) is required", http.StatusBadRequest)
		return
	}

	imageID, ok := imageStore[functionID]
	if !ok {
		http.Error(w, "function ID not found", http.StatusNotFound)
		return
	}

	output, err := docker.RunDockerContainer(imageID)
	if err != nil {
		http.Error(w, "failed to run docker container", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Output: %s\n", output)))
}
