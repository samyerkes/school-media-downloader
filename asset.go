package main

import (
	"errors"
	"io"
	"net/http"
	"os"
)

type Asset interface {
	GetDownloadURL() string
	GetFilename() string
	GetID() string
}

type ErrorSkipped struct{}

func (e *ErrorSkipped) Error() string {
	return "file already exists, skipping download"
}

// Download an asset (photo or video) and saves it to the filesystem
func Download(a Asset) error {
	err := CheckIfFileExists(a)
	if err != nil {
		return err
	}

	out, err := CreateFile(a)
	if err != nil {
		return err
	}
	defer out.Close()

	data, err := ReadFile(a)
	if err != nil {
		return err
	}
	defer data.Close()

	err = SaveFile(out, data)
	if err != nil {
		return err
	}

	return nil
}

// CheckIfFileExists checks if the file for the asset already exists
func CheckIfFileExists(a Asset) error {
	filename := a.GetFilename()
	if _, err := os.Stat(filename); err == nil {
		return &ErrorSkipped{} // File exists, return error to skip download
	}
	return nil
}

// CreateFile creates a new file for the asset and returns the file handle
func CreateFile(a Asset) (*os.File, error) {
	filename := a.GetFilename()
	out, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReadFile fetches the asset from its download URL and returns an io.ReadCloser
func ReadFile(a Asset) (io.ReadCloser, error) {
	resp, err := http.Get(a.GetDownloadURL())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Failed to download file: " + resp.Status)
	}
	return resp.Body, nil
}

// SaveFile writes the data from the ReadCloser to the provided file
func SaveFile(out *os.File, data io.ReadCloser) error {
	_, err := io.Copy(out, data)
	if err != nil {
		return err
	}
	return nil
}
