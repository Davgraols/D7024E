package main

import (
	"fmt"
	"time"
)

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
	var newKlist []Contact
	kContact := RT.FindClosestContacts(target, K) // TODO aquire RT mutex
	concan := ContactCandidates{
		contacts: newKlist,
	}

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
		msg1 := <-alphachannel1
		msg2 := <-alphachannel2
		msg3 := <-alphachannel3
		for respond < 3 {

			select{
				case msg1 = <- respond{
					concan.Append(msg1.klist)
					concan.Sort()
					respond = respond + 1
				}
				case msg3 = <- respond{
					concan.Append(msg1.klist)
					concan.Sort()
					respond = respond + 1
				}
				case msg2 = <- respond{
					concan.Append(msg1.klist)
					concan.Sort()
					respond = respond + 1
				}
			}
		}
	}

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
	targetID := NewRandomHash(string(data))

	FileLock.Lock()
	Files[*targetID] = data
	FileLock.Unlock()

	RTLock.Lock()
	closetsContacts := RT.FindClosestContacts(targetID, K)
	RTLock.Unlock()
	for _, contact := range closetsContacts {
		Net.SendStoreMessage(data, &contact)
	}
}

func (kademlia *Kademlia) republish(fileHash KademliaID, after time.Duration) {
	time.Sleep(after * time.Second)
	FileLock.Lock()
	file := Files[fileHash]
	FileLock.Unlock()
	kademlia.Store(file)
	fmt.Println("Republished file: ", string(Files[fileHash]))
	kademlia.republish(fileHash, after)
}
