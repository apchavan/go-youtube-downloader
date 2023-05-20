package runner

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

/*
Definitions for creating request object structure needed by YouTube's internal APIs.
*/

// Structure of client in 'ContextRequestMapStruct'.
type ClientRequestMapStruct struct {
	Hl                string `json:"hl"`
	Gl                string `json:"gl"`
	ClientName        string `json:"clientName"`
	ClientVersion     string `json:"clientVersion"`
	ClientScreen      string `json:"clientScreen"`
	AndroidSdkVersion int    `json:"androidSdkVersion"`
}

// Structure of embedUrl in 'ContextRequestMapStruct'.
type ThirdPartyRequestMapStruct struct {
	EmbedUrl string `json:"embedUrl"`
}

// Structure of context in 'RequestBodyStruct'.
type ContextRequestMapStruct struct {
	Client     ClientRequestMapStruct     `json:"client"`
	ThirdParty ThirdPartyRequestMapStruct `json:"thirdParty"`
}

/*
// Structure of contentPlaybackContext in 'PlaybackContextStruct'.
type ContentPlaybackContextStruct struct {
	SignatureTimestamp int `json:"signatureTimestamp"`
}

// Structure of playbackContext in 'RequestBodyStruct'.
//
// This is mostly used when accessing YouTube videos with Age-Restrictions.
type PlaybackContextStruct struct {
	ContentPlaybackContext ContentPlaybackContextStruct `json:"contentPlaybackContext"`
}
*/

// Structure of request body needed by YouTube's internal APIs.
type RequestBodyStruct struct {
	Context ContextRequestMapStruct `json:"context"`
	// PlaybackContext PlaybackContextStruct   `json:"playbackContext"`
	VideoId        string `json:"videoId"`
	RacyCheckOk    bool   `json:"racyCheckOk"`
	ContentCheckOk bool   `json:"contentCheckOk"`
}

// Extract YouTube's video ID & return string.
func getYouTubeVideoID_FromURL(videoURL string) string {
	videoURL = strings.TrimSpace(videoURL)

	// If total number of characters is 11, then it's likely to be a video ID itself.
	if len(videoURL) == 11 {
		return videoURL
	}

	// Extract video's ID using video URL.
	videoURL = strings.TrimSuffix(videoURL, "/")

	videoURL = strings.ReplaceAll(videoURL, "https://www.youtube.com/watch?v=", "")

	videoURL = strings.ReplaceAll(videoURL, "https://youtu.be/", "")

	return videoURL
}

// Function that sends a POST request to YouTube's internal API endpoint, YouTubei,
// then returns the fetched metadata map to the caller.
func GetVideoMetadataFromYouTubei(videoURL string) map[string]interface{} {

	// Extract ID from `videoURL`
	videoID := getYouTubeVideoID_FromURL(videoURL)

	// Create HTTP client
	client := &http.Client{}

	// Create request body object
	reqBodyObj := &RequestBodyStruct{
		Context: ContextRequestMapStruct{
			Client: ClientRequestMapStruct{
				Hl:                "en",
				Gl:                "US",
				ClientName:        "TVHTML5_SIMPLY_EMBEDDED_PLAYER",
				ClientVersion:     "2.0",
				ClientScreen:      "WATCH",
				AndroidSdkVersion: 33,
			},
			ThirdParty: ThirdPartyRequestMapStruct{
				EmbedUrl: "https://www.youtube.com/",
			},
		},
		/*
			PlaybackContext: PlaybackContextStruct{
				ContentPlaybackContext: ContentPlaybackContextStruct{
					SignatureTimestamp: 19487,
				},
			},
		*/
		VideoId:        videoID,
		RacyCheckOk:    true,
		ContentCheckOk: true,
	}

	// Create a bytes buffer to write & store the `reqBodyObj`
	reqBodyBytesBuffer := new(bytes.Buffer)

	// Write JSON encoding of `reqBodyObj` to `reqBodyBytesBuffer`
	json.NewEncoder(reqBodyBytesBuffer).Encode(reqBodyObj)

	// Create a strings.Reader object using `reqBodyBytesBuffer` to pass with POST request
	reqBodyStringsReader := strings.NewReader(reqBodyBytesBuffer.String())

	// The URL for YouTube's internal API with key
	youtubeiUrl := "https://www.youtube.com/youtubei/v1/player?key=AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8"

	// Prepare POST request object
	postReq, err := http.NewRequest("POST", youtubeiUrl, reqBodyStringsReader)
	if err != nil {
		log.Fatal(err)
	}

	// Set required headers
	postReq.Header.Set("Host", "www.youtube.com")
	postReq.Header.Set("Content-Type", "application/json")
	postReq.Header.Set("Accept", "*/*")
	postReq.Header.Set("Origin", "https://www.youtube.com")
	postReq.Header.Set("Referer", "https://www.youtube.com/")

	// Send POST request
	resp, err := client.Do(postReq)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read the response body
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	responseBodyMap := make(map[string]interface{})
	json.Unmarshal(bodyText, &responseBodyMap)

	return responseBodyMap
}
