package main

import "fmt"

type Kademlia struct {
	routingTB RoutingTable
	//file      []byte
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	fmt.Println("im in LookupContact")
	Kcontact := kademlia.routingTB.FindClosestContacts(target.ID, 20)
	fmt.Println(Kcontact)
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
