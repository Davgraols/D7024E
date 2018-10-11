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

func (kademlia *Kademlia) LookupContact(target *KademliaID) []Contact {
	fmt.Println("im in LookupContact")
	index := 0
	rounds := 0
	currentcheck := 0
	serial := NewRandomSerial()
	alphachannel1 := make(chan RPC)
	alphachannel2 := make(chan RPC)
	alphachannel3 := make(chan RPC)
	lastchenel := make(chan RPC, K)
	var newKlist []Contact
	var hasret []Contact
	kContact := RT.FindClosestContacts(target, Alpha) // TODO aquire RT mutex
	concan := ContactCandidates{
		contacts: newKlist,
	}

	for { // TODO make chanels an mutex
		kContact = newAlpha(hasret, kContact)

		if len(kContact) >= 1 {
			serial = NewRandomSerial()
			go Net.sendLookupKmessage(kContact[index], target, serial)
			ConnectionLock.Lock()
			Connections[serial] = alphachannel1
			ConnectionLock.Unlock()
			currentcheck++
			index = index + 1
		}

		if len(kContact) >= 2 {
			serial = NewRandomSerial()
			go Net.sendLookupKmessage(kContact[index], target, serial)
			ConnectionLock.Lock()
			Connections[serial] = alphachannel2
			ConnectionLock.Unlock()
			currentcheck++
			index = index + 1
		}

		if len(kContact) >= 3 {
			serial = NewRandomSerial()
			go Net.sendLookupKmessage(kContact[index], target, serial)
			ConnectionLock.Lock()
			Connections[serial] = alphachannel3
			ConnectionLock.Unlock()
			currentcheck++
			index = index + 1
		}

		respond := 0

		for respond < currentcheck {

			select {
			case msg1 := <-alphachannel1:
				tempK := makeKlist(msg1.Klist)
				concan.Append(tempK)
				hasret = append(hasret, NewContact(IdFromBytes(msg1.SenderId), msg1.SenderIp))
				respond = respond + 1

			case msg2 := <-alphachannel2:
				tempK2 := makeKlist(msg2.Klist)
				concan.Append(tempK2)
				hasret = append(hasret, NewContact(IdFromBytes(msg2.SenderId), msg2.SenderIp))
				respond = respond + 1

			case msg3 := <-alphachannel3:
				tempK3 := makeKlist(msg3.Klist)
				concan.Append(tempK3)
				hasret = append(hasret, NewContact(IdFromBytes(msg3.SenderId), msg3.SenderIp))
				respond = respond + 1

			}
		}

		currentcheck = 0
		concan.Sort()
		newKlist = concan.GetContacts(K)

		index = 0
		rounds = rounds + 1
		if (rounds == 3) || (kContact[0] == newKlist[0]) {
			for _, contact := range newKlist {
				serial := NewRandomSerial()
				go Net.sendLookupKmessage(contact, target, serial)
				ConnectionLock.Lock()
				Connections[serial] = lastchenel
				ConnectionLock.Unlock()
			}
			break
		}

		kContact = newKlist
	}

	for res := 0; res < K; res++ {
		msglast := <-lastchenel
		tempKlast := makeKlist(msglast.Klist)
		concan.Append(tempKlast)
	}
	concan.Sort()
	newKlist = concan.GetContacts(K)
	fmt.Println("Lookup returned: ", newKlist)
	return newKlist
}

func newAlpha(checked, klist []Contact) []Contact {
	var templist []Contact
	for _, contact := range klist {
		exist := false
		for _, checkedcon := range checked {
			if contact.ID.Equals(checkedcon.ID) {
				exist = true
			}
		}
		if !exist {
			templist = append(templist, contact)
		}
	}
	return templist
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
