# **MARLIN** Box
This project aims to take another approach at RFID-Jukeboxes like the [toniebox](https://tonies.com/).

Our goal is to empower parents to make their kids happy - on their own terms.

## Prerequisites
We are currently at the beginning of the development.

If you wish to try the **MARLIN** Box in this state you will need the following:
- Raspberry Pi (at least 3B)
  - Raspbian OS
    - Packages:
      - `portaudio-devel.x86_64`
- RFID-Reader (**EM4100**)
- Some **.mp3-Files**

## Setup
- Clone this repository to your Computer or directly onto the Raspberry Pi - on whatever device you would like to build the software you will need **golang** installed
  - `git clone https://github.com/itsscb/marlinbox`
  - `cd marlinbox/cmd`
- Build the software
  - `go build main.go && mv main marlinbox`
- Create or modify the file `playlist.json` in the same directory 
- Run it with `sudo` [^1]
  - `sudo ./marlinbox`

[^1]: Required because we access `/dev/input/` devices which require root privileges.
