package runner

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Checks if the FFmpeg is installed in the system.
// If installed, then it'll run the commands to merge provided video & audio files
// to create combined single output file.
// Otherwise, it'll show status message to get FFmpeg &
// then manually run the required command.
//
// Reference : https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
func CheckAndMergeWithFFmpeg(
	videoFilePath string,
	audioFilePath string,
	videoTitle string,
	mergeProgressMsgChannel chan string,
) {
	// Make sure `mergeProgressMsgChannel` will be closed
	defer close(mergeProgressMsgChannel)

	ffmpegBinaryName := "ffmpeg"
	execPath, err := exec.LookPath(ffmpegBinaryName)

	if errors.Is(err, exec.ErrDot) {
		// If FFmpeg binary is found in project root directory,
		// then get absolute path of that binary.
		execPath, _ = filepath.Abs(execPath)
		err = nil
	} else if err != nil {
		mergeProgressMsgChannel <- GetFFmpegNotFoundText(videoFilePath, audioFilePath)
		return
	}

	// Remove special characters from name & create a file name to write merged output data.
	videoTitle = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(videoTitle, "")
	outputFileName := strings.TrimSpace(videoTitle) + ".mp4"

	mergeProgressMsgChannel <- "🔃 Merging video & audio files with FFmpeg..."

	// Create FFmpeg command of below form to execute :
	// ffmpeg -i 'audio.mp4' -i 'video.mp4' -c:a aac -c:v libx265 -preset ultrafast 'output.mp4'
	// Reference : https://stackoverflow.com/a/45688183
	cmdArguments := []string{
		"-i", audioFilePath,
		"-i", videoFilePath,
		"-c:a", "aac",
		"-c:v", "libx265",
		"-preset", "ultrafast",
		outputFileName,
	}

	ffmpegCmd := exec.Command(
		execPath,
		cmdArguments...,
	)

	ffmpegStartTime := time.Now()

	// Execute the command
	err = ffmpegCmd.Run()
	if err != nil {
		mergeProgressMsgChannel <- fmt.Sprintf("❌ Unable to merge with FFmpeg.\nerr : %v\n", err)
		return
	}
	ffmpegEndTime := time.Now()

	// Delete separately downloaded video file & audio file once merged successfully...
	os.Remove(videoFilePath)
	os.Remove(audioFilePath)

	mergeProgressMsgChannel <- (GetFFmpegMergeSuccessText(outputFileName) +
		fmt.Sprintf("\nTime taken to merge : %v", ffmpegEndTime.Sub(ffmpegStartTime)))
}
