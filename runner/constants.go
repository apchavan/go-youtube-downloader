package runner

import "fmt"

// Returns the string of application's name for title.
func GetAppNameTitle() string {
	return " Go YouTube Downloader "
}

// Return label for `Video Link` input field.
func GetVideoLinkInputFieldLabel() string {
	return "Paste/Type YouTube Video Link : "
}

// Return label for `Video Quality` dropdown.
func GetVideoQualityDropdownLabel() string {
	return "Select Video Quality : "
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

// Return status message for Download Finished.
func GetDownloadFinishedMessage() string {
	return "Download Finished!"
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
	return "Paste/type the YouTube video link in above input field. Then click on appeared 'Download' button."
}

// Return label for Quit info.
func GetQuitInfoLabel() string {
	return "How to quit ? : "
}

// Return string text for Quit info.
func GetQuitInfoText() string {
	return "Press keyboard shortcut <Ctrl + c>"
}
