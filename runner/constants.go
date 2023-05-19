package runner

import "fmt"

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
func GetDownloadButtonProgressLabel() string {
	return "Download in progress..."
}

// Return label for Quit button.
func GetQuitButtonLabel() string {
	return "Quit"
}

// Return label for Status.
func GetStatusLabel() string {
	return "Status : "
}

// Return status message for Downloading Video.
func GetDownloadingMessage(downloadingVideoTitle string) string {
	return fmt.Sprintf("Downloading '%s'...", downloadingVideoTitle)
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
	var line3 string = "(3) Then click on newly appeared 'Download' button."
	return (line1 + line2 + line3)
}

// Return label for Quit info.
func GetQuitInfoLabel() string {
	return "Quit app using : "
}

// Return string text for Quit info.
func GetQuitInfoText() string {
	return "Keyboard shortcut '<Ctrl> + c'"
}

// Return label for age-restricted note.
func GetAgeRestrictedNoteLabel() string {
	return "Note : "
}

// Return string text for age-restricted note.
func GetAgeRestrictedText() string {
	return "Age-Restricted videos are NOT supported..."
}
