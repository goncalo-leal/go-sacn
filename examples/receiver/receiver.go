package main

import (
	"fmt"
	"net"
	"time"

	"gitlab.com/patopest/go-sacn"
	"gitlab.com/patopest/go-sacn/packet"
)

func main() {
	fmt.Println("hello")

	itf, _ := net.InterfaceByName("enxd0c0bf4d488a") // specific to your machine
	receiver := sacn.NewReceiver(itf)

	receiver.JoinUniverse(sacn.DISCOVERY_UNIVERSE)
	receiver.RegisterPacketCallback(packet.PacketTypeDiscovery, discoveryPacketCallback)

	receiver.JoinUniverse(1)
	receiver.RegisterPacketCallback(packet.PacketTypeData, dataPacketCallback)
	receiver.RegisterTerminationCallback(universeTerminatedCallback)

	receiver.Start()

	for {
		time.Sleep(1 * time.Nanosecond)
	}
}

func discoveryPacketCallback(p packet.SACNPacket, source string) {
	fmt.Printf("at least it call the function\n")

	d, ok := p.(*packet.DiscoveryPacket)
	if !ok {
		return
	}

	fmt.Printf("Discovered universes from %s:\n", string(d.SourceName[:]))
	for i := 0; i < d.GetNumUniverses(); i++ {
		fmt.Printf("%d, ", d.Universes[i])
	}
}

func dataPacketCallback(p packet.SACNPacket, source string) {
	d, ok := p.(*packet.DataPacket)
	if !ok {
		return
	}
	fmt.Printf("Received Data Packet for universe %d from %s\n", d.Universe, source)

	data, err := d.GetData()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Data: %v\n", data)

	data, err = d.GetData(1, 5)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Data (limited): %v\n", data)
}

func universeTerminatedCallback(universe uint16) {
	fmt.Printf("Universe %d is terminated", universe)
}
