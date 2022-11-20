package marlinbox

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	evdev "github.com/gvalkov/golang-evdev"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

type MarlinBox struct {
	DeviceName      string `json:"devicename"`
	DevicePath      string `json:"devicepath,omitempty"`
	CurrentID       string
	Volume          float64        `json:"volume,omitempty"`
	Playlist        []*PlayCard    `json:"playlist,omitempty"`
	ControlCards    []*ControlCard `json:"controlcards,omitempty"`
	Device          *evdev.InputDevice
	CurrentPlayCard *PlayCard
	Player          oto.Player
	PlayerContext   *oto.Context
}

type RFIDCard struct {
	ID string `json:"id"`
}

type PlayCard struct {
	RFIDCard
	File string `json:"file,omitempty"`
}
type ControlCard struct {
	RFIDCard
	Function string `json:"function,omitempty"`
}

var KeyMap = map[uint16]string{
	2:  "1",
	3:  "2",
	4:  "3",
	5:  "4",
	6:  "5",
	7:  "6",
	8:  "7",
	9:  "8",
	10: "9",
	11: "0",
	28: "ENTER",
}

func New(path string) *MarlinBox {
	var mb *MarlinBox

	f, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(f, &mb)
	if err != nil {
		panic(err)
	}

	if mb.Volume == 0 {
		mb.Volume = 1
	}

	if mb.DeviceName == "" {
		panic("No DeviceName given")
	}

	devices, err := evdev.ListInputDevices()

	for _, dev := range devices {
		if dev.Name == mb.DeviceName {
			mb.DevicePath = dev.Fn
		}
	}

	if mb.DevicePath == "" {
		panic("Device not found")
	}

	return mb
}

func (mb *MarlinBox) Run() {
	var err error
	mb.Device, err = evdev.Open(mb.DevicePath)
	if err != nil {
		panic(err)
	}
	mb.Device.Grab()

	go func() {
		for {
			events, err := mb.Device.Read()
			if err != nil {
				log.Println(err)
			}

			for _, ev := range events {
				if ev.Value != 0x1 {
					continue
				}
				val, ok := KeyMap[ev.Code]
				if !ok {
					continue
				}

				if val == "ENTER" {
					err := mb.GetCurrentCard()
					if err != nil {
						log.Println(err)
					}
					fmt.Println(mb.CurrentPlayCard)
					mb.CurrentID = ""
					// err =
					// if err != nil {
					// 	log.Println(err)
					// 	panic(err)
					// }
					continue
				}

				mb.CurrentID += val
			}
		}
	}()
}

func (mb *MarlinBox) GetCurrentCard() error {
	fmt.Println(mb.ControlCards)
	for _, c := range mb.ControlCards {
		if mb.CurrentID == c.ID {
			switch c.Function {
			case "vol+":
				fmt.Println("vol+")
				if mb.Volume < 1.0 {
					mb.Volume = mb.Volume + 0.2
				}
			case "vol-":
				fmt.Println("vol-")
				if mb.Volume > 0 {
					mb.Volume = mb.Volume - 0.2
				}
			default:
				return nil
			}
			return nil
		}
	}
	for _, c := range mb.Playlist {
		if mb.CurrentID == c.ID {
			mb.CurrentPlayCard = c
			if mb.PlayerContext != nil {
				mb.Player.Reset()
			}
			mb.Play()
			return nil
		}
	}
	return errors.New("Card not found: " + mb.CurrentID)
}

func (mb *MarlinBox) Play() {
	go func() {
		var err error
		var ready chan struct{}
		playingID := mb.CurrentPlayCard.ID
		f, err := os.Open(mb.CurrentPlayCard.File)
		if err != nil {
			log.Println(err)
			return
		}
		defer f.Close()

		d, err := mp3.NewDecoder(f)
		if err != nil {
			log.Println(err)
			return
		}

		mb.PlayerContext, ready, err = oto.NewContext(d.SampleRate(), 2, 2)
		if err != nil {
			log.Println(err)
			return
		}
		<-ready

		mb.Player = mb.PlayerContext.NewPlayer(d)
		defer mb.Player.Close()
		mb.Player.SetVolume(mb.Volume)
		mb.Player.Play()

		for {
			if mb.CurrentPlayCard.ID != playingID {
				break
			}
			if mb.Volume != mb.Player.Volume() {
				mb.Player.SetVolume(mb.Volume)
			}
			time.Sleep(time.Second)
			if !mb.Player.IsPlaying() {
				break
			}
		}
	}()

}
