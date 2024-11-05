package main

import (
	"log"
	"math/rand"
	"time"

	"gitlab.com/patopest/go-sacn"
	"gitlab.com/patopest/go-sacn/packet"
)

func main() {
	log.Println("Hello, I'll send some random data to the universe 1")

	sender, err := sacn.NewSender("192.168.26.115", &sacn.SenderOptions{}) // Create sender
	if err != nil {
		log.Fatal(err)
	}

	// Initialise universe
	var uni uint16 = 1
	universe, err := sender.StartUniverse(uni)
	if err != nil {
		log.Fatal(err)
	}
	sender.SetMulticast(uni, true)

	// Create a random seed
	src := rand.NewSource(time.Now().UnixNano()) // Create a new source
	r := rand.New(src)                           // Create a new random generator with the source

	// Create new packet and fill it up with data
	p := packet.NewDataPacket()
	p.SetData([]uint8{
		// Fixture 1
		byte(r.Intn(256)), // Red
		byte(r.Intn(256)), // Green
		byte(r.Intn(256)), // Blue
		// Fixture 2
		byte(r.Intn(256)), // Red
		byte(r.Intn(256)), // Green
		byte(r.Intn(256)), // Blue
	})
	log.Println("Sending packet")

	for i := 0; i < 10; i++ {
		universe <- p // send the packet

		time.Sleep(1 * time.Second)

		p.SetData([]uint8{
			// Fixture 1
			byte(r.Intn(256)), // Red
			byte(r.Intn(256)), // Green
			byte(r.Intn(256)), // Blue
			// Fixture 2
			byte(r.Intn(256)), // Red
			byte(r.Intn(256)), // Green
			byte(r.Intn(256)), // Blue
		})
	}

	// To stop the universe and advertise termination to receivers
	close(universe)

	time.Sleep(1 * time.Second)

	// To close the sender altogether
	sender.Close()
}
