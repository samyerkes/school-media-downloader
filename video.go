package main

import "fmt"

type Video struct {
	ID       string `json:"id"`
	VideoURL string `json:"video_file_url"`
}

type VideoResponse struct {
	Videos []Video `json:"videos"`
}

// GetDownloadURL returns the URL to download the video
func (v Video) GetDownloadURL() string {
	return v.VideoURL
}

// GetFilename returns the filename where the video will be saved
func (v Video) GetFilename() string {
	return fmt.Sprintf("%s/%s/%s.mp4", MediaDir, DownloadDate, v.ID)
}

// GetID returns the unique identifier for the video
func (v Video) GetID() string {
	return v.ID
}
