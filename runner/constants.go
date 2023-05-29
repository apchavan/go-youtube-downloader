package runner

import "fmt"

// Returns the byte value 83886.08 for Range header in GET request
// i.e. almost equal to 0.08 MB.
func GetStepSizeBytes() float64 {
	return float64(83886.08)
}

// Returns the string of application's name for title.
func GetAppNameTitle() string {
	return " Go YouTube Downloader "
}

// Return label for `Video ID/Link` input field.
func GetVideoIDLink_InputFieldLabel() string {
	return "Paste/Type YouTube Video ID/Link : "
}

// Return label for video related dropdown.
func GetVideoQuality_FPS_Size_Type_DropdownLabel() string {
	return "Video Quality, FPS, Size & Type : "
}

// Return label for audio related dropdown.
func GetAudioBitrate_Size_Type_DropdownLabel() string {
	return "Audio Bitrate, Size & Type : "
}

// Return label for Download button.
func GetDownloadButtonLabel() string {
	return "Download"
}

// Return label for Download button when download has been started.
func GetDownloadButtonProcessingLabel() string {
	return "Processing..."
}

// Return label for Quit button.
func GetQuitButtonLabel() string {
	return "Quit"
}

// Return label for Status.
func GetStatusLabel() string {
	return "Status : "
}

// Return label for YouTube video title.
func GetYouTubeVideoTitleLabel() string {
	return "YouTube Video Title : "
}

// Return status message for Starting download.
func GetStartingDownloadMessage(downloadingVideoTitle string) string {
	return fmt.Sprintf("Starting download for '%s'...", downloadingVideoTitle)
}

// Return status message for Invalid Video URL.
func GetInvalidURLMessage(videoURL string) string {
	return fmt.Sprintf("No YouTube video found at :> '%s'", videoURL)
}

// Return status message for Download Finished.
func GetDownloadFinishedMessage(downloadingVideoTitle string) string {
	return fmt.Sprintf("Download Finished for '%s'!", downloadingVideoTitle)
}

// Returns the string of info title.
func GetAboutAppTitle() string {
	return " About App "
}

// Return label for GitHub Repo.
func GetRepoLabel() string {
	return "GitHub Repo : "
}

// Return string text for GitHub Repo link.
func GetRepoLinkText() string {
	return "https://github.com/apchavan/go-youtube-downloader"
}

// Return label for app usage info.
func GetUsageInfoLabel() string {
	return "How to use ? : "
}

// Return string text for app usage info.
func GetUsageInfoText() string {
	var line1 string = "(1) Paste/type the YouTube video ID or link in above input field.\n"
	var line2 string = "(2) Select the required 'Video Quality, FPS, Size & Type' & 'Audio Bitrate, Size & Type' from their dropdowns.\n"
	var line3 string = "(3) Then click on newly appeared 'Download' button. Also visit link mentioned in 'Important Notes & Limitations' below."
	return (line1 + line2 + line3)
}

// Return label for Quit info.
func GetQuitInfoLabel() string {
	return "How to Quit ? : "
}

// Return string text for Quit info.
func GetQuitInfoText() string {
	return "Keyboard shortcut '<Ctrl> + c' or if downloading started then close Terminal directly"
}

// Return label for important notes/limitations.
func GetImportantNotesLimitsLabel() string {
	return "Important Notes & Limitations : "
}

// Return string text for important notes/limitations.
func GetImportantNotesLimitsText() string {
	return "https://github.com/apchavan/go-youtube-downloader#important-notes"
}

// Return string text for FFmpeg not found in system.
func GetFFmpegNotFoundText(
	videoFilePath string,
	audioFilePath string,
) string {
	var line1 string = "❌ FFmpeg not found in environment's PATH.\n"
	var line2 string = "(1) Install latest FFmpeg package from https://ffmpeg.org/.\n"
	var line3 string = "(2) Then try below command with absolute file paths to manually merge video & audio : \n"
	var line4 string = "ffmpeg -i \"audio_file_path\" -i \"video_file_path\" -c:a aac -c:v libx265 -preset ultrafast \"output_file_path\"\n"
	return (line1 + line2 + line3 + line4)
}

// Return string text for files successfully merged to output using FFmpeg.
func GetFFmpegMergeSuccessText(
	outputFileName string,
) string {
	return fmt.Sprintf("🥳🎉 Final output stored as '%s'", outputFileName)
}
