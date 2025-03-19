package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func CreateTempDir() (string, error) {
	return os.MkdirTemp("", "zip-extract")
}

func CleanupTempDir(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		log.Printf("failed to cleanup the temp directory %v\n", err)
	}
}

/*
SaveZipFile saves the contents of an io.Reader to a zip file in the specified temporary directory.
The function returns the full path to the saved zip file, or an empty string if an error occurs.
The temporary directory is assumed to be valid and writable, ofcourse.
*/
func SaveZipFile(tempDir, filename string, file io.Reader) string {
	zipPath := filepath.Join(tempDir, filename)
	outFile, err := os.Create(zipPath)
	if err != nil {
		return ""
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		return ""
	}

	return zipPath
}

/*
ExtractZip extracts the contents of a zip file to a temporary directory.

Parameters:
- zipPath (string): The path to the zip file to be extracted.
- tempDir (string): The path to the temporary directory where the extraction will take place.

Returns:
- string: The path to the directory where the zip file contents were extracted.
- error: Any error that occurred during the extraction process.

The function creates a new directory named "extracted" within the temporary directory, opens
the zip file, and extracts its contents to the "extracted" directory. If any errors occur during
the extraction, the function returns an empty string and the error.
*/
func ExtractZip(zipPath, tempDir string) (string, error) {
	extractDir := filepath.Join(tempDir, "extracted")
	err := os.Mkdir(extractDir, 0755)
	if err != nil {
		return "", err
	}

	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	// actual extraction of the files from the zip file
	for _, file := range reader.File {
		path := filepath.Join(extractDir, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return "", err
		}
		defer outFile.Close()

		zipFile, err := file.Open()
		if err != nil {
			return "", err
		}
		defer zipFile.Close()

		_, err = io.Copy(outFile, zipFile)
		if err != nil {
			return "", err
		}
	}

	return extractDir, nil
}

func DetectHandlerFile(dir string) (string, string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", "", err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		switch filepath.Ext(file.Name()) {
		case ".py":
			return file.Name(), "python", nil
		case ".go":
			return file.Name(), "golang", nil
		}
	}

	return "", "", fmt.Errorf("no valid handler file format found (expected '.py' or '.go')")
}
