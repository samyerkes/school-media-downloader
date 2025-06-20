package main

import "fmt"

type Photo struct {
	ID      string `json:"id"`
	MainURL string `json:"main_url"`
}

type PhotosResponse struct {
	Photos []Photo `json:"photos"`
}

// GetDownloadURL returns the URL to download the photo
func (p Photo) GetDownloadURL() string {
	return p.MainURL
}

// GetFilename returns the filename where the photo will be saved
func (p Photo) GetFilename() string {
	return fmt.Sprintf("%s/%s/%s.jpg", MediaDir, DownloadDate, p.ID)
}

// GetID returns the unique identifier for the photo
func (p Photo) GetID() string {
	return p.ID
}
