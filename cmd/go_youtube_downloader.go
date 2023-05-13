package main

import (
	"fmt"
	"strings"

	runnerpkg "github.com/apchavan/go-youtube-downloader/runner"
)

func main() {
	_ = runnerpkg.GetTuiAppLayout()
	fmt.Printf("\n- '%s' exited...\n", strings.TrimSpace(runnerpkg.GetAppNameTitle()))
}
