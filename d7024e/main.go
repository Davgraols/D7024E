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
	Files    = make(map[string]string)
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
		kademlia := Kademlia{
			routingTB: *routingTable,
		}

		//go network.SendPingMessage(&bootstrapNode)
		go kademlia.LookupContact(&me)
		go Listen("127.0.0.1", 4000)
	} else {
		go Listen("127.0.0.1", 4000)
	}

	for {
		me := NewContact(MyId, "kademliaNodes")
		routingTable := NewRoutingTable(me)
		bootstrapId := NewKademliaID("77ff0a3a0ec73e10ff408ece8728f84ae1af7bbf")
		bootstrapNode := NewContact(bootstrapId, "kademliaBootstrap")
		routingTable.AddContact(bootstrapNode)
		kademlia := Kademlia{
			routingTB: *routingTable,
		}
		msg := <-Requests
		switch msg.RpcType {
		case 0:
			fmt.Println("Received PING from: ", msg.SenderIp)
			//senderIp := strings.Split(msg.SenderIp, ":")[0]
			//contact := NewContact(IdFromBytes(msg.SenderId), senderIp)
			//go network.SendPingResponseMessage(&contact)
			//go Kademlia.LookupContact(&contact)
		case 1:
			fmt.Println("Received PONG from: ", msg.SenderIp)
			//senderIp := strings.Split(msg.SenderIp, ":")[0]
			//contact := NewContact(IdFromBytes(msg.SenderId), senderIp)
			//go network.SendPingMessage(&contact)

		case 2:
			fmt.Println("Received STORE_REQ from: ", msg.SenderIp)
			//data := msg.Value

		case 3:
			fmt.Println("Received STORE_RES from: ", msg.SenderIp)

		case 4:
			//rpc för hitta k närmsta
			fmt.Println("Received FIND_NODE_REQ from: ", msg.SenderIp)
			senderIp := strings.Split(msg.SenderIp, ":")[0]
			contact := NewContact(IdFromBytes(msg.SenderId), senderIp)
			network.sendLookupKresp(MyId, &kademlia, &contact)
		case 5:
			//rpc svar för hittaa k närmsta
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
