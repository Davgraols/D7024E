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
	finalReqCount := 0
	serial := NewRandomSerial()
	alphachannel1 := make(chan RPC)
	alphachannel2 := make(chan RPC)
	alphachannel3 := make(chan RPC)
	lastchenel := make(chan RPC, K)
	var newKlist []Contact
	hasret := make(map[KademliaID]Contact)
	RTLock.Lock()
	hasret[*MyId] = RT.me
	kContact := RT.FindClosestContacts(target, Alpha) // TODO aquire RT mutex
	fmt.Println("Lookup current k closests: ", kContact)
	RTLock.Unlock()
	concan := ContactCandidates{
		contacts: newKlist,
	}
	concan.Append(kContact)

	for { // TODO make chanels an mutex
		kContact = newAlpha(hasret, kContact)
		currentcheck = 0
		fmt.Printf("Starting round: %d with %d contacts\n", rounds, len(kContact))
		if len(kContact) >= 1 {
			serial = NewRandomSerial()
			fmt.Println("Sending lookup 1")
			go Net.sendLookupKmessage(kContact[index], target, serial)
			ConnectionLock.Lock()
			Connections[serial] = alphachannel1
			ConnectionLock.Unlock()
			currentcheck++
			index++
		}

		if len(kContact) >= 2 {
			serial = NewRandomSerial()
			fmt.Println("Sending lookup 2")
			go Net.sendLookupKmessage(kContact[index], target, serial)
			ConnectionLock.Lock()
			Connections[serial] = alphachannel2
			ConnectionLock.Unlock()
			currentcheck++
			index++
		}

		if len(kContact) >= 3 {
			serial = NewRandomSerial()
			fmt.Println("Sending lookup 3")
			go Net.sendLookupKmessage(kContact[index], target, serial)
			ConnectionLock.Lock()
			Connections[serial] = alphachannel3
			ConnectionLock.Unlock()
			currentcheck++
			index++
		}

		respond := 0

		fmt.Printf("Wating for %d responses\n", currentcheck)
		for respond < currentcheck {

			select {
			case msg1 := <-alphachannel1:
				fmt.Println("Received response on aplhachannel1")
				tempK := makeKlist(msg1.Klist)
				concan.Append(tempK)
				fmt.Printf("Added %d contacts to concan. Current size: %d\n", len(tempK), concan.Len())
				tempContact := NewContact(IdFromBytes(msg1.SenderId), msg1.SenderIp)
				hasret[*tempContact.ID] = tempContact
				//hasret = append(hasret, NewContact(IdFromBytes(msg1.SenderId), msg1.SenderIp))
				respond++

			case msg2 := <-alphachannel2:
				fmt.Println("Received response on aplhachannel2")
				tempK2 := makeKlist(msg2.Klist)
				concan.Append(tempK2)
				fmt.Printf("Added %d contacts to concan. Current size: %d\n", len(tempK2), concan.Len())
				tempContact := NewContact(IdFromBytes(msg2.SenderId), msg2.SenderIp)
				hasret[*tempContact.ID] = tempContact
				//hasret = append(hasret, NewContact(IdFromBytes(msg2.SenderId), msg2.SenderIp))
				respond++

			case msg3 := <-alphachannel3:
				fmt.Println("Received response on aplhachannel3")
				tempK3 := makeKlist(msg3.Klist)
				concan.Append(tempK3)
				fmt.Printf("Added %d contacts to concan. Current size: %d\n", len(tempK3), concan.Len())
				tempContact := NewContact(IdFromBytes(msg3.SenderId), msg3.SenderIp)
				hasret[*tempContact.ID] = tempContact
				//hasret = append(hasret, NewContact(IdFromBytes(msg3.SenderId), msg3.SenderIp))
				respond++
			}
		}

		concan.calcDistances(target)
		concan.removeDuplicates()
		concan.Sort()
		newKlist = concan.GetContacts(K)

		index = 0
		rounds = rounds + 1
		if (rounds == 3) || (len(kContact) < 1) || (kContact[0] == newKlist[0]) {
			fmt.Println("In final requests")
			for _, contact := range newKlist {
				_, contacted := hasret[*contact.ID] // checks if contact has already been contacted
				if !contacted {                     // Only send to contacts that has not been contacted
					serial := NewRandomSerial()
					go Net.sendLookupKmessage(contact, target, serial)
					ConnectionLock.Lock()
					Connections[serial] = lastchenel
					ConnectionLock.Unlock()
					finalReqCount++
				}
			}
			fmt.Printf("Sent %d final requests\n", finalReqCount)
			break
		}

		kContact = newKlist
	}

	for res := 0; res < finalReqCount; res++ {
		msglast := <-lastchenel
		tempKlast := makeKlist(msglast.Klist)
		concan.Append(tempKlast)
	}
	fmt.Printf("Received %d final requests\n", finalReqCount)
	concan.calcDistances(target)
	concan.Sort()
	concan.RemoveContact(MyId)
	concan.removeDuplicates()
	newKlist = concan.GetContacts(K)
	fmt.Println("Lookup returned: ", newKlist)
	return newKlist
}

