package main

import "fmt"

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
	// TODO
}
