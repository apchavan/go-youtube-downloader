package runner

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rivo/tview"
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

// Function for 'YouTubeVideoDetailsStruct' that downloads the audio file from YouTube.
func (youtubeVideoDetails *YouTubeVideoDetailsStruct) DownloadYouTubeAudioFile(
	downloadProgressMsgChannel chan string,
	isDownloadFinished *bool,
	application *tview.Application,
) {
	// Get video title
	videoDetailsMap := youtubeVideoDetails.VideoMetaData["videoDetails"]
	videoTitle := videoDetailsMap.(map[string]interface{})["title"].(string)

	// Video file extension
	fileExtension := ""

	// Total download size in MB
	totalDownloadSizeMB := 0.0

	for idx, subString := range strings.Split(youtubeVideoDetails.SelectedAudioQuality, " | ") {
		// Index 1 is video size
		if idx == 1 {
			totalDownloadSizeMB, _ = strconv.ParseFloat(strings.Split(subString, " ")[0], 64)

		}
		// Index 2 is audio file extension
		if idx == 2 {
			fileExtension = strings.TrimSpace(subString)
		}
	}

	// Remove special characters from name
	videoTitle = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(videoTitle, "")
	// Generate a file name to save video
	fileName := strings.TrimSpace(videoTitle) + "_AF_AUD." + strings.ToLower(fileExtension)

	// Create blank file to save video data
	audioFile, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("\n(*) Error while creating the audio file... \n%v\n", err)
	}
	defer audioFile.Close()

	// Set download URL
	downloadUrl := youtubeVideoDetails.VideoQualitiesMap[youtubeVideoDetails.SelectedAudioQuality]

	// Get total download size in bytes to set range header while downloading
	totalDownloadSizeBytes := totalDownloadSizeMB * (1024 * 1024)

	// Create HTTP client for GET request
	httpClient := &http.Client{}

	getReq, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		log.Fatalf("\n(*) Error while creating GET request object... \n%v\n", err)
	}

	// Below stepSizeBytes is almost equal to 1.25 MB
	stepSizeBytes := float64(1310720)
	byteAdditionFactor := stepSizeBytes - 1.0
	totalBytesCopied := float64(0.0)

	// Loop using Range header to download `stepSizeBytes` bytes chunks of video data in `fileName` from `downloadUrl`
	for idx := float64(0.0); idx < totalDownloadSizeBytes; idx += stepSizeBytes {
		getStartTime := time.Now()
		// Set header to get range of bytes
		getReq.Header.Set("Range", fmt.Sprintf("bytes=%.0f-%.0f", idx, idx+byteAdditionFactor))

		respVideoData, err := httpClient.Do(getReq)
		if err != nil {
			log.Fatalf("\n(*) Error while fetching video data... \n%v\n", err)
		}
		defer respVideoData.Body.Close()
		getEndTime := time.Now()

		// Copy/store the downloaded video data
		// currentBytesCopied, err := io.Copy(videoFile, respVideoData.Body)
		copyStartTime := time.Now()
		currentBytesCopied, err := io.CopyN(audioFile, respVideoData.Body, respVideoData.ContentLength)
		if err != nil {
			log.Fatalf("\n(*) Error while copying video data... \n%v\n", err)
		}
		copyEndTime := time.Now()

		totalBytesCopied += float64(currentBytesCopied)

		downloadProgressMsg := fmt.Sprintf(
			"%f MB audio file downloaded of %f MB (Time taken: Download=%v, Storage Write=%v)",
			totalBytesCopied/(1024*1024),
			totalDownloadSizeMB,
			getEndTime.Sub(getStartTime),
			copyEndTime.Sub(copyStartTime),
		)
		downloadProgressMsgChannel <- downloadProgressMsg
	}

	// Set the downloading status as completed
	*isDownloadFinished = true

	// Close channel to indicate task is completed & this channel will not provide anything in future
	close(downloadProgressMsgChannel)
}