func newAlpha(checked map[KademliaID]Contact, klist []Contact) []Contact {
	var templist []Contact
	for _, contact := range klist {

		_, exist := checked[*contact.ID] // checks if contact has already been contacted
		if !exist {
			templist = append(templist, contact)
		}
	}
	return templist
}

func makeKlist(klist []*RPCKnearest) []Contact {
	var newKlist []Contact

	for i := 0; i < len(klist); i++ {
		id := klist[i].Id
		ip := klist[i].Ip
		newid := IdFromBytes(id)
		newnode := NewContact(newid, string(ip))
		newKlist = append(newKlist, newnode)
		//fmt.Println("Added contact: ", newnode.String())
	}

	return newKlist
}

func (kademlia *Kademlia) LookupData(target *KademliaID) []byte {
	fmt.Println("im in LookupData")
	index := 0
	rounds := 0
	currentcheck := 0
	finalReqCount := 0
	serial := NewRandomSerial()
	alphachannel1 := make(chan RPC)
	alphachannel2 := make(chan RPC)
	alphachannel3 := make(chan RPC)
	lastchenel := make(chan RPC, K)
	var newKlist []Contact
	hasret := make(map[KademliaID]Contact)
	RTLock.Lock()
	hasret[*MyId] = RT.me
	kContact := RT.FindClosestContacts(target, Alpha) // TODO aquire RT mutex
	RTLock.Unlock()
	concan := ContactCandidates{
		contacts: newKlist,
	}

	concan.Append(kContact)

	for { // TODO make chanels an mutex
		kContact = newAlpha(hasret, kContact)
		currentcheck = 0

		fmt.Printf("Starting round: %d with %d contacts\n", rounds, len(kContact))
		if len(kContact) >= 1 {
			serial = NewRandomSerial()
			fmt.Println("Sending lookup 1")
			go Net.SendFindDataMessage(target, &kContact[index], serial)
			ConnectionLock.Lock()
			Connections[serial] = alphachannel1
			ConnectionLock.Unlock()
			currentcheck++
			index++
		}

		if len(kContact) >= 2 {
			serial = NewRandomSerial()
			fmt.Println("Sending lookup 2")
			go Net.SendFindDataMessage(target, &kContact[index], serial)
			ConnectionLock.Lock()
			Connections[serial] = alphachannel2
			ConnectionLock.Unlock()
			currentcheck++
			index++
		}

		if len(kContact) >= 3 {
			serial = NewRandomSerial()
			fmt.Println("Sending lookup 3")
			go Net.SendFindDataMessage(target, &kContact[index], serial)
			ConnectionLock.Lock()
			Connections[serial] = alphachannel3
			ConnectionLock.Unlock()
			currentcheck++
			index++
		}

		respond := 0

		fmt.Printf("Wating for %d responses\n", currentcheck)
		for respond < currentcheck {

			select {
			case msg1 := <-alphachannel1:
				fmt.Println("LD Received response on aplhachannel1")
				if msg1.Value != nil {
					fmt.Println("LookupData returned value: ", string(msg1.Value))
					return msg1.Value
				} else {
					tempK := makeKlist(msg1.Klist)
					concan.Append(tempK)
					fmt.Printf("Added %d contacts to concan. Current size: %d\n", len(tempK), concan.Len())
				}
				tempContact := NewContact(IdFromBytes(msg1.SenderId), msg1.SenderIp)
				hasret[*tempContact.ID] = tempContact
				respond++

			case msg2 := <-alphachannel2:
				fmt.Println("LD Received response on aplhachannel2")
				if msg2.Value != nil {
					fmt.Println("LookupData returned value: ", string(msg2.Value))
					return msg2.Value
				} else {
					tempK := makeKlist(msg2.Klist)
					concan.Append(tempK)
					fmt.Printf("Added %d contacts to concan. Current size: %d\n", len(tempK), concan.Len())
				}
				tempContact := NewContact(IdFromBytes(msg2.SenderId), msg2.SenderIp)
				hasret[*tempContact.ID] = tempContact
				respond++

			case msg3 := <-alphachannel3:
				fmt.Println("LD Received response on aplhachannel3")
				if msg3.Value != nil {
					fmt.Println("LookupData returned value: ", string(msg3.Value))
					return msg3.Value
				} else {
					tempK := makeKlist(msg3.Klist)
					concan.Append(tempK)
					fmt.Printf("Added %d contacts to concan. Current size: %d\n", len(tempK), concan.Len())
				}
				tempContact := NewContact(IdFromBytes(msg3.SenderId), msg3.SenderIp)
				hasret[*tempContact.ID] = tempContact
				respond++
			}
		}

		concan.calcDistances(target)
		concan.Sort()
		concan.RemoveContact(MyId)
		concan.removeDuplicates()
		newKlist = concan.GetContacts(K)

		index = 0
		rounds = rounds + 1
		if (rounds == 3) || (len(kContact) < 1) || (kContact[0] == newKlist[0]) {
			fmt.Println("In final requests")
			for _, contact := range newKlist {
				_, contacted := hasret[*contact.ID] // checks if contact has already been contacted
				if !contacted {                     // Only send to contacts that has not been contacted
					serial := NewRandomSerial()
					go Net.SendFindDataMessage(target, &kContact[index], serial)
					ConnectionLock.Lock()
					Connections[serial] = lastchenel
					ConnectionLock.Unlock()
					finalReqCount++
				}
			}
			fmt.Printf("Sent %d final requests\n", finalReqCount)
			break
		}

		kContact = newKlist
	}

	for res := 0; res < finalReqCount; res++ {
		msglast := <-lastchenel
		fmt.Println()
		if msglast.Value != nil {
			fmt.Println("LookupData returned value: ", string(msglast.Value))
			return msglast.Value
		}
	}
	fmt.Printf("Received %d final requests\n", finalReqCount)
	fmt.Println("LookupData returned nil")
	return nil
}

// Store stores a "file" data
func (kademlia *Kademlia) Store(data []byte) {
	fmt.Println("Starting store procedure")
	targetID := NewRandomHash(string(data))

	FileLock.Lock()
	Files[*targetID] = data
	FileLock.Unlock()
	fmt.Printf("Stored file: %s with id: %s\n", string(data), targetID.String())

	go kademlia.republish(targetID, OwnerRepublish)

	closetsContacts := kademlia.LookupContact(targetID)

	for _, contact := range closetsContacts {
		Net.SendStoreMessage(data, &contact)
	}
}

func (kademlia *Kademlia) republish(fileHash *KademliaID, after time.Duration) {
	fmt.Println("Starting republish procedure")
	time.Sleep(after)
	FileLock.Lock()
	file := Files[*fileHash]
	FileLock.Unlock()
	kademlia.Store(file)
	fmt.Println("Republished file: ", string(Files[*fileHash]))
	kademlia.republish(fileHash, after)
}
