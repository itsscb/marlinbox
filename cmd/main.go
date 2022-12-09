package main

import (
//	"time"
	marlinbox "github.com/itsscb/marlinbox"
)

func main() {
	mb := marlinbox.New("playlist.json")
	mb.Run()
/*	for {
		time.Sleep(time.Second * 1)
	}
*/
}
