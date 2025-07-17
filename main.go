package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"
)

type Results struct {
	PhotosDownloaded int
	PhotosSkipped    int
	VideosDownloaded int
	VideosSkipped    int
}

const (
	DateFormat = "2006-01-02"
	MediaDir   = "media"
)

var (
	API_BASE     = os.Getenv("API_BASE_URL")
	AuthToken    = os.Getenv("AUTH_TOKEN")
	Logger       = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	R            Results
	Timer        = time.Now()
	Version      = "development"
	MediaList    = []Asset{}
	DownloadDate string
)

func main() {
	Logger.Info("Starting...", "version", Version)

	CheckRequiredEnvVars("API_BASE_URL", "AUTH_TOKEN")

	Date := flag.String("date", "", "Date in YYYY-MM-DD format to download photos and videos for")
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	if *debug {
		Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
		Logger.Debug("Debug logging enabled")
	}

	var D time.Time
	if *Date != "" {
		D, _ = time.Parse(DateFormat, *Date)
		DownloadDate = D.Format(DateFormat)
		Logger.Info("Using provided date for download", "date", DownloadDate)
	} else {
		DownloadDate = time.Now().Format(DateFormat)
		Logger.Info("No date provided, using today's date", "date", DownloadDate)
	}

	URLPhotos := API_BASE + "/web/parent/photos/?page=1&filters%5Bphoto%5D%5Bdatetime_from%5D=" + DownloadDate + "%2000%3A00^&filters%5Bphoto%5D%5Bdatetime_to%5D=" + DownloadDate + "%2023%3A59"

	URLVideos := API_BASE + "/web/parent/videos/?page=1&filters%5Bvideo%5D%5Bdatetime_from%5D=" + DownloadDate + "%2000%3A00&filters%5Bvideo%5D%5Bdatetime_to%5D=" + DownloadDate + "%2023%3A59"

	if err := os.MkdirAll(fmt.Sprintf("%s/%s", MediaDir, DownloadDate), 0o755); err != nil {
		Logger.Error(err.Error())
	}

	// Fetch photos json
	photoResponse, err := sendRequest(URLPhotos)
	if err != nil {
		Logger.Error(fmt.Sprintf("Error sending request for photos: %v", err))
		os.Exit(1)
	}

	var photosResponseData PhotosResponse
	if err = decodeResponse(photoResponse, &photosResponseData); err != nil {
		Logger.Error(fmt.Sprintf("Error decoding photos response: %v", err))
		os.Exit(1)
	}

	for _, p := range photosResponseData.Photos {
		MediaList = append(MediaList, p)
	}

	// Fetch videos json
	videoResponse, err := sendRequest(URLVideos)
	if err != nil {
		Logger.Error(fmt.Sprintf("Error sending request for photos: %v", err))
		os.Exit(1)
	}
	var videosResponseData VideoResponse
	if err = decodeResponse(videoResponse, &videosResponseData); err != nil {
		Logger.Error(fmt.Sprintf("Error decoding photos response: %v", err))
		os.Exit(1)
	}

	for _, p := range videosResponseData.Videos {
		MediaList = append(MediaList, p)
	}

	for _, m := range MediaList {
		err = Download(m)
		if err != nil {
			var eSkipped *ErrorSkipped
			switch errors.As(err, &eSkipped) {
			case true:
				switch m.(type) {
				case Photo:
					R.PhotosSkipped++
				case Video:
					R.VideosSkipped++
				}
				Logger.Debug("File already exists, skipping download", "filename", m.GetFilename())
			default:
				Logger.Error(err.Error())
			}
		}

		switch m.(type) {
		case Photo:
			R.PhotosDownloaded++
		case Video:
			R.VideosDownloaded++
		}
		Logger.Debug("Finished processing media", "filename", m.GetFilename())
	}

	Logger.Info("Photos done", "downloaded", R.PhotosDownloaded, "skipped", R.PhotosSkipped)
	Logger.Info("Videos done", "downloaded", R.VideosDownloaded, "skipped", R.VideosSkipped)
	Logger.Info("Done", "time", time.Since(Timer).String())
}

func CheckRequiredEnvVars(v ...string) {
	for _, env := range v {
		if os.Getenv(env) == "" {
			Logger.Error(fmt.Sprintf("%s environment variable is not set", env))
			os.Exit(1)
		}
	}
}
