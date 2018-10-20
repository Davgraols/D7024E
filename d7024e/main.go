package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

const (
	K              = 3
	Alpha          = 3
	NodeRepublish  = time.Second * 10
	OwnerRepublish = time.Second * 35
	TimeOut        = time.Second * 15
	MainDebug      = false
	NetworkDebug   = false
	KademliaDebug  = true
	FileStoreDebug = false
)

var (
	MyId = NewRandomKademliaID()

	// Global instances
	RT          = &RoutingTable{} // Needs mutex
	Connections = make(map[int32]chan RPC)
	FS          = NewFileStore()
	Filemap     = make(map[KademliaID][]byte)
	Requests    = make(chan RPC, 5)
	Net         = Network{Port: "4000", BootstrapIP: "127.0.0.1"}
	KademliaObj = Kademlia{}

	//RTLock Global Locks
	RTLock         = &sync.Mutex{}
	FileLock       = &sync.Mutex{}
	ConnectionLock = &sync.Mutex{}

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

	if !bootstrap {
		me := NewContact(MyId, "kademliaNodes")
		RT = NewRoutingTable(me)
		bootstrapId := NewKademliaID("77ff0a3a0ec73e10ff408ece8728f84ae1af7bbf")
		bootstrapNode := NewContact(bootstrapId, "kademliaBootstrap")
		RT.AddContact(bootstrapNode)
		go Listen("127.0.0.1", 4000)
		go nodeInit()

	} else {
		MyId = NewKademliaID("77ff0a3a0ec73e10ff408ece8728f84ae1af7bbf")
		me := NewContact(MyId, "kademliaBootstrap")
		RT = NewRoutingTable(me)
		go Listen("127.0.0.1", 4000)
		go bootstrapInit()
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

func nodeInit() {
	time.Sleep(1 * time.Second)
	go KademliaObj.LookupContact(MyId)
}

func bootstrapInit() {
	time.Sleep(NodeRepublish)
	KademliaObj.Store([]byte("hello asd"), &RT.me)
	time.Sleep(NodeRepublish * 2)
	KademliaObj.Unpin(NewRandomHash("hello asd"))
}

func handlePingReq(msg RPC) {
	serialnr := NewRandomSerial()
	if MainDebug {
		fmt.Println("Received PING from: ", msg.SenderIp)
	}
	contact := NewContact(IdFromBytes(msg.SenderId), msg.SenderIp)
	go Net.SendPingResponseMessage(&contact, serialnr)
}

func handlePingRes(msg RPC) {
	if MainDebug {
		fmt.Println("Received PONG from: ", msg.SenderIp)
	}
}

func handleStoreReq(msg RPC) {
	if MainDebug {
		fmt.Printf("Received STORE_REQ from: %s with serial: %d\n", msg.SenderIp, msg.Ser)
	}

	owner := NewContact(IdFromBytes(msg.LookupId), msg.OwnerIp)
	if !owner.ID.Equals(RT.me.ID) {
		fileID := NewRandomHash(string(msg.Value))
		FS.StoreFile(msg.Value, &owner)
		FS.SetRepublished(fileID, true)
		if MainDebug {
			fmt.Println("Stored file: ", string(msg.Value))
		}
	} else {
		if MainDebug {
			fmt.Println("Did not store file. Am owner")
		}
	}

	contact := NewContact(IdFromBytes(msg.SenderId), msg.SenderIp)

	RTLock.Lock()
	RT.AddContact(contact)
	RTLock.Unlock()

	Net.SendStoreResponseMessage(&contact, msg.Ser)
}

func handleStoreRes(msg RPC) {
	if MainDebug {
		fmt.Printf("Received STORE_RES from: %s with serial: %d\n", msg.SenderIp, msg.Ser)
	}
}

//RPC4
func handleFindNodeReq(msg RPC) {
	//rpc för hitta k närmsta
	if MainDebug {
		fmt.Printf("Received FIND_NODE_REQ from: %s with serial: %d\n", msg.SenderIp, msg.Ser)
	}
	contact := NewContact(IdFromBytes(msg.SenderId), msg.SenderIp)

	RTLock.Lock()
	RT.AddContact(contact)
	RTLock.Unlock()
	//fmt.Println("Added contact: ", contact.String())
	Net.sendLookupKresp(IdFromBytes(msg.LookupId), &contact, msg.Ser)
}

//RPC5
func handleFindNodeRes(msg RPC) {
	//rpc svar för hitta k närmsta
	if MainDebug {
		fmt.Printf("Received FIND_NODE_RES from: %s with serial: %d\n", msg.SenderIp, msg.Ser)
	}
	ConnectionLock.Lock()
	Connections[msg.Ser] <- msg
	ConnectionLock.Unlock()
}

func handleFindValueReq(msg RPC) {
	if MainDebug {
		fmt.Printf("Received FIND_VALUE_REQ from: %s with serial: %d\n", msg.SenderIp, msg.Ser)
	}
	fileId := IdFromBytes(msg.LookupId)
	file, exists := FS.getFile(fileId)
	if MainDebug {
		fmt.Printf("File id: %s exists: %t\n", fileId.String(), exists)
	}
	contact := NewContact(IdFromBytes(msg.SenderId), msg.SenderIp)

	RTLock.Lock()
	RT.AddContact(contact)
	RTLock.Unlock()

	if exists {
		if MainDebug {
			fmt.Println("File exists!")
		}
		owner := file.owner
		Net.SendFindDataResponseMessage(file.content, nil, &contact, msg.Ser, &owner)
	} else {
		if MainDebug {
			fmt.Println("File does not exist")
		}
		RTLock.Lock()
		contacts := RT.FindClosestContacts(IdFromBytes(msg.LookupId), 20)
		RTLock.Unlock()
		Net.SendFindDataResponseMessage(nil, contacts, &contact, msg.Ser, nil)
	}
}

func handleFindValueRes(msg RPC) {
	if MainDebug {
		fmt.Printf("Received FIND_VALUE_RES from: %s with serial: %d\n", msg.SenderIp, msg.Ser)
	}

	ConnectionLock.Lock()
	Connections[msg.Ser] <- msg
	ConnectionLock.Unlock()
}
