package main

import (
	"flag"
	"fmt"
	"time"
)

var (
	// Global variables
	K        = 20
	Alpha    = 3
	RT       = RoutingTable{} // Needs mutex
	MyId     = NewRandomKademliaID()
	Requests = make(chan RPC, 5)
	Files    = make(map[string]string) // Needs mutex
	Net      = Network{Port: "4000", BootstrapIP: "127.0.0.1"}

	// Local Variables
	mode = flag.String("m", "server", "mode: client or server")
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
	kademlia := Kademlia{}

	if !bootstrap {
		me := NewContact(MyId, "kademliaNodes")
		RT = *NewRoutingTable(me)
		bootstrapId := NewKademliaID("77ff0a3a0ec73e10ff408ece8728f84ae1af7bbf")
		bootstrapNode := NewContact(bootstrapId, "kademliaBootstrap")
		RT.AddContact(bootstrapNode)
		time.Sleep(1 * time.Second)
		go Listen("127.0.0.1", 4000)
		//go Net.SendPingMessage(&bootstrapNode)
		go kademlia.LookupContact(me.ID)

	} else {
		me := NewContact(MyId, "kademliaBootstrap")
		RT = *NewRoutingTable(me)
		go Listen("127.0.0.1", 4000)
	}

	for {
		msg := <-Requests
		switch msg.RpcType {
		case 0:
			go handlePingReq(msg)
		case 1:
			go handlePingRes(msg)
		case 2:
			go handleStoreReq(msg)
		case 3:
			go handleStoreRes(msg)
		case 4:
			go handleFindNodeReq(msg)
		case 5:
			go handleFindNodeRes(msg)
		case 6:
			go handleFindValueReq(msg)
		case 7:
			go handleFindValueRes(msg)

		default:
			fmt.Println("default")
		}
	}
}

func handlePingReq(msg RPC) {
	fmt.Println("Received PING from: ", msg.SenderIp)
	contact := NewContact(IdFromBytes(msg.SenderId), msg.SenderIp)
	go Net.SendPingResponseMessage(&contact)
}

func handlePingRes(msg RPC) {
	fmt.Println("Received PONG from: ", msg.SenderIp)
	//senderIp := strings.Split(msg.SenderIp, ":")[0]
	//contact := NewContact(IdFromBytes(msg.SenderId), senderIp)
	//go network.SendPingMessage(&contact)
}

func handleStoreReq(msg RPC) {
	fmt.Println("Received STORE_REQ from: ", msg.SenderIp)
	//data := msg.Value
}

func handleStoreRes(msg RPC) {
	fmt.Println("Received STORE_RES from: ", msg.SenderIp)
}

func handleFindNodeReq(msg RPC) {
	//rpc för hitta k närmsta
	fmt.Println("Received FIND_NODE_REQ from: ", msg.SenderIp)
	contact := NewContact(IdFromBytes(msg.SenderId), msg.SenderIp)

	// TODO aquire RT mutex
	RT.AddContact(contact)
	fmt.Println("Added contact: ", contact.String())
	Net.sendLookupKresp(IdFromBytes(msg.LookupId), &contact)
}

func handleFindNodeRes(msg RPC) {
	//rpc svar för hita k närmsta
	fmt.Println("Received FIND_NODE_RES from: ", msg.SenderIp)
	klist := msg.Klist

	// TODO aquire RT mutex
	for i := 0; i < len(klist); i++ {
		id := klist[i].Id
		ip := klist[i].Ip
		newid := IdFromBytes(id)
		newnode := NewContact(newid, string(ip))
		RT.AddContact(newnode)
		fmt.Println("Added contact: ", newnode.String())
	}
}

func handleFindValueReq(msg RPC) {
	fmt.Println("Received FIND_VALUE_REQ from: ", msg.SenderIp)
}

func handleFindValueRes(msg RPC) {
	fmt.Println("Received FIND_VALUE_RES from: ", msg.SenderIp)
}
