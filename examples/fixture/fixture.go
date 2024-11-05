package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"gitlab.com/patopest/go-sacn"
	"gitlab.com/patopest/go-sacn/packet"
)

const (
	UNIVERSE   = 1 // the fixture should listen to this universe
	N_CHANNELS = 3 // the fixture should listen to this number of channels
)

var (
	receiver      *sacn.Receiver
	isFirstSquare = true
	hasTerminated = true
	channel       = uint8(1) // the fixture should listen to this channel
)

func universeTerminatedCallback(universe uint16) {
	if universe != UNIVERSE {
		return
	}

	fmt.Printf("Universe %d is terminated\n", universe)
	hasTerminated = true
}

// drawSquare takes 3 bytes (RGB) and prints a colored square in the same position
func drawSquare(r, g, b byte, isFirst *bool) {
	color := fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b) // Set background color
	reset := "\033[0m"                                   // Reset color

	squareLength := 5 // Length of the square

	// Clear the screen if this is not the first square
	if !*isFirst {
		// Move the cursor up by L rows
		fmt.Printf("\033[%dA", squareLength)

		fmt.Printf("\r")
	}

	// Draw a LxL square in the specified position
	for i := 0; i < squareLength; i++ {
		for j := 0; j < squareLength; j++ {
			fmt.Print(color + "  " + reset) // Each "  " simulates a block
		}
		fmt.Println()
	}

	// change isFirst to false
	*isFirst = false
}

func dataPacketCallback(p packet.SACNPacket, source string) {
	d, ok := p.(*packet.DataPacket)
	if !ok {
		return
	}

	data, err := d.GetData(channel, N_CHANNELS)
	if err != nil {
		fmt.Println(err)
		return
	}

	drawSquare(data[0], data[1], data[2], &isFirstSquare)
}

func discoveryPacketCallback(p packet.SACNPacket, source string) {
	d, ok := p.(*packet.DiscoveryPacket)
	if !ok {
		return
	}

	// Check if the universe we are interested in is in the list of discovered universes
	var i int
	for i = 0; i < d.GetNumUniverses(); i++ {
		if d.Universes[i] == UNIVERSE {
			break
		}
	}
	if i == d.GetNumUniverses() {
		return
	}

	// If the universe is available and we don't already have established connection, join it
	if hasTerminated {
		fmt.Printf("Universe %d is available\n", UNIVERSE)

		// Join the universe
		receiver.JoinUniverse(UNIVERSE)
		receiver.RegisterPacketCallback(packet.PacketTypeData, dataPacketCallback)
		receiver.RegisterTerminationCallback(universeTerminatedCallback)
	}

	// Change hasTerminated to false
	hasTerminated = false
}

func main() {

	if len(os.Args) > 2 {
		fmt.Println("Usage: go run fixture.go [channel] (optional, default channel is 1)")
		os.Exit(1)
	}

	if len(os.Args) == 2 {
		channel = uint8(os.Args[1][0] - '0')
	}

	fmt.Printf("Hello, I'm waiting to receive sACN discovery packets. If active, I'll be listening to universe %d in channel %d\n", UNIVERSE, channel)

	// itf, _ := net.InterfaceByName("enxd0c0bf4d488a") // specific to your machine
	itf, _ := net.InterfaceByName("wlp2s0") // specific to your machine

	receiver = sacn.NewReceiver(itf)

	receiver.JoinUniverse(sacn.DISCOVERY_UNIVERSE)
	receiver.RegisterPacketCallback(packet.PacketTypeDiscovery, discoveryPacketCallback)
	receiver.RegisterTerminationCallback(universeTerminatedCallback)

	receiver.Start()

	for {
		time.Sleep(1 * time.Nanosecond)
	}
}
