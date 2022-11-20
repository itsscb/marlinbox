package main

import (
	"time"

	marlinbox "github.com/itsscb/marlinbox"
)

func main() {
	mb := marlinbox.New("playlist.json")
	mb.Run()
	time.Sleep(time.Minute * 5)
}
