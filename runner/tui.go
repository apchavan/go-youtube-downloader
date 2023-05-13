package runner

import (
	"github.com/rivo/tview"
)

// Construct & return the application instance.
func GetTuiAppLayout() *tview.Application {
	// Create pointer object of `YouTubeVideoOperations` struct
	youtubeVideoOperations := &YouTubeVideoOperations{VideoUrl: "", VideoQuality: ""}

	// Create application object
	application := tview.NewApplication()

	// Create flex layout
	flexLayout := tview.NewFlex()

	// Get input form
	inputForm := GetInputForm(application, youtubeVideoOperations)

	// Get info form
	infoForm := GetInfoForm(application, youtubeVideoOperations)

	// Add both forms to the `flexLayout`
	flexLayout = flexLayout.AddItem(
		tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(inputForm, 0, 2, true).
			AddItem(infoForm, 0, 1, false),
		0, 1, true)

	// Run the application with mouse support enabled
	if err := application.SetRoot(flexLayout, true).SetFocus(flexLayout).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
	return application
}

// Create & return input form
func GetInputForm(application *tview.Application, youtubeVideoOperations *YouTubeVideoOperations) *tview.Form {
	inputForm := tview.NewForm()

	inputForm = inputForm.AddInputField(
		GetVideoLinkInputFieldLabel(), "", 0, nil,
		func(text string) {
			youtubeVideoOperations.VideoUrl = text

			// TODO: Fetch, check video's metadata & set the quality options array of string if video is valid
			if text != "" && inputForm.GetFormItemIndex(GetVideoQualityDropdownLabel()) == -1 {
				// Create the Video Quality dropdown when proper link text is present in the input field
				// and no existing input dropdown is already exist.
				inputForm = inputForm.AddDropDown(
					GetVideoQualityDropdownLabel(),
					[]string{"144p", "240p", "480p", "720p", "1080p"},
					0,

					func(option string, optionIndex int) {
						youtubeVideoOperations.VideoQuality = option
					})

				// Add `Download` button if already not present.
				if inputForm.GetButtonIndex(GetDownloadButtonLabel()) == -1 {
					inputForm = inputForm.AddButton(
						GetDownloadButtonLabel(),

						func() {
							if youtubeVideoOperations.VideoUrl != "" && youtubeVideoOperations.VideoQuality != "" {
								// Clear previous status textviews if found any
								for inputForm.GetFormItemIndex(GetStatusLabel()) != -1 {
									inputForm.RemoveFormItem(inputForm.GetFormItemIndex(GetStatusLabel()))
								}

								// Before starting download, change the label & set button disabled
								inputForm.GetButton(inputForm.GetButtonIndex(GetDownloadButtonLabel())).
									SetLabel(GetDownloadButtonProgressLabel()).SetDisabled(true)

								// Add & set textview to show download in progress...
								inputForm = inputForm.AddTextView(GetStatusLabel(),
									GetDownloadingMessage(youtubeVideoOperations.VideoUrl),
									0, 0, true, true)

								// Very important to re-draw the screen before actually starting the video download
								application.ForceDraw().Sync()

								// TODO: Handle download action
								downloadStatusChannel := make(chan bool)
								go youtubeVideoOperations.DownloadYouTubeVideo(downloadStatusChannel)

								downloadStatus := false
								for value := range downloadStatusChannel {
									downloadStatus = value
								}

								// Change textview to show download finished...
								if downloadStatus {
									// Clear previous status textviews if found any
									for inputForm.GetFormItemIndex(GetStatusLabel()) != -1 {
										inputForm.RemoveFormItem(inputForm.GetFormItemIndex(GetStatusLabel()))
									}
									inputForm = inputForm.AddTextView(GetStatusLabel(), GetDownloadFinishedMessage(),
										0, 0, true, true)

									// In the end, reset the download button.
									// Before starting download, change the label & set button disabled
									inputForm.GetButton(inputForm.GetButtonIndex(GetDownloadButtonProgressLabel())).
										SetLabel(GetDownloadButtonLabel()).SetDisabled(false)
								}
							}
						}).SetButtonsAlign(tview.AlignRight)
				}
			} else if text == "" && inputForm.GetFormItemIndex(GetVideoQualityDropdownLabel()) != -1 {
				// Remove the Video Quality dropdown when no text is present in the input field
				// and atleast 1 existing input dropdown is already exist.
				inputForm.RemoveFormItem(inputForm.GetFormItemIndex(GetVideoQualityDropdownLabel()))

				// Clear previous status textviews if found any
				for inputForm.GetFormItemIndex(GetStatusLabel()) != -1 {
					inputForm.RemoveFormItem(inputForm.GetFormItemIndex(GetStatusLabel()))
				}

				// Remove download button when no text is present.
				if inputForm.GetButtonIndex(GetDownloadButtonLabel()) != -1 {
					inputForm.RemoveButton(inputForm.GetButtonIndex(GetDownloadButtonLabel()))
				}
			}
		})

	// Add `Quit` button to input form
	/*
		inputForm = inputForm.AddButton(
			GetQuitButtonLabel(),

			func() {
				application.Stop()
				application.ForceDraw().Sync()
			}).SetButtonsAlign(tview.AlignRight)
	*/

	inputForm = inputForm.SetFocus(inputForm.GetFormItemIndex(GetVideoLinkInputFieldLabel()))
	inputForm.SetBorder(true).SetTitle(GetAppNameTitle()).SetTitleAlign(tview.AlignCenter)
	return inputForm
}

// Create & return info form
func GetInfoForm(application *tview.Application, youtubeVideoOperations *YouTubeVideoOperations) *tview.Form {
	infoForm := tview.NewForm()

	// Add info about how to use the app
	infoForm = infoForm.AddTextView(
		GetUsageInfoLabel(),
		GetUsageInfoText(),
		0, 1, true, true)

	// Add info about quitting the app
	infoForm = infoForm.AddTextView(
		GetQuitInfoLabel(),
		GetQuitInfoText(),
		0, 1, true, true)

	// Add GitHub repo link
	infoForm = infoForm.AddTextView(
		GetRepoLabel(),
		GetRepoLinkText(),
		0, 1, true, true)

	infoForm.SetBorder(true).SetTitle(GetAboutAppTitle()).SetTitleAlign(tview.AlignCenter)
	return infoForm
}
