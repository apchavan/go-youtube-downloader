
# Go YouTube Downloader

Command line app written in [Go](https://go.dev) to download videos using YouTube URLs.

## Build Binary:

From project root enter command,

- On Linux/UNIX :

    `go build -o go_youtube_downloader ./cmd/go_youtube_downloader.go`

- On Windows :

    `go build -o go_youtube_downloader.exe ./cmd/go_youtube_downloader.go`

https://github.com/apchavan/go-youtube-downloader/assets/49102443/6ded3ee7-c5ed-49d3-bdf4-57147fef3c18

## Run Directly with Source Code:

After installing [Go](https://go.dev), clone/download this project & from project root enter below command,

`go run ./cmd/go_youtube_downloader.go`

## Important Notes:

- Systems must have [FFmpeg](https://ffmpeg.org/) installed or have static binary in project directory to merge downloaded separate video & audio streams into a single file.

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
