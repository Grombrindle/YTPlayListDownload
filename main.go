package main

import (
	YT "YTPlayListDownload/playlist"
	"log"

	"github.com/gen2brain/beeep"
	// "os"
	// "time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	YT.OpenLink()
	err := beeep.Notify("YTPlayListDownload", "Program has finished", "")
	if err != nil {
		log.Fatal(err)
	}
}
