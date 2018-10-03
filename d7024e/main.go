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
	Files    = make(map[string][]byte)
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
		//go network.SendPingMessage(&bootstrapNode)
		fmt.Println(len(routingTable.FindClosestContacts(bootstrapId, 1)))
		go network.SendStoreMessage([]byte("Hello World!"), routingTable)
		go Listen("127.0.0.1", 4000)

	} else {
		me := NewContact(MyId, "kademliaBootstrap")
		routingTable := NewRoutingTable(me)
		fmt.Println("me: ", routingTable.me)
		go Listen("127.0.0.1", 4000)
	}

	for {
		msg := <-Requests
		switch msg.RpcType {
		case 0:
			fmt.Println("Received PING from: ", msg.SenderIp)
			senderIp := strings.Split(msg.SenderIp, ":")[0]
			contact := NewContact(IdFromBytes(msg.SenderId), senderIp)
			go network.SendPingResponseMessage(&contact)

		case 1:
			fmt.Println("Received PONG from: ", msg.SenderIp)
			//senderIp := strings.Split(msg.SenderIp, ":")[0]
			//contact := NewContact(IdFromBytes(msg.SenderId), senderIp)
			//go network.SendPingMessage(&contact)

		case 2:
			fmt.Println("Received STORE_REQ from: ", msg.SenderIp)
			messageString := string(msg.Value)
			fileHashID := NewRandomHash(messageString)
			Files[fileHashID.String()] = msg.Value
			fmt.Println("Stored file: ", messageString)
			senderIp := strings.Split(msg.SenderIp, ":")[0]
			contact := NewContact(IdFromBytes(msg.SenderId), senderIp)
			go network.SendStoreResponseMessage(&contact)

		case 3:
			fmt.Println("Received STORE_RES from: ", msg.SenderIp)

		case 4:
			fmt.Println("Received FIND_NODE_REQ from: ", msg.SenderIp)

		case 5:
			fmt.Println("Received FIND_NODE_RES from: ", msg.SenderIp)

		case 6:
			fmt.Println("Received FIND_VALUE_REQ from: ", msg.SenderIp)

		case 7:
			fmt.Println("Received FIND_VALUE_RES from: ", msg.SenderIp)

		default:
			fmt.Println("default")
		}
	}
}
