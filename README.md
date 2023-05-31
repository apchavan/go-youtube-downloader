
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/gomods/athens) ![](https://img.shields.io/badge/OS-Linux-orange) ![](https://img.shields.io/badge/OS-macOS-black) ![](https://img.shields.io/badge/OS-Windows-blue) [![GoReportCard](https://goreportcard.com/badge/github.com/apchavan/go-youtube-downloader)](https://goreportcard.com/report/github.com/apchavan/go-youtube-downloader)

# Go YouTube Downloader

Command line app written in [Go](https://go.dev) to download Shorts & Videos using YouTube URLs/IDs.

## Main Features:

- YouTube Shorts & Videos downloading.

- Ability to select from different content qualities.

## Project Dependencies:

At present, the project has 2 dependencies,

1. [tview](https://github.com/rivo/tview) - Terminal UI library with rich, interactive widgets - written in Golang.

2. [FFmpeg](https://ffmpeg.org/) - The leading cross-platform multimedia framework. It should be either installed in system or have latest static binary in project's root directory.

## Working Demo:

1. When pasted either YouTube Video or Shorts ID/URL the app fetches the metadata from YouTube's internal APIs.

2. Then depending on quality selections for video & audio, the application downloads the Video/Shorts content by making of small sized data requests to the fetched content URLs.

3. In the end, if the [FFmpeg](https://ffmpeg.org/) exist, then both separate video & audio stream files are merged into single output file.

https://github.com/apchavan/go-youtube-downloader/assets/49102443/e177f755-b607-40be-8d22-05f4850e97a7

## Build Binary:

From project root enter command,

- On Linux/UNIX :-

    `go build -o go_youtube_downloader ./cmd/go_youtube_downloader.go`

- On Windows :-

    `go build -o go_youtube_downloader.exe ./cmd/go_youtube_downloader.go`

https://github.com/apchavan/go-youtube-downloader/assets/49102443/6ded3ee7-c5ed-49d3-bdf4-57147fef3c18

## Run Directly with Source Code:

After installing [Go](https://go.dev), clone/download this project & from project root enter below command,

`go run ./cmd/go_youtube_downloader.go`

https://github.com/apchavan/go-youtube-downloader/assets/49102443/1b1c4fd9-f0fe-4590-86e1-4a456a012d5f

## Important Notes:

- Systems must have [FFmpeg](https://ffmpeg.org/) installed or have latest static binary in project's root directory to merge downloaded separate video & audio streams into a single file.

- _Age-restricted_ videos can not be downloaded due to YouTube's Signature Ciphering.

- YouTube have bandwidth limitations for each incoming request, around 10 MB per request. If any request gets more data than this size limit, then further requests will throttle download or connection may get terminated. So, to get better performance when downloading data & writing it to output file, it's divided into smaller chunks for consistency. Based on selected quality & size, the download time would be more or less.

## Special Thanks to Resources:

- [Reverse-Engineering YouTube: Revisited](https://tyrrrz.me/blog/reverse-engineering-youtube-revisited) - Blog explaining YouTube's internal APIs.

- [YouTube-Internal-Clients](https://github.com/zerodytrash/YouTube-Internal-Clients) - A python script that discovers hidden YouTube API clients. Just a research project.

- [YT-DLP](https://github.com/yt-dlp/yt-dlp) - A youtube-dl fork with additional features and fixes

- [Youtubei](https://github.com/SuspiciousLookingOwl/youtubei) - Get Youtube data such as videos, playlists, channels, video information & comments, related videos, up next video, and more!

- [Efficient File Download In Golang: A Comprehensive Guide](https://marketsplash.com/tutorials/go/golang-download/)

- [Golang Download Files Example](https://golangdocs.com/golang-download-files)

- [Team NewPipe](https://github.com/TeamNewPipe)
