package runner

import "time"

// Structure that holds :
//
// (1) VideoUrl - A valid URL of YouTube to download video from.
//
// (2) VideoQuality - The `VideoQuality` of provided video.
type YouTubeVideoOperations struct {
	VideoUrl     string
	VideoQuality string
}

// Function for 'YouTubeVideoOperations' that actually downloads the video.
func (youtubeVideoOperations *YouTubeVideoOperations) DownloadYouTubeVideo(downloadStatusChannel chan bool) {
	// Simulate some process time
	time.Sleep(time.Second * 5)

	// Send true status on channel to denote the downloading is completed
	downloadStatusChannel <- true

	// Close channel to indicate task is completed & this channel will not provide anything in future
	close(downloadStatusChannel)
}
