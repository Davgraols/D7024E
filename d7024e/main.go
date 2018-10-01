package main

import (
	"flag"
	"fmt"
)

var (
	mode     = flag.String("m", "server", "mode: client or server")
	MyId     = NewRandomKademliaID()
	network  = Network{Port: "4000", BootstrapIP: "127.0.0.1"}
	Requests = make(chan RPC, 5)
)

func main() {

	flag.Parse()
	switch *mode {
	case "server":
		run(true)
	case "client":
		run(false)
	}
}

func run(bootstrap bool) {

	if !bootstrap {
		me := NewContact(MyId, "kademliaNodes")
		routingTable := NewRoutingTable(me)
		bootstrapId := NewKademliaID("77ff0a3a0ec73e10ff408ece8728f84ae1af7bbf")
		bootstrapNode := NewContact(bootstrapId, "kademliaBootstrap")
		routingTable.AddContact(bootstrapNode)
		go network.SendPingMessage(&bootstrapNode)
		go Listen("127.0.0.1", 4000)
	} else {
		me := NewContact(MyId, "kademliaBootstrap")
		routingTable := NewRoutingTable(me)
		fmt.Println(routingTable.me.String())
		go Listen("127.0.0.1", 4000)
	}

	for {
		msg := <-Requests
		switch msg.RpcType {
		case 1:
			fmt.Println("PING", msg.SenderIp)
			contact := NewContact(IdFromBytes(msg.SenderId), msg.SenderIp)
			go network.SendPingResponseMessage(&contact)
		case 2:
			fmt.Println("PONG", msg.SenderIp)
		default:
			fmt.Println("default")
		}
	}
}
