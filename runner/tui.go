package runner

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

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

			if urlText != "" {
				videoValidityMsgChannel := make(chan string)
				isValidVideo := false
				videoValidityMsg := ""

				// Fetch & check video's metadata to see if video is valid
				go youtubeVideoDetails.IsValidYouTubeURL(videoValidityMsgChannel, &isValidVideo)

				for msg := range videoValidityMsgChannel {
					videoValidityMsg = msg
				}

				if !isValidVideo {
					// Clear controls if already exist
					removeControls(inputForm, true, true, true, true, true)

					// Show status message stored by `videoValidityMsg`
					inputForm = inputForm.AddTextView(GetStatusLabel(),
						videoValidityMsg,
						0, 0, true, true)
					return
				} else {
					// Clear previous status TextViews if found any
					removeControls(inputForm, false, false, false, true, false)
				}

				// Show video title in TextView
				videoDetailsMap := youtubeVideoDetails.VideoMetaData["videoDetails"]
				videoTitle := videoDetailsMap.(map[string]interface{})["title"].(string)
				inputForm = inputForm.AddTextView(GetYouTubeVideoTitleLabel(),
					videoTitle,
					0, 0, true, true)

				// Create the Video dropdown when proper link text is present in the input field
				// and no existing input dropdown is already exist.
				if inputForm.GetFormItemIndex(GetVideoQuality_FPS_Size_Type_DropdownLabel()) == -1 {
					inputForm = inputForm.AddDropDown(
						GetVideoQuality_FPS_Size_Type_DropdownLabel(),
						getDescendingSize_VideoQualities(youtubeVideoDetails.VideoQualitiesMap),
						-1,

						func(option string, optionIndex int) {
							youtubeVideoDetails.SelectedVideoQuality = option
						})
				}

				// Create the Audio dropdown when proper link text is present in the input field
				// and no existing input dropdown is already exist.
				if inputForm.GetFormItemIndex(GetAudioBitrate_Size_Type_DropdownLabel()) == -1 {
					inputForm = inputForm.AddDropDown(
						GetAudioBitrate_Size_Type_DropdownLabel(),
						getDescendingSize_AudioQualities(youtubeVideoDetails.AudioQualitiesMap),
						-1,

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
								/* 1. Handle video downloading */

								// Clear previous status TextViews if found any
								removeControls(inputForm, false, false, false, true, false)

								// Before starting download, change the label & set button disabled
								inputForm.GetButton(inputForm.GetButtonIndex(GetDownloadButtonLabel())).
									SetLabel(GetDownloadButtonProgressLabel()).SetDisabled(true)

								// Add & set TextView to show download in progress...
								inputForm = inputForm.AddTextView(GetStatusLabel(),
									GetStartingDownloadMessage(videoTitle),
									0, 0, true, true)

								// Very important to re-draw the screen before actually starting the video download
								application.ForceDraw().Sync()

								// Handle video file download action
								downloadProgressMsgChannel := make(chan string)
								isDownloadFinished := false
								go youtubeVideoDetails.DownloadYouTubeVideoFile(
									downloadProgressMsgChannel,
									&isDownloadFinished,
									application,
								)

								// Clear previous status TextViews if found any
								removeControls(inputForm, false, false, false, true, false)

								// Set TextView for showing download status
								videoDownloadProgressTextView := tview.NewTextView().
									SetLabel(GetStatusLabel()).
									SetText(fmt.Sprintf("Downloading video file for %s", videoTitle)).
									SetChangedFunc(func() {
										// Update the screen as per download progress
										application.ForceDraw().Sync()
									})

								// Add TextView to show video download progress
								inputForm = inputForm.AddFormItem(videoDownloadProgressTextView)

								// Check the download progress & update the TextView
								for downloadProgressMsg := range downloadProgressMsgChannel {
									videoDownloadProgressTextView.SetText(downloadProgressMsg)
								}

								// Change TextView to show download finished...
								if isDownloadFinished {
									// Clear previous status TextViews if found any
									removeControls(inputForm, false, false, false, true, false)

									// Pass download finished message with video title
									inputForm = inputForm.AddTextView(GetStatusLabel(),
										GetDownloadFinishedMessage(videoTitle),
										0, 0, true, true)

									// In the end, reset the download button.
									// Before starting download, change the label & set button disabled
									// inputForm.GetButton(inputForm.GetButtonIndex(GetDownloadButtonProgressLabel())).
									// 	SetLabel(GetDownloadButtonLabel()).SetDisabled(false)
								}

								/* 2. Handle audio downloading */
								// Handle video file download action
								downloadProgressMsgChannel = make(chan string)
								isDownloadFinished = false
								go youtubeVideoDetails.DownloadYouTubeAudioFile(
									downloadProgressMsgChannel,
									&isDownloadFinished,
									application,
								)

								// Clear previous status TextViews if found any
								removeControls(inputForm, false, false, false, true, false)

								// Set TextView for showing download status
								audioDownloadProgressTextView := tview.NewTextView().
									SetLabel(GetStatusLabel()).
									SetText(fmt.Sprintf("Downloading audio file for %s", videoTitle)).
									SetChangedFunc(func() {
										// Update the screen as per download progress
										application.ForceDraw().Sync()
									})

								// Add TextView to show video download progress
								inputForm = inputForm.AddFormItem(audioDownloadProgressTextView)

								// Check the download progress & update the TextView
								for downloadProgressMsg := range downloadProgressMsgChannel {
									audioDownloadProgressTextView.SetText(downloadProgressMsg)
								}

								// Change TextView to show download finished...
								if isDownloadFinished {
									// Clear previous status TextViews if found any
									removeControls(inputForm, false, false, false, true, false)

									// Pass download finished message with video title
									inputForm = inputForm.AddTextView(GetStatusLabel(),
										GetDownloadFinishedMessage(videoTitle),
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
				removeControls(inputForm, true, true, true, true, true)
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
	removeVideoDropdown, removeAudioDropdown, removeYouTubeVideoTitleTextView, removeStatusTextView, removeDownloadButton bool) {
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

	// Clear previous YouTube Video Title TextViews if found any
	if removeYouTubeVideoTitleTextView {
		for inputForm.GetFormItemIndex(GetYouTubeVideoTitleLabel()) != -1 {
			inputForm.RemoveFormItem(inputForm.GetFormItemIndex(GetYouTubeVideoTitleLabel()))
		}
	}

	// Clear previous status TextViews if found any
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

// Returns the video qualities in descending order based on size.
func getDescendingSize_AudioQualities(audioQualitiesMap map[string]string) []string {
	if len(audioQualitiesMap) == 0 {
		return make([]string, 0)
	}

	audioSizesCollection := make([]float64, 0)
	for option := range audioQualitiesMap {
		// fmt.Printf("\n\nV- %s => %s\n", option, url)
		for idx, subString := range strings.Split(option, " | ") {
			// Index 1 is audio size
			if idx == 1 {
				// fmt.Printf("%v) -- %v\n", idx, strings.Split(v, " "))
				floatSize, _ := strconv.ParseFloat(strings.Split(subString, " ")[0], 64)
				audioSizesCollection = append(
					audioSizesCollection,
					floatSize,
				)
			}
		}
	}
	// fmt.Printf("\nBefore :- %v\n", audioSizesCollection)
	sort.Sort(sort.Reverse(sort.Float64Slice(audioSizesCollection)))
	// fmt.Printf("\nAfter :- %v\n", audioSizesCollection)

	descendingSizeAudioQualities := make([]string, 0)
	for _, floatSize := range audioSizesCollection {
		sizeString := fmt.Sprint(floatSize)
		for key := range audioQualitiesMap {
			// fmt.Printf("\nstrings.Contains(%v, %v)", key, sizeString)
			if strings.Contains(key, sizeString) {
				descendingSizeAudioQualities = append(descendingSizeAudioQualities, key)
			}
		}
	}

	return descendingSizeAudioQualities
}

// Returns the video qualities in descending order based on size.
func getDescendingSize_VideoQualities(videoQualitiesMap map[string]string) []string {
	if len(videoQualitiesMap) == 0 {
		return make([]string, 0)
	}

	videoSizesCollection := make([]float64, 0)
	for option := range videoQualitiesMap {
		// fmt.Printf("\n\nV- %s => %s\n", option, url)
		for idx, subString := range strings.Split(option, " | ") {
			// Index 2 is video size
			if idx == 2 {
				// fmt.Printf("%v) -- %v\n", idx, strings.Split(v, " "))
				floatSize, _ := strconv.ParseFloat(strings.Split(subString, " ")[0], 64)
				videoSizesCollection = append(
					videoSizesCollection,
					floatSize,
				)
			}
		}
	}
	// fmt.Printf("\nBefore :- %v\n", videoSizeCollection)
	sort.Sort(sort.Reverse(sort.Float64Slice(videoSizesCollection)))
	// fmt.Printf("\nAfter :- %v\n", videoSizeCollection)

	descendingSizeVideoQualities := make([]string, 0)
	for _, floatSize := range videoSizesCollection {
		sizeString := fmt.Sprint(floatSize)
		for key := range videoQualitiesMap {
			// fmt.Printf("\nstrings.Contains(%v, %v)", key, sizeString)
			if strings.Contains(key, sizeString) {
				descendingSizeVideoQualities = append(descendingSizeVideoQualities, key)
			}
		}
	}

	return descendingSizeVideoQualities
}
