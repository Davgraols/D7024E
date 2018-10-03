package main

import "fmt"

type Kademlia struct {
	routingTB RoutingTable
	//file      []byte
}

func (kademlia *Kademlia) LookupContact(target *Contact, net *Network) {
	fmt.Println("im in LookupContact")
	Kcontact := kademlia.routingTB.FindClosestContacts(target.ID, 20)
	//fmt.Println(Kcontact)
	for i := 0; i < len(Kcontact); i++ {
		fmt.Println("im in forloop")
		go target.sendLookupKmessage(Kcontact[i])
	}
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
