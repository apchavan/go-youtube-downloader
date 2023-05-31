package runner

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rivo/tview"
)

// Structure that holds required details about YouTube Shorts/Video.
type YouTubeDetailsStruct struct {
	VideoUrlOrID            string                 // Store the YouTube's video URL or ID
	SelectedVideoQuality    string                 // Store the user's selected video quality
	SelectedAudioQuality    string                 // Store the user's selected audio quality
	VideoMetaData           map[string]interface{} // Store the metadata of `VideoUrlOrID` fetched from YouTube
	VideoQualitiesMap       map[string]string      // Store map of displayed video qualities to their download URLs
	AudioQualitiesMap       map[string]string      // Store map of displayed audio qualities to their download URLs
	DownloadedVideoFilePath string                 // Store file path where video is downloaded
	DownloadedAudioFilePath string                 // Store file path where audio is downloaded
}

// Function to validate & return a boolean status of whether the passed URL is right.
func (youtubeDetails *YouTubeDetailsStruct) IsValidYouTubeURL(
	videoValidityMsgChannel chan string,
	isValidVideo *bool,
) {
	// Make sure the `videoValidityMsgChannel` will be closed
	defer close(videoValidityMsgChannel)

	// Perform basic check if YouTube itself is missing from video URL,
	// or total characters != 11 (to consider as video ID);
	// if yes, then it's invalid URL!
	if !strings.Contains(
		strings.ToLower(youtubeDetails.VideoUrlOrID),
		"youtube.com",
	) && !strings.Contains(
		strings.ToLower(youtubeDetails.VideoUrlOrID),
		"youtu.be",
	) && len(youtubeDetails.VideoUrlOrID) != 11 {
		videoValidityMsgChannel <- GetInvalidURLMessage(youtubeDetails.VideoUrlOrID)
		*isValidVideo = false
		return
	}

	// Send POST request to YouTube's internal API to get & store map of metadata in `VideoMetaData`
	youtubeDetails.VideoMetaData = GetVideoMetadataFromYouTubei(youtubeDetails.VideoUrlOrID)

	// If 'status' under 'playabilityStatus' is not set to "OK" then video is not available
	responseBodyMap := youtubeDetails.VideoMetaData
	playabilityStatusMap := responseBodyMap["playabilityStatus"]

	// Check if video status is not set to 'OK', return false
	if playabilityStatusMap.(map[string]interface{})["status"].(string) != "OK" {
		videoValidityMsgChannel <- GetInvalidURLMessage(youtubeDetails.VideoUrlOrID)
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
			videoDetailsMap := youtubeDetails.VideoMetaData["videoDetails"]
			videoTitle := videoDetailsMap.(map[string]interface{})["title"].(string)

			videoValidityMsgChannel <- ("WARNING : '" + videoTitle + "' is marked as Age-Restricted video & so can't be downloaded...")
			*isValidVideo = false
			return
		}
	}

	// Extract & store video/audio information highest first, lowest last order
	youtubeDetails.extractAdaptiveAudioVideoQualities()

	videoValidityMsgChannel <- ""
	*isValidVideo = true
}

// Extract & store video/audio information highest first, lowest last order.
func (youtubeDetails *YouTubeDetailsStruct) extractAdaptiveAudioVideoQualities() {
	responseBodyMap := youtubeDetails.VideoMetaData
	streamingDataMap := responseBodyMap["streamingData"]

	// Allocate respective maps
	youtubeDetails.AudioQualitiesMap = make(map[string]string)
	youtubeDetails.VideoQualitiesMap = make(map[string]string)

	// Iterate over 'adaptiveFormats' key & get the required data.
	for _, rawAdaptiveFormat := range streamingDataMap.(map[string]interface{})["adaptiveFormats"].([]interface{}) {
		currentAdaptiveFormat := rawAdaptiveFormat.(map[string]interface{})

		formatType := currentAdaptiveFormat["mimeType"].(string)

		if strings.HasPrefix(formatType, "audio") {
			// Process the discovered audio format.
			// Audio Bitrate, Size & Type : 131073 bits/s | 5862305 (Converted to MB) | MP4/WEBM

			// Convert contentLength from bytes to MB
			audioSize, _ := strconv.ParseFloat(currentAdaptiveFormat["contentLength"].(string), 64)
			audioSize /= (1024 * 1024)

			audioFileType := ""
			if strings.Contains(formatType, "mp4") {
				audioFileType = "MP4"
			} else if strings.Contains(formatType, "webm") {
				audioFileType = "WEBM"
			}

			youtubeDetails.AudioQualitiesMap[fmt.Sprintf("%d", int(currentAdaptiveFormat["bitrate"].(float64)))+" bits/s | "+fmt.Sprintf("%f", audioSize)+" MB | "+audioFileType] = currentAdaptiveFormat["url"].(string)

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

			youtubeDetails.VideoQualitiesMap[currentAdaptiveFormat["qualityLabel"].(string)+" | "+fmt.Sprintf("%d", int(currentAdaptiveFormat["fps"].(float64)))+" fps | "+fmt.Sprintf("%f", videoSize)+" MB | "+videoFileType] = currentAdaptiveFormat["url"].(string)
		}
	}
}

