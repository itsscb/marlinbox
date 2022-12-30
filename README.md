# **MARLIN** Box
This project aims to take another approach at RFID-Jukeboxes like the [toniebox](https://tonies.com/).

Our goal is to empower parents to make their kids happy - on their own terms.

## Prerequisites
The project is currently in the development but is already used in production.

If you wish to try the **MARLIN** Box in this state you will need the following:
- [Raspberry Pi Zero 2W](https://www.berrybase.de/detail/index/sArticle/9357)
  - Raspbian OS
    - Packages:
      - `portaudio19-dev`
- Powerbank
- [USB-Hub](https://www.amazon.de/gp/product/B01K7RR3W8/ref=ppx_yo_dt_b_asin_title_o04_s01?ie=UTF8&psc=1)
- [Speakers](https://www.amazon.de/gp/product/B00JRW0M32/ref=ppx_yo_dt_b_asin_title_o04_s00?ie=UTF8&psc=1)
- [USB-Soundcard](https://www.berrybase.de/usb-2.0-soundkarte-mit-stereo-kopfhoerer-ausgang-und-mikrofon-eingang)
- [USB-RFID-Reader (**EM4100**)](https://www.amazon.de/gp/product/B018OYOR3E/ref=ppx_yo_dt_b_asin_title_o05_s01?ie=UTF8&psc=1)
- [RFID-Cards](https://www.amazon.de/gp/product/B07TRSR3VB/ref=ppx_yo_dt_b_asin_title_o04_s02?ie=UTF8&psc=1)
- Some cables
- Some **.mp3-Files**
- *Optional*:
  - [Power-Button](https://www.amazon.de/gp/product/B08VH4SMLT/ref=ppx_yo_dt_b_asin_title_o05_s00?ie=UTF8&psc=1)

## Setup
- Setup your **Raspberry Pi Zero 2 W** as usual or use the [Raspberry Pi Imager](https://www.raspberrypi.com/software/) on your Computer and insert the configured SD-Card into your **Raspberry Pi Zero 2 W**
- Once online and connected to your network *ssh* into the **Raspberry**
- Update and Restart the Pi using `sudo apt-get -y update && sudo apt-get -y upgrade && sudo reboot now` 
- *ssh* back
- Install `curl` and `portaudio19-dev` if it's not already there using `sudo apt-get -y install curl portaudio19-dev`
- Download **golang** using `curl -OL https://go.dev/dl/go1.19.4.linux-armv6l.tar.gz`
- Extract the archive using `tar -xf go.19.4.linux-arm6l.tar.gz`
- Clone this repository
  - `git clone https://github.com/itsscb/marlinbox`
  - `cd marlinbox/cmd`
- Build the software `~/go/bin/go build main.go && mv cmd marlinbox`
- Create or modify the file `playlist.json` in the same directory
- Create the `systemd service` with `sudo nano /etc/systemd/system/marlinbox.service`
  - `[Unit]
Description=Runs the binary of the marlinbox

[Service]
ExecStart=/home/marlinbox/marlinbox/cmd/cmd
WorkingDirectory=/home/marlinbox/marlinbox/cmd
Type=simple
User=root
`
- Enable the `systemd service` with `sudo systemctl enable marlinbox.service`
- Restart the Pi with `sudo reboot now`

Your **MARLINBOX** should be up and running.
