package d7024e

import (
	"fmt"
)

var (
	MyId          = NewRandomKademliaID()
	me            = NewContact(MyId, "127.0.0.1")
	routingTable  = NewRoutingTable(me)
	bootstrapId   = NewKademliaID("77ff0a3a0ec73e10ff408ece8728f84ae1af7bbf")
	bootstrapNode = NewContact(bootstrapId, "127.0.0.1")
	network       = Network{Port: "4000", BootstrapIP: "127.0.0.1"}
	Requests      = make(chan RPC, 5)
)

func Init(bootstrap bool) {

	if !bootstrap {
		routingTable.AddContact(bootstrapNode)
		go network.SendPingMessage(&bootstrapNode)
		go Listen("127.0.0.1", 4001)
	} else {
		go Listen("127.0.0.1", 4000)
	}

	for {
		msg := <-Requests
		switch msg.RpcType {
		case 1:
			fmt.Println("PING", msg.SenderIp)
			contact := NewContact(IdFromBytes(msg.SenderId), string(msg.SenderIp))
			go network.SendPingResponseMessage(&contact)
		case 2:
			fmt.Println("PING_PONG", msg.SenderIp)
		default:
			fmt.Println("default")
		}
	}

}
