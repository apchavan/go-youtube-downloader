package runner

import (
	"github.com/rivo/tview"
)

// Construct & return the application instance.
func GetTuiAppLayout() *tview.Application {
	// Create pointer object of `YouTubeVideoDetailsStruct` struct
	youtubeVideoDetails := &YouTubeVideoDetailsStruct{}

	// Create application object
	application := tview.NewApplication()

	// Create flex layout
	flexLayout := tview.NewFlex()

	// Get input form
	inputForm := getInputForm(application, youtubeVideoDetails)

	// Get info form
	infoForm := getInfoForm()

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
func getInputForm(application *tview.Application, youtubeVideoDetails *YouTubeVideoDetailsStruct) *tview.Form {
	inputForm := tview.NewForm()

	inputForm = inputForm.AddInputField(
		GetVideoIDLink_InputFieldLabel(), "", 0, nil,
		func(urlText string) {
			youtubeVideoDetails.VideoUrlOrID = urlText

			// TODO: Fetch, check video's metadata & set the quality & FPS options array of string if video is valid
			if urlText != "" {
				if !youtubeVideoDetails.IsValidYouTubeURL() {
					// Clear controls if already exist
					removeControls(inputForm, true, true, true, true)

					// TODO: Pass download finished message with video title
					inputForm = inputForm.AddTextView(GetStatusLabel(),
						GetInvalidURLMessage(youtubeVideoDetails.VideoUrlOrID),
						0, 0, true, true)
					return
				} else {
					// Clear previous status textviews if found any
					removeControls(inputForm, false, false, true, false)
				}

				// Create the Video dropdown when proper link text is present in the input field
				// and no existing input dropdown is already exist.
				// TODO: Pass actual video qualities highest first, lowest last
				videoQualities := make([]string, len(youtubeVideoDetails.VideoQualitiesMap))
				idx := 0
				for qualityKey := range youtubeVideoDetails.VideoQualitiesMap {
					videoQualities[idx] = qualityKey
					idx++
				}
				if inputForm.GetFormItemIndex(GetVideoQuality_FPS_Size_Type_DropdownLabel()) == -1 {
					inputForm = inputForm.AddDropDown(
						GetVideoQuality_FPS_Size_Type_DropdownLabel(),
						videoQualities,
						/*[]string{
							" 1440p60 | 60 fps | 68.6 MB | MP4 ",
							" 1080p60 HDR | 60 fps | 58.6 MB | WEBM ",
							" 720p60 | 60 fps | 48.6 MB | MP4 ",
							" 480p | 30 fps | 38.6 MB | WEBM ",
							" 360p | 24 fps | 28.6 MB | MP4 ",
							" 144p | 24 fps | 18.6 MB | WEBM ",
						},*/
						0,

						func(option string, optionIndex int) {
							youtubeVideoDetails.SelectedVideoQuality = option
						})
				}

				// Create the Audio dropdown when proper link text is present in the input field
				// and no existing input dropdown is already exist.
				// TODO: Pass actual audio qualities highest first, lowest last
				audioQualities := make([]string, len(youtubeVideoDetails.AudioQualitiesMap))
				idx = 0
				for qualityKey := range youtubeVideoDetails.AudioQualitiesMap {
					audioQualities[idx] = qualityKey
					idx++
				}
				if inputForm.GetFormItemIndex(GetAudioBitrate_Size_Type_DropdownLabel()) == -1 {
					inputForm = inputForm.AddDropDown(
						GetAudioBitrate_Size_Type_DropdownLabel(),
						audioQualities,
						/*[]string{
							" 631073 | 68.6 MB | MP4 ",
							" 531073 | 58.6 MB | WEBM ",
							" 431073 | 48.6 MB | MP4 ",
							" 331073 | 38.6 MB | WEBM ",
							" 231073 | 28.6 MB | MP4 ",
							" 131073 | 18.6 MB | WEBM ",
						},*/
						0,

						func(option string, optionIndex int) {
							youtubeVideoDetails.SelectedAudioQuality = option
						})
				}

				// Add `Download` button if already not present.
				if inputForm.GetButtonIndex(GetDownloadButtonLabel()) == -1 {
					inputForm = inputForm.AddButton(
						GetDownloadButtonLabel(),

						func() {
							if youtubeVideoDetails.VideoUrlOrID != "" &&
								youtubeVideoDetails.SelectedVideoQuality != "" &&
								youtubeVideoDetails.SelectedAudioQuality != "" {
								// Clear previous status textviews if found any
								removeControls(inputForm, false, false, true, false)

								// Before starting download, change the label & set button disabled
								inputForm.GetButton(inputForm.GetButtonIndex(GetDownloadButtonLabel())).
									SetLabel(GetDownloadButtonProgressLabel()).SetDisabled(true)

								// Add & set textview to show download in progress...
								// TODO: Pass downloading message with video title
								inputForm = inputForm.AddTextView(GetStatusLabel(),
									GetDownloadingMessage(youtubeVideoDetails.VideoUrlOrID),
									0, 0, true, true)

								// Very important to re-draw the screen before actually starting the video download
								application.ForceDraw().Sync()

								// TODO: Handle download action
								downloadStateChannel := make(chan bool)
								go youtubeVideoDetails.DownloadYouTubeVideo(downloadStateChannel)

								// Check the status whether download has finished
								isDownloadFinished := false
								for value := range downloadStateChannel {
									isDownloadFinished = value
								}

								// Change textview to show download finished...
								if isDownloadFinished {
									// Clear previous status textviews if found any
									removeControls(inputForm, false, false, true, false)

									// TODO: Pass download finished message with video title
									inputForm = inputForm.AddTextView(GetStatusLabel(),
										GetDownloadFinishedMessage(youtubeVideoDetails.VideoUrlOrID),
										0, 0, true, true)

									// In the end, reset the download button.
									// Before starting download, change the label & set button disabled
									inputForm.GetButton(inputForm.GetButtonIndex(GetDownloadButtonProgressLabel())).
										SetLabel(GetDownloadButtonLabel()).SetDisabled(false)
								}
							}
						}).SetButtonsAlign(tview.AlignRight)
				}
			} else if urlText == "" {
				// Clear controls if already exist
				removeControls(inputForm, true, true, true, true)
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

	inputForm = inputForm.SetFocus(inputForm.GetFormItemIndex(GetVideoIDLink_InputFieldLabel()))
	inputForm.SetBorder(true).SetTitle(GetAppNameTitle()).SetTitleAlign(tview.AlignCenter)
	return inputForm
}

// Create & return info form
func getInfoForm() *tview.Form {
	infoForm := tview.NewForm()

	// Add info about how to use the app
	infoForm = infoForm.AddTextView(
		GetUsageInfoLabel(),
		GetUsageInfoText(),
		0, 3, true, true)

	// Add info about quitting the app
	infoForm = infoForm.AddTextView(
		GetQuitInfoLabel(),
		GetQuitInfoText(),
		0, 1, true, true)

	// Add note about age-restricted videos not supported
	infoForm = infoForm.AddTextView(
		GetAgeRestrictedNoteLabel(),
		GetAgeRestrictedText(),
		0, 1, true, true)

	// Add GitHub repo link
	infoForm = infoForm.AddTextView(
		GetRepoLabel(),
		GetRepoLinkText(),
		0, 1, true, true)

	infoForm.SetBorder(true).SetTitle(GetAboutAppTitle()).SetTitleAlign(tview.AlignCenter)
	return infoForm
}

// Clears the controls specified by boolean parameters
func removeControls(
	inputForm *tview.Form,
	removeVideoDropdown, removeAudioDropdown, removeStatusTextView, removeDownloadButton bool) {
	// Remove the Video dropdown when no text is present in the input field
	// and atleast 1 existing input dropdown is already exist.
	if removeVideoDropdown {
		for inputForm.GetFormItemIndex(GetVideoQuality_FPS_Size_Type_DropdownLabel()) != -1 {
			inputForm.RemoveFormItem(inputForm.GetFormItemIndex(GetVideoQuality_FPS_Size_Type_DropdownLabel()))
		}
	}

	// Remove the Audio dropdown when no text is present in the input field
	// and atleast 1 existing input dropdown is already exist.
	if removeAudioDropdown {
		for inputForm.GetFormItemIndex(GetAudioBitrate_Size_Type_DropdownLabel()) != -1 {
			inputForm.RemoveFormItem(inputForm.GetFormItemIndex(GetAudioBitrate_Size_Type_DropdownLabel()))
		}
	}

	// Clear previous status textviews if found any
	if removeStatusTextView {
		for inputForm.GetFormItemIndex(GetStatusLabel()) != -1 {
			inputForm.RemoveFormItem(inputForm.GetFormItemIndex(GetStatusLabel()))
		}
	}

	// Remove download button when no text is present.
	if removeDownloadButton {
		for inputForm.GetButtonIndex(GetDownloadButtonLabel()) != -1 {
			inputForm.RemoveButton(inputForm.GetButtonIndex(GetDownloadButtonLabel()))
		}
	}
}