// Function for 'YouTubeDetailsStruct' that downloads the audio file from YouTube.
func (youtubeDetails *YouTubeDetailsStruct) DownloadYouTubeAudioFile(
	downloadProgressMsgChannel chan string,
	isDownloadFinished *bool,
	application *tview.Application,
) {
	// Make sure `downloadProgressMsgChannel` will be closed
	defer close(downloadProgressMsgChannel)

	// Get video title
	videoDetailsMap := youtubeDetails.VideoMetaData["videoDetails"]
	audioTitle := videoDetailsMap.(map[string]interface{})["title"].(string)

	// Video file extension
	fileExtension := ""

	// Total download size in MB
	totalDownloadSizeMB := 0.0

	for idx, subString := range strings.Split(youtubeDetails.SelectedAudioQuality, " | ") {
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
	audioTitle = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(audioTitle, "")
	// Generate a file name to save video
	fileName := strings.TrimSpace(audioTitle) + "_AF_AUD." + strings.ToLower(fileExtension)

	// Create blank file to save video data
	audioFile, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("\n(*) Error while creating the audio file... \n%v\n", err)
	}
	defer audioFile.Close()

	// Set download URL
	downloadUrl := youtubeDetails.AudioQualitiesMap[youtubeDetails.SelectedAudioQuality]

	// Get total download size in bytes to set range header while downloading
	totalDownloadSizeBytes := totalDownloadSizeMB * (1024 * 1024)

	// Create HTTP client for GET request
	httpClient := &http.Client{}

	getReq, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		log.Fatalf("\n(*) Error while creating GET request object... \n%v\n", err)
	}
	// getReq.Close = true // https://stackoverflow.com/a/19006050

	stepSizeBytes := GetStepSizeBytes()
	byteAdditionFactor := stepSizeBytes - 1.0
	totalBytesCopied := float64(0.0)

	downloadStartTime := time.Now()

	// Loop using Range header to download `stepSizeBytes` bytes chunks of audio data in `fileName` from `downloadUrl`
	for idx := float64(0.0); idx < totalDownloadSizeBytes; idx += stepSizeBytes {

		// Set header to get range of bytes
		getReq.Header.Set("Range", fmt.Sprintf("bytes=%.0f-%.0f", idx, idx+byteAdditionFactor))

		respAudioData, err := httpClient.Do(getReq)
		if err != nil {
			// Sometimes the request get failed, then create new request object and
			// again try to continue fetching the video bytes
			getReq, err = http.NewRequest("GET", downloadUrl, nil)

			// Handle the worst case by getting new `downloadUrl`
			// and then re-construct new request object
			if err != nil {
				videoValidityMsgChannel := make(chan string)
				isValidVideo := false
				youtubeDetails.IsValidYouTubeURL(
					videoValidityMsgChannel,
					&isValidVideo,
				)
				<-videoValidityMsgChannel

				// Set download URL
				downloadUrl := youtubeDetails.AudioQualitiesMap[youtubeDetails.SelectedAudioQuality]

				getReq, _ = http.NewRequest("GET", downloadUrl, nil)
			}

			// To match with correct bytes range, subtract `stepSizeBytes`
			// so it'll more precise to missing bytes when continuing loop
			idx -= stepSizeBytes

			downloadProgressMsgChannel <- "Re-created new HTTP request object..."
			continue
		}

		// Copy/store the downloaded audio data
		fileOpStartTime := time.Now()

		// audioFile.ReadFrom(respAudioData.Body)
		// io.CopyN(audioFile, respAudioData.Body, respAudioData.ContentLength)
		io.Copy(audioFile, respAudioData.Body)

		fileOpEndTime := time.Now()

		totalBytesCopied += float64(respAudioData.ContentLength)
		downloadProgressMsg := fmt.Sprintf(
			"Downloading audio file; done %.2f MB of %.2f MB...\nAudio data file writing time: %v",
			totalBytesCopied/(1024*1024),
			totalDownloadSizeMB,
			fileOpEndTime.Sub(fileOpStartTime),
		)

		downloadProgressMsgChannel <- downloadProgressMsg

		// Close the response body before starting next iteration.
		// This is more performant than using `defer` statement.
		respAudioData.Body.Close()
	}

	downloadEndTime := time.Now()

	downloadProgressMsgChannel <- fmt.Sprintf(
		"Written audio file as '%s'...\nDownload time : %v\n",
		audioFile.Name(),
		downloadEndTime.Sub(downloadStartTime),
	)

	// Store the file path where audio is downloaded.
	youtubeDetails.DownloadedAudioFilePath, _ = filepath.Abs(audioFile.Name())

	// Sleep for 5 seconds just to make sure above message is displayed.
	time.Sleep(time.Second * 5)

	// Set the downloading status as completed
	*isDownloadFinished = true
}

