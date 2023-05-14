package runner

import (
	"strings"
	"time"
)

// Structure that holds :
//
// VideoUrl - A valid URL of YouTube to download video from.
//
// SelectedVideoQuality - User's video quality preference.
//
// SelectedAudioQuality - User's audio quality preference.
//
// VideoMetaData - The map of metadata information about video returned by YouTube's internal API.
//
// VideoQualities - Collection of all available video qualities for provided video.
//
// AudioQualities - Collection of all available audio qualities for provided video.
type YouTubeVideoDetailsStruct struct {
	VideoUrl             string
	SelectedVideoQuality string
	SelectedAudioQuality string
	VideoMetaData        map[string]interface{}
	VideoQualities       []string
	AudioQualities       []string
}

// Function to validate & return a boolean status of whether the passed URL is right.
func (youtubeVideoDetails *YouTubeVideoDetailsStruct) IsValidYouTubeURL() bool {
	// Perform basic check if YouTube is missing from video URL, if so, then it's invalid URL!
	if !strings.Contains(
		strings.ToLower(youtubeVideoDetails.VideoUrl),
		"youtube.com",
	) && !strings.Contains(
		strings.ToLower(youtubeVideoDetails.VideoUrl),
		"youtu.be",
	) {
		return false
	}

	// TODO: Send POST request to YouTube's internal API

	// TODO: If 'status' under 'playabilityStatus' is not set to "OK" then video is not available

	// TODO: Otherwise convert response & store map of metadata in `VideoMetaData`
	// youtubeVideoDetails.VideoMetaData =

	// TODO: Extract & store video/audio information highest first, lowest last
	// youtubeVideoDetails.VideoQualities = []string{}
	// youtubeVideoDetails.AudioQualities = []string{}

	// Since everything is good, return true
	return true
}

// Function for 'YouTubeVideoDetailsStruct' that actually downloads the video.
func (youtubeVideoDetails *YouTubeVideoDetailsStruct) DownloadYouTubeVideo(downloadStateChannel chan bool) {
	// Simulate some process time
	time.Sleep(time.Second * 5)

	// Send true status on channel to denote the downloading is completed
	downloadStateChannel <- true

	// Close channel to indicate task is completed & this channel will not provide anything in future
	close(downloadStateChannel)
}
