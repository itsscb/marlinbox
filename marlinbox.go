package marlinbox

import (
	"encoding/json"
	"log"
	"os"
	"time"

	evdev "github.com/gvalkov/golang-evdev"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

type MarlinBox struct {
	DeviceName      string             `json:"devicename"`
	DevicePath      string             `json:"devicepath,omitempty"`
	ConfigPath      string             `json:"configpath,omitempty"`
	CurrentID       string             `json:"-"`
	Volume          float64            `json:"-"`
	NextSong        bool               `json:"-"`
	Playlist        []*PlayCard        `json:"playlist,omitempty"`
	ControlCards    []*ControlCard     `json:"controlcards,omitempty"`
	Device          *evdev.InputDevice `json:"-"`
	CurrentPlayCard *PlayCard          `json:"-"`
	Player          oto.Player         `json:"-"`
	PlayerContext   *oto.Context       `json:"-"`
}

type RFIDCard struct {
	ID string `json:"id"`
}

type PlayCard struct {
	RFIDCard
	File []string `json:"file"`
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
	if err != nil {
		log.Panicf("Error listing Input Devices: %s", err)
	}

	for _, dev := range devices {
		if dev.Name == mb.DeviceName {
			mb.DevicePath = dev.Fn
		}
	}

	if mb.DevicePath == "" {
		panic("Device not found")
	}

	mb.ConfigPath = path

	return mb
}

func (mb *MarlinBox) Run() {
	var err error
	mb.Device, err = evdev.Open(mb.DevicePath)
	if err != nil {
		panic(err)
	}
	mb.Device.Grab()

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

			if val == "ENTER" || len(mb.CurrentID) == 9 {
				if val != "ENTER" {
					mb.CurrentID += val
				}
				mb.GetCurrentCard()
				mb.CurrentID = ""
				continue
			}

			mb.CurrentID += val
		}
	}
}

func (mb *MarlinBox) GetCurrentCard() {
	for _, c := range mb.ControlCards {
		if mb.CurrentID == c.ID {
			switch c.Function {
			case "vol+":
				if mb.Volume < 1.0 {
					mb.Volume = mb.Volume + 0.2
					return
				}
			case "vol-":
				if mb.Volume >= 0.2 {
					mb.Volume = mb.Volume - 0.2
					return
				}
			case "nextSong":
				if mb.Player != nil {
					mb.NextSong = true
					return
				}
				return
			default:
				return
			}
		}
	}
	for _, c := range mb.Playlist {
		if mb.CurrentID == c.ID {
			if mb.CurrentPlayCard != nil {
				if mb.CurrentID == mb.CurrentPlayCard.ID {
					return
				}
			}
			mb.CurrentPlayCard = c
			if mb.PlayerContext != nil {
				mb.Player.Reset()
			}
			mb.Play()
			return
		}
	}
	mb.Playlist = append(mb.Playlist, &PlayCard{
		RFIDCard: RFIDCard{
			ID: mb.CurrentID,
		},
	})

	config, err := json.MarshalIndent(&mb, "", " ")
	if err != nil {
		log.Printf("\nCould not marshal output '%+v': %s", config, err)
	}

	jsonFile, err := os.Create(mb.ConfigPath)
	if err != nil {
		log.Printf("\nCould not create json-File '%s': %s", mb.ConfigPath, err)
	}

	_, err = jsonFile.Write(config)
	if err != nil {
		log.Printf("\nCould not write json-File '%s': %s", mb.ConfigPath, err)
	}
}

func (mb *MarlinBox) Play() {
	go func() {
		var ready chan struct{}
		playingID := mb.CurrentPlayCard.ID

		for _, file := range mb.CurrentPlayCard.File {
			if playingID != mb.CurrentPlayCard.ID {
				return
			}
			f, err := os.Open(file)
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

				if mb.NextSong {
					mb.Player.Pause()
					mb.Player.Close()
					mb.NextSong = false
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
		}
	}()

}
