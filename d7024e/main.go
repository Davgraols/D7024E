package main

import (
	"flag"
	"fmt"
	"strings"
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
			fmt.Println("Received PING from: ", msg.SenderIp)
			senderIp := strings.Split(msg.SenderIp, ":")[0]
			contact := NewContact(IdFromBytes(msg.SenderId), senderIp)
			go network.SendPingResponseMessage(&contact)
		case 2:
			fmt.Println("Received PONG from: ", msg.SenderIp)
			//senderIp := strings.Split(msg.SenderIp, ":")[0]
			//contact := NewContact(IdFromBytes(msg.SenderId), senderIp)
			//go network.SendPingMessage(&contact)
		default:
			fmt.Println("default")
		}
	}
}
