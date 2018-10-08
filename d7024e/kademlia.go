package main

import "fmt"

type Kademlia struct {
	//routingTB RoutingTable replaced with global variable
	//file      []byte
}

//struct to nodeobject for the lookup algoritm
type Nodeobj struct {
	ser      int
	lookedup bool //false = not checked
	Contact  *Contact
	distance *KademliaID
	id       *KademliaID
}

func (kademlia *Kademlia) LookupContact(target *KademliaID) {
	fmt.Println("im in LookupContact")
	index := 0
	ser := 0
	rounds := 0
	alphachannel1 := make(chan RPC)
	alphachannel2 := make(chan RPC)
	alphachannel3 := make(chan RPC)

	kContact := RT.FindClosestContacts(target, K) // TODO aquire RT mutex

	for i := 0; i < len(kContact); i++ { // TODO make chanels an mutex
		if rounds == 3 {
			break
		} else {
			targetcontact := NewContact(kContact[index].ID, kContact[index].Address) //adress = IP
			go Net.sendLookupKmessage(targetcontact, kContact[index].ID)
			index = index + 1
			ser = ser + 1
			targetcontact = NewContact(kContact[index].ID, kContact[index].Address) //adress = IP
			go Net.sendLookupKmessage(targetcontact, kContact[index].ID)
			index = index + 1
			ser = ser + 1
			targetcontact = NewContact(kContact[index].ID, kContact[index].Address) //adress = IP
			go Net.sendLookupKmessage(targetcontact, kContact[index].ID)
			index = index + 1
			ser = ser + 1
			rounds = rounds + 1
		}
		respond := 0
		for respond < 3 {

			for i := range alphachannel1 {

				respond = respond + 1
			}
			for i := range alphachannel2 {

				respond = respond + 1
			}
			for i := range alphachannel3 {

				respond = respond + 1
			}
		}
	}

}
func compareLen(klist, ktarget []Contact) []Contact {

}

func makenodeobj(contact *Contact, sernr int) Nodeobj {
	node := Nodeobj{
		ser:      sernr,
		lookedup: false,
		Contact:  contact,
		distance: contact.ID.CalcDistance(MyId),
		id:       contact.ID,
	}

	return node
}

func findKClosest(alphalist Nodeobj) {

}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
