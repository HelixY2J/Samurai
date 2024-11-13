package main

import (
	"fmt"
	"net"
	"os"

	"github.com/yobert/alsa"
)

// const (
// 	StreamTypePlayback = C.SND_PCM_STREAM_PLAYBACK // playback stream
// 	//StreamTypeCapture  = C.SND_PCM_STREAM_CAPTURE
// )

func main() {

	// connect to the server

	conn, err := net.Dial("tcp", "localhost:8087")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}

	defer conn.Close()
	fmt.Println("Connected to server")

	cards, err := alsa.OpenCards()
	if err != nil {
		fmt.Println("Error opening ALSA device:", err)
		os.Exit(1)
	}
	defer alsa.CloseCards(cards)

	defaultCard := cards[1]
	fmt.Println("def:", defaultCard)

	// Step 3: Access devices on the chosen card
	device, err := defaultCard.Devices()
	if err != nil || len(device) == 0 {
		fmt.Println("Error accessing devices:", err)
		os.Exit(1)
	}
	fmt.Println("devices", device)
	// We will choose the first available playback device for simplicity
	playbackDevice := device[0]

	// Open the playback device for audio output
	if err := playbackDevice.Open(); err != nil {
		fmt.Println("Error opening device:", err)
		os.Exit(1)
	}
	defer playbackDevice.Close()
	fmt.Println("playback", playbackDevice)

	// Step 3: Set hardware parameters (16-bit PCM, 48000 Hz, stereo)
	_, err = playbackDevice.NegotiateFormat(alsa.S16_LE) // 16-bit little-endian format
	if err != nil {
		fmt.Println("Error setting format:", err)
		os.Exit(1)
	}

	_, err = playbackDevice.NegotiateRate(48000) // Set sample rate to 48000 Hz
	if err != nil {
		fmt.Println("Error setting rate:", err)
		os.Exit(1)
	}

	_, err = playbackDevice.NegotiateChannels(2) // Stereo output
	if err != nil {
		fmt.Println("Error setting channels:", err)
		os.Exit(1)
	}

	// Step 4: Prepare the device
	if err := playbackDevice.Prepare(); err != nil {
		fmt.Println("Error preparing device:", err)
		os.Exit(1)
	}

	// At this point, you should be ready to play audio to the default device.
	fmt.Println("Device is ready, you can start playing audio.")

	fmt.Println("----- Card Details -----")
	fmt.Printf("Cards: %s\n", cards)
	fmt.Printf("def: %s\n", defaultCard)

	// Device details
	fmt.Println("----- Device Details -----")

	fmt.Printf("Devices: %s\n", device)
	fmt.Printf("play: %s\n", playbackDevice)

	// Step 6: Simulate reading audio from a network connection
	// buf_buf := make([]byte, 65536*20)
	buf := make([]byte, 32768)
	for {
		n, err := conn.Read(buf)

		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("EOF reached. Server closed the connection.")
				break
			}
			fmt.Println("Error reading data from server:", err)
			break
		}
		if err := playbackDevice.Write(buf[:n], n); err != nil {
			fmt.Println("Error wriiting data to device: ", err)
			break
		}
		fmt.Printf("Read %d bytes: %s\n", n, string(buf[:2]))
		fmt.Println("Playing audio........")
	}

	// Step 7: Close the connection and ALSA resources
	fmt.Println("Playback finished.")

}
