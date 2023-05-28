
Command line app written in [Go](https://go.dev) to download videos using YouTube URLs.

---

## Important Note:

- At present, this app do not work for videos having _age-restrictions_.

## Build Binary:

From project root enter command,

- On Linux/UNIX :

    `go build -o go_youtube_downloader ./cmd/go_youtube_downloader.go`

- On Windows :

    `go build -o go_youtube_downloader.exe ./cmd/go_youtube_downloader.go`

## Run Directly with Source Code:

- Linux/UNIX/Windows :

    `go run ./cmd/go_youtube_downloader.go`

---

## Special Thanks to Resources:

- [Reverse-Engineering YouTube: Revisited](https://tyrrrz.me/blog/reverse-engineering-youtube-revisited) - Blog explaining YouTube's internal APIs.

- [YouTube-Internal-Clients](https://github.com/zerodytrash/YouTube-Internal-Clients) - A python script that discovers hidden YouTube API clients. Just a research project.

- [YT-DLP](https://github.com/yt-dlp/yt-dlp) - A youtube-dl fork with additional features and fixes

- [Youtubei](https://github.com/SuspiciousLookingOwl/youtubei) - Get Youtube data such as videos, playlists, channels, video information & comments, related videos, up next video, and more!

- [Efficient File Download In Golang: A Comprehensive Guide](https://marketsplash.com/tutorials/go/golang-download/)
