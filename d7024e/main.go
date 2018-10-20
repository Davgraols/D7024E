package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

const (
	K              = 20
	Alpha          = 3
	NodeRepublish  = time.Second * 20
	OwnerRepublish = time.Second * 60
)

var (
	MyId = NewRandomKademliaID()

	// Global instances
	RT          = &RoutingTable{} // Needs mutex
	Connections = make(map[int32]chan RPC)
	Files       = make(map[KademliaID][]byte)
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

	Mainrest()
	/*
		flag.Parse()
		switch *mode {
		case "server":
			run(true)
		case "client":
			run(false)
		}*/
}

func run(bootstrap bool) {
	Mainrest()
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
	file := []byte("Hello world")
	go KademliaObj.Store(file)
	time.Sleep(10 * time.Second)
	//go kademlia.LookupContact(me.ID)
	fileId := NewRandomHash("Hello world")
	fmt.Println("Looking up file with id: ", fileId.String())
	go KademliaObj.LookupData(fileId)
}

func bootstrapInit() {
	//file := []byte("Hello world")
	//fileId := NewRandomHash(string(file))
	//Files[*fileId] = file
	//fmt.Printf("Stored file: %s with key: %s, to string: %s \n", string(file), fileId, fileId.String())
	//time.Sleep(time.Second * 5)
	//go kademlia.Store([]byte("Hello world"))
}

func handlePingReq(msg RPC) {
	serialnr := NewRandomSerial()
	fmt.Println("Received PING from: ", msg.SenderIp)
	contact := NewContact(IdFromBytes(msg.SenderId), msg.SenderIp)
	go Net.SendPingResponseMessage(&contact, serialnr)
}

func handlePingRes(msg RPC) {
	fmt.Println("Received PONG from: ", msg.SenderIp)
	//senderIp := strings.Split(msg.SenderIp, ":")[0]
	//contact := NewContact(IdFromBytes(msg.SenderId), senderIp)
	//go network.SendPingMessage(&contact)
}

func handleStoreReq(msg RPC) {
	fmt.Printf("Received STORE_REQ from: %s with serial: %d\n", msg.SenderIp, msg.Ser)

	fileHash := NewRandomHash(string(msg.Value))
	FileLock.Lock()
	Files[*fileHash] = msg.Value
	FileLock.Unlock()
	go KademliaObj.republish(fileHash, NodeRepublish)
	fmt.Println("Stored file: ", string(msg.Value))
	contact := NewContact(IdFromBytes(msg.SenderId), msg.SenderIp)

	RTLock.Lock()
	RT.AddContact(contact)
	RTLock.Unlock()

	Net.SendStoreResponseMessage(&contact, msg.Ser)
}

func handleStoreRes(msg RPC) {
	fmt.Printf("Received STORE_RES from: %s with serial: %d\n", msg.SenderIp, msg.Ser)
}

//RPC4
func handleFindNodeReq(msg RPC) {
	//rpc för hitta k närmsta
	fmt.Printf("Received FIND_NODE_REQ from: %s with serial: %d\n", msg.SenderIp, msg.Ser)
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
	fmt.Printf("Received FIND_NODE_RES from: %s with serial: %d\n", msg.SenderIp, msg.Ser)
	ConnectionLock.Lock()
	Connections[msg.Ser] <- msg
	ConnectionLock.Unlock()
}

func handleFindValueReq(msg RPC) {
	fmt.Printf("Received FIND_VALUE_REQ from: %s with serial: %d\n", msg.SenderIp, msg.Ser)
	fileId := IdFromBytes(msg.LookupId)
	FileLock.Lock()
	file, exists := Files[*fileId]
	FileLock.Unlock()
	fmt.Printf("File id: %s exists: %t\n", fileId.String(), exists)
	contact := NewContact(IdFromBytes(msg.SenderId), msg.SenderIp)

	RTLock.Lock()
	RT.AddContact(contact)
	RTLock.Unlock()

	if exists {
		fmt.Println("File exists!")
		Net.SendFindDataResponseMessage(file, nil, &contact, msg.Ser)
	} else {
		fmt.Println("File does not exist")
		RTLock.Lock()
		contacts := RT.FindClosestContacts(IdFromBytes(msg.LookupId), 20)
		RTLock.Unlock()
		Net.SendFindDataResponseMessage(nil, contacts, &contact, msg.Ser)
	}
}

func handleFindValueRes(msg RPC) {
	fmt.Printf("Received FIND_VALUE_RES from: %s with serial: %d\n", msg.SenderIp, msg.Ser)

	ConnectionLock.Lock()
	Connections[msg.Ser] <- msg
	ConnectionLock.Unlock()
}
