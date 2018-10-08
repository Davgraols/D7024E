package main

import (
	"fmt"
	"time"
)

type Kademlia struct {
	//routingTB RoutingTable replaced with global variable
	//file      []byte
}

func (kademlia *Kademlia) LookupContact(target *KademliaID) {
	fmt.Println("im in LookupContact")

	// TODO aquire RT mutex
	kContact := RT.FindClosestContacts(target, K)
	//fmt.Println(Kcontact)

	// implement full lookup procedure so that only Alpha (global variable) calls are made in parallel.
	for i := 0; i < len(kContact); i++ {
		fmt.Println("im in forloop")
		go Net.sendLookupKmessage(kContact[i], target)
	}

}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	targetID := NewRandomHash(string(data))

	FileLock.Lock()
	Files[targetID.String()] = data
	FileLock.Unlock()

	RTLock.Lock()
	closetsContacts := RT.FindClosestContacts(targetID, K)
	RTLock.Unlock()
	for _, contact := range closetsContacts {
		Net.SendStoreMessage(data, &contact)
	}
}

func (kademlia *Kademlia) republish(fileHash string, after time.Duration) {
	time.Sleep(after * time.Second)
	FileLock.Lock()
	file := Files[fileHash]
	FileLock.Unlock()
	kademlia.Store(file)
	fmt.Println("Republished file: ", string(Files[fileHash]))
	kademlia.republish(fileHash, after)
}