// Function for 'YouTubeVideoDetailsStruct' that downloads the video file from YouTube.
func (youtubeVideoDetails *YouTubeVideoDetailsStruct) DownloadYouTubeVideoFile(
	downloadProgressMsgChannel chan string,
	isDownloadFinished *bool,
	application *tview.Application,
) {
	// Get video title
	videoDetailsMap := youtubeVideoDetails.VideoMetaData["videoDetails"]
	videoTitle := videoDetailsMap.(map[string]interface{})["title"].(string)

	// Video file extension
	fileExtension := ""

	// Total download size in MB
	totalDownloadSizeMB := 0.0

	for idx, subString := range strings.Split(youtubeVideoDetails.SelectedVideoQuality, " | ") {
		// Index 2 is video size
		if idx == 2 {
			totalDownloadSizeMB, _ = strconv.ParseFloat(strings.Split(subString, " ")[0], 64)

		}
		// Index 3 is video file extension
		if idx == 3 {
			fileExtension = strings.TrimSpace(subString)
		}
	}

	// Remove special characters from name
	videoTitle = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(videoTitle, "")
	// Generate a file name to save video
	fileName := strings.TrimSpace(videoTitle) + "_AF_VID." + strings.ToLower(fileExtension)

	// Create blank file to save video data
	videoFile, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("\n(*) Error while creating the video file... \n%v\n", err)
	}
	defer videoFile.Close()

	// Set download URL
	downloadUrl := youtubeVideoDetails.VideoQualitiesMap[youtubeVideoDetails.SelectedVideoQuality]

	// Get total download size in bytes to set range header while downloading
	totalDownloadSizeBytes := totalDownloadSizeMB * (1024 * 1024)

	// Create HTTP client for GET request
	httpClient := &http.Client{}

	getReq, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		log.Fatalf("\n(*) Error while creating GET request object... \n%v\n", err)
	}

	// stepSizeBytes := uint64(10000000)
	// stepSizeBytes := float64(524288)
	// stepSizeBytes := float64(10240)

	// Below stepSizeBytes is almost equal to 1.25 MB
	stepSizeBytes := float64(1310720)
	byteAdditionFactor := stepSizeBytes - 1.0
	totalBytesCopied := float64(0.0)

	// Loop using Range header to download `stepSizeBytes` bytes chunks of video data in `fileName` from `downloadUrl`
	for idx := float64(0.0); idx < totalDownloadSizeBytes; idx += stepSizeBytes {
		getStartTime := time.Now()
		// Set header to get range of bytes
		getReq.Header.Set("Range", fmt.Sprintf("bytes=%.0f-%.0f", idx, idx+byteAdditionFactor))

		respVideoData, err := httpClient.Do(getReq)
		if err != nil {
			log.Fatalf("\n(*) Error while fetching video data... \n%v\n", err)
		}
		defer respVideoData.Body.Close()
		getEndTime := time.Now()

		// Copy/store the downloaded video data
		// currentBytesCopied, err := io.Copy(videoFile, respVideoData.Body)
		copyStartTime := time.Now()
		currentBytesCopied, err := io.CopyN(videoFile, respVideoData.Body, respVideoData.ContentLength)
		if err != nil {
			log.Fatalf("\n(*) Error while copying video data... \n%v\n", err)
		}
		copyEndTime := time.Now()

		totalBytesCopied += float64(currentBytesCopied)

		downloadProgressMsg := fmt.Sprintf(
			"%f MB video file downloaded of %f MB (Time taken: Download=%v, Storage Write=%v)",
			totalBytesCopied/(1024*1024),
			totalDownloadSizeMB,
			getEndTime.Sub(getStartTime),
			copyEndTime.Sub(copyStartTime),
		)
		downloadProgressMsgChannel <- downloadProgressMsg
	}

	// Set the downloading status as completed
	*isDownloadFinished = true

	// Close channel to indicate task is completed & this channel will not provide anything in future
	close(downloadProgressMsgChannel)
}
