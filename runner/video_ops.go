package runner

import (
	"fmt"
	"strconv"
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
	VideoUrlOrID         string
	SelectedVideoQuality string
	SelectedAudioQuality string
	VideoMetaData        map[string]interface{}
	VideoQualitiesMap    map[string]string
	AudioQualitiesMap    map[string]string
}

// Function to validate & return a boolean status of whether the passed URL is right.
func (youtubeVideoDetails *YouTubeVideoDetailsStruct) IsValidYouTubeURL(
	videoValidityMsgChannel chan string,
	isValidVideo *bool,
) {
	// Make sure the `videoValidityMsgChannel` will be closed
	defer close(videoValidityMsgChannel)

	// Perform basic check if YouTube itself is missing from video URL,
	// or total characters != 11 (to consider as video ID);
	// if yes, then it's invalid URL!
	if !strings.Contains(
		strings.ToLower(youtubeVideoDetails.VideoUrlOrID),
		"youtube.com",
	) && !strings.Contains(
		strings.ToLower(youtubeVideoDetails.VideoUrlOrID),
		"youtu.be",
	) && len(youtubeVideoDetails.VideoUrlOrID) != 11 {
		videoValidityMsgChannel <- GetInvalidURLMessage(youtubeVideoDetails.VideoUrlOrID)
		*isValidVideo = false
		return
	}

	// Send POST request to YouTube's internal API to get & store map of metadata in `VideoMetaData`
	youtubeVideoDetails.VideoMetaData = GetVideoMetadataFromYouTubei(youtubeVideoDetails.VideoUrlOrID)

	// If 'status' under 'playabilityStatus' is not set to "OK" then video is not available
	responseBodyMap := youtubeVideoDetails.VideoMetaData
	playabilityStatusMap := responseBodyMap["playabilityStatus"]

	// Check if video status is not set to 'OK', return false
	if playabilityStatusMap.(map[string]interface{})["status"].(string) != "OK" {
		videoValidityMsgChannel <- GetInvalidURLMessage(youtubeVideoDetails.VideoUrlOrID)
		*isValidVideo = false
		return
	}

	// If 'signatureCipher' exist in metadata, then the video is Age-Restricted.
	streamingDataMap := responseBodyMap["streamingData"]
	// Iterate over 'adaptiveFormats' key & get the required data.
	for _, rawAdaptiveFormat := range streamingDataMap.(map[string]interface{})["adaptiveFormats"].([]interface{}) {
		currentAdaptiveFormat := rawAdaptiveFormat.(map[string]interface{})

		if _, isFound := currentAdaptiveFormat["signatureCipher"]; isFound {
			// Show warning with video title saying it can't be downloaded.
			videoDetailsMap := youtubeVideoDetails.VideoMetaData["videoDetails"]
			videoTitle := videoDetailsMap.(map[string]interface{})["title"].(string)
			videoValidityMsgChannel <- ("WARNING : '" + videoTitle + "' is marked as Age-Restricted video & so can't be downloaded...")
			*isValidVideo = false
			return
		}
	}

	// Extract & store video/audio information highest first, lowest last order
	youtubeVideoDetails.extractAdaptiveAudioVideoQualities()

	videoValidityMsgChannel <- ""
	*isValidVideo = true
}

// Extract & store video/audio information highest first, lowest last order.
func (youtubeVideoDetails *YouTubeVideoDetailsStruct) extractAdaptiveAudioVideoQualities() {
	responseBodyMap := youtubeVideoDetails.VideoMetaData
	streamingDataMap := responseBodyMap["streamingData"]

	// Allocate respective maps
	youtubeVideoDetails.AudioQualitiesMap = make(map[string]string)
	youtubeVideoDetails.VideoQualitiesMap = make(map[string]string)

	// Iterate over 'adaptiveFormats' key & get the required data.
	for _, rawAdaptiveFormat := range streamingDataMap.(map[string]interface{})["adaptiveFormats"].([]interface{}) {
		currentAdaptiveFormat := rawAdaptiveFormat.(map[string]interface{})

		formatType := currentAdaptiveFormat["mimeType"].(string)

		if strings.HasPrefix(formatType, "audio") {
			// Process the discovered audio format.
			// Audio Bitrate, Size & Type : 131073 | 5862305 (Converted to MB) | MP4/WEBM

			// Convert contentLength from bytes to MB
			audioSize, _ := strconv.ParseFloat(currentAdaptiveFormat["contentLength"].(string), 64)
			audioSize /= (1024 * 1024)

			audioFileType := ""
			if strings.Contains(formatType, "mp4") {
				audioFileType = "MP4"
			} else if strings.Contains(formatType, "webm") {
				audioFileType = "WEBM"
			}

			youtubeVideoDetails.AudioQualitiesMap[fmt.Sprintf("%d", int(currentAdaptiveFormat["bitrate"].(float64)))+" bits/s | "+fmt.Sprintf("%f", audioSize)+" MB | "+audioFileType] = currentAdaptiveFormat["url"].(string)

		} else if strings.HasPrefix(formatType, "video") {
			// Process the discovered video format.
			// Video Quality, FPS, Size & Type : 1080p60 HDR | 60 fps | 5862305 (Converted to MB) | MP4/WEBM

			// Convert contentLength from bytes to MB
			videoSize, _ := strconv.ParseFloat(currentAdaptiveFormat["contentLength"].(string), 64)
			videoSize /= (1024 * 1024)

			videoFileType := ""
			if strings.Contains(formatType, "mp4") {
				videoFileType = "MP4"
			} else if strings.Contains(formatType, "webm") {
				videoFileType = "WEBM"
			}

			youtubeVideoDetails.VideoQualitiesMap[currentAdaptiveFormat["qualityLabel"].(string)+" | "+fmt.Sprintf("%d", int(currentAdaptiveFormat["fps"].(float64)))+" fps | "+fmt.Sprintf("%f", videoSize)+" MB | "+videoFileType] = currentAdaptiveFormat["url"].(string)
		}
	}
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
