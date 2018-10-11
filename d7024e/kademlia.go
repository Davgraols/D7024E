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
	id       *KademliaID
}

func (kademlia *Kademlia) LookupContact(target *KademliaID) {
	fmt.Println("im in LookupContact")
	index := 0
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
			if !(RT.getBucketIndex(kContact[index].ID)) {
				targetcontact := NewContact(kContact[index].ID, kContact[index].Address) //adress = IP
				go Net.sendLookupKmessage(targetcontact, kContact[index].ID)
				Connections[Serial] = alphachannel1
				index = index + 1
			}
			if !(RT.getBucketIndex(kContact[index].ID)) {
				targetcontact = NewContact(kContact[index].ID, kContact[index].Address) //adress = IP
				go Net.sendLookupKmessage(targetcontact, kContact[index].ID)
				Connections[Serial] = alphachannel1
				index = index + 1
			}
			if !(RT.getBucketIndex(kContact[index].ID)) {
				targetcontact = NewContact(kContact[index].ID, kContact[index].Address) //adress = IP
				go Net.sendLookupKmessage(targetcontact, kContact[index].ID)
				Connections[Serial] = alphachannel1
				index = index + 1
				rounds = rounds + 1
			}

		}
		respond := 0

		for respond < 3 {

			select {
			case msg1 = <-alphachannel1:
				tempK := msg1.klist
				concan.Append(makeKlist(tempK))
				concan.Sort()
				for l := 0; l < len(tempK); l++ {
					RT.AddContact(tempK[l])
				}
				respond = respond + 1

			case msg2 = <-alphachannel2:
				tempK2 := msg2.klist
				concan.Append(makeKlist(tempK2))
				concan.Sort()
				for l := 0; l < len(tempK2); l++ {
					RT.AddContact(tempK[l])
				}
				respond = respond + 1

			case msg3 = <-alphachannel3:
				tempK3 := msg3.klist
				concan.Append(makeKlist(tempK3))
				concan.Sort()
				for l := 0; l < len(tempK3); l++ {
					RT.AddContact(tempK[l])
				}
				respond = respond + 1
			}
		}
		newKlist = concan.GetContacts(20)
		if (kContact[0] == newKlist[0]) && (kContact[1] == newKlist[1]) && (kContact[2] == newKlist[2]) {
			return newKlist
		}
	} //end of round
	return newKlist
}

func makenodeobj(contact *Contact, sernr int) Nodeobj {
	node := Nodeobj{
		ser:      sernr,
		lookedup: false,
		Contact:  contact,
		id:       contact.ID,
	}

	return node
}

func makeKlist(klist []*RPCKnearest) []Contact {
	var newKlist []Contact

	// TODO aquire RT mutex
	for i := 0; i < len(klist); i++ {
		id := klist[i].Id
		ip := klist[i].Ip
		newid := IdFromBytes(id)
		newnode := NewContact(newid, string(ip))
		newKlist = append(newKlist, newnode)
		fmt.Println("Added contact: ", newnode.String())
	}

	return newKlist

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