// Function for 'YouTubeDetailsStruct' that downloads the video file from YouTube.
func (youtubeDetails *YouTubeDetailsStruct) DownloadYouTubeVideoFile(
	downloadProgressMsgChannel chan string,
	isDownloadFinished *bool,
	application *tview.Application,
) {
	// Make sure `downloadProgressMsgChannel` will be closed
	defer close(downloadProgressMsgChannel)

	// Get video title
	videoDetailsMap := youtubeDetails.VideoMetaData["videoDetails"]
	videoTitle := videoDetailsMap.(map[string]interface{})["title"].(string)

	// Video file extension
	fileExtension := ""

	// Total download size in MB
	totalDownloadSizeMB := 0.0

	for idx, subString := range strings.Split(youtubeDetails.SelectedVideoQuality, " | ") {
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
	downloadUrl := youtubeDetails.VideoQualitiesMap[youtubeDetails.SelectedVideoQuality]

	// Get total download size in bytes to set range header while downloading
	totalDownloadSizeBytes := totalDownloadSizeMB * (1024 * 1024)

	// Create HTTP client for GET request
	httpClient := &http.Client{}

	getReq, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		log.Fatalf("\n(*) Error while creating GET request object... \n%v\n", err)
	}
	// getReq.Close = true // https://stackoverflow.com/a/19006050

	stepSizeBytes := GetStepSizeBytes()
	byteAdditionFactor := stepSizeBytes - 1.0
	totalBytesCopied := float64(0.0)

	downloadStartTime := time.Now()

	// Loop using Range header to download `stepSizeBytes` bytes chunks of video data in `fileName` from `downloadUrl`
	for idx := float64(0.0); idx < totalDownloadSizeBytes; idx += stepSizeBytes {

		// Set header to get range of bytes
		getReq.Header.Set("Range", fmt.Sprintf("bytes=%.0f-%.0f", idx, idx+byteAdditionFactor))

		respVideoData, err := httpClient.Do(getReq)
		if err != nil {
			// Sometimes the request get failed, then create new request object and
			// again try to continue fetching the video bytes
			getReq, err = http.NewRequest("GET", downloadUrl, nil)

			// Handle the worst case by getting new `downloadUrl`
			// and then re-construct new request object
			if err != nil {
				videoValidityMsgChannel := make(chan string)
				isValidVideo := false
				youtubeDetails.IsValidYouTubeURL(
					videoValidityMsgChannel,
					&isValidVideo,
				)
				<-videoValidityMsgChannel

				// Set download URL
				downloadUrl := youtubeDetails.VideoQualitiesMap[youtubeDetails.SelectedVideoQuality]

				getReq, _ = http.NewRequest("GET", downloadUrl, nil)
			}

			// To match with correct bytes range, subtract `stepSizeBytes`
			// so it'll more precise to missing bytes when continuing loop
			idx -= stepSizeBytes

			downloadProgressMsgChannel <- "Re-created new HTTP request object..."
			continue
		}
		// Copy/store the downloaded video data in file
		fileOpStartTime := time.Now()

		// videoFile.ReadFrom(respVideoData.Body)
		// io.CopyN(videoFile, respVideoData.Body, respVideoData.ContentLength)
		io.Copy(videoFile, respVideoData.Body)

		fileOpEndTime := time.Now()

		totalBytesCopied += float64(respVideoData.ContentLength)
		downloadProgressMsg := fmt.Sprintf(
			"Downloading video file; done %.2f MB of %.2f MB...\nVideo data file writing time: %v",
			totalBytesCopied/(1024*1024),
			totalDownloadSizeMB,
			fileOpEndTime.Sub(fileOpStartTime),
		)
		downloadProgressMsgChannel <- downloadProgressMsg

		// Close the response body before starting next iteration.
		// This is more performant than using `defer` statement.
		respVideoData.Body.Close()
	}

	downloadEndTime := time.Now()

	downloadProgressMsgChannel <- fmt.Sprintf(
		"Written video file as '%s'...\nDownload time : %v\n",
		videoFile.Name(),
		downloadEndTime.Sub(downloadStartTime),
	)

	// Store the file path where video is downloaded.
	youtubeDetails.DownloadedVideoFilePath, _ = filepath.Abs(videoFile.Name())

	// Sleep for 5 seconds just to make sure above message is displayed.
	time.Sleep(time.Second * 5)

	// Set the downloading status as completed
	*isDownloadFinished = true
}
