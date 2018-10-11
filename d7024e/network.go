package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
)

// Network ADSASD
type Network struct {
	Port        string
	BootstrapIP string
}

const messageSize int = 1024

func Listen(ip string, port int) {

	pc, err := net.ListenPacket("udp", ":"+strconv.Itoa(port))
	CheckError(err)
	defer pc.Close()

	buf := make([]byte, messageSize)

	for {
		n, addr, err := pc.ReadFrom(buf)
		rpc := &RPC{}
		err = proto.Unmarshal(buf[0:n], rpc)
		rpc.SenderIp = strings.Split(addr.String(), ":")[0]

		Requests <- *rpc
		CheckError(err)
	}
}

// SendPingMessage ASDASD
func (network *Network) SendPingMessage(contact *Contact) {

	rpc := RPC{
		RpcType:  0,
		Ser:      1337,
		SenderId: MyId.ToBytes(),
	}

	data, err := proto.Marshal(&rpc)
	if err != nil {
		log.Fatal("marshalling error: ", err)
	}

	buf := []byte(data)

	conn, err := net.Dial("udp", contact.Address+":4000")
	CheckError(err)
	defer conn.Close()

	conn.Write(buf)
	fmt.Printf("sending PING with id %s to %s", hex.EncodeToString(rpc.SenderId), contact.Address)

	if err != nil {
		log.Println(err)
	}

}

func (network *Network) SendPingResponseMessage(contact *Contact) {

	rpc := RPC{
		RpcType:  1,
		Ser:      1337,
		SenderId: MyId.ToBytes(),
	}

	data, err := proto.Marshal(&rpc)
	if err != nil {
		log.Fatal("marshalling error: ", err)
	}
	buf := []byte(data)

	conn, err := net.Dial("udp", contact.Address+":4000")
	CheckError(err)
	defer conn.Close()

	conn.Write(buf)
	fmt.Printf("sending PONG with id %s to %s", hex.EncodeToString(rpc.SenderId), contact.Address)

}

func (network *Network) sendLookupKmessage(Kcontact Contact, target *KademliaID) {
	rpc := RPC{
		RpcType:  4,
		Ser:      1337,
		SenderId: MyId.ToBytes(),
		LookupId: target.ToBytes(),
	}
	data, err := proto.Marshal(&rpc)
	if err != nil {
		log.Fatal("marshalling error: ", err)
	}
	buf := []byte(data)

	conn, err := net.Dial("udp", Kcontact.Address+":4000")
	CheckError(err)
	defer conn.Close()

	conn.Write(buf)
}

func contactListToRpc(contactList []Contact) []*RPCKnearest {
	var rpcklist []*RPCKnearest

	for i := 0; i < len(contactList); i++ {
		rpcnearest := RPCKnearest{
			Id: contactList[i].ID.ToBytes(),
			Ip: []byte(contactList[i].Address),
		}
		rpcklist = append(rpcklist, &rpcnearest)
	}

	return rpcklist
}

func (network *Network) sendLookupKresp(target *KademliaID, contact *Contact) {
	// TODO aquire RT mutex
	Kcontact := RT.FindClosestContacts(target, K)

	fmt.Printf("In sendLookupKresp. Found %d contacts in RT", len(Kcontact))
	//var rpcklist []*RPCKnearest
	rpcklist := contactListToRpc(Kcontact)
	/*
		for i := 0; i < len(Kcontact); i++ {
			rpcnearest := RPCKnearest{
				Id: Kcontact[i].ID.ToBytes(),
				Ip: []byte(Kcontact[i].Address),
			}
			rpcklist = append(rpcklist, &rpcnearest)
		}*/

	rpc := RPC{
		RpcType:  5,
		Ser:      1337,
		SenderId: MyId.ToBytes(),
		Klist:    rpcklist,
	}

	data, err := proto.Marshal(&rpc)
	if err != nil {
		log.Fatal("marshalling error: ", err)
	}
	buf := []byte(data)

	conn, err := net.Dial("udp", contact.Address+":4000")
	CheckError(err)
	defer conn.Close()

	conn.Write(buf)
}

func (contac *Contact) makeRPClist(Kcontact *Contact) {

}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte, contact *Contact) {

	SerialLock.Lock()
	serial := Serial
	Serial++
	SerialLock.Unlock()

	rpc := RPC{
		RpcType:  2,
		Ser:      serial,
		SenderId: MyId.ToBytes(),
		Value:    data,
	}

	rpcData, err := proto.Marshal(&rpc)
	if err != nil {
		log.Fatal("marshalling error: ", err)
	}
	buf := []byte(rpcData)

	conn, err := net.Dial("udp", contact.Address+":4000")
	CheckError(err)
	defer conn.Close()
	conn.Write(buf)
	fmt.Printf("sending STORE_REQ with id %s to %s serial: %d", hex.EncodeToString(rpc.SenderId), contact.Address, serial)
}

func (network *Network) SendStoreResponseMessage(contact *Contact) {
	rpc := RPC{
		RpcType:  3,
		Ser:      1337,
		SenderId: MyId.ToBytes(),
	}

	data, err := proto.Marshal(&rpc)
	if err != nil {
		log.Fatal("marshalling error: ", err)
	}
	buf := []byte(data)

	conn, err := net.Dial("udp", contact.Address+":4000")
	CheckError(err)
	defer conn.Close()

	conn.Write(buf)
	fmt.Printf("sending STORE_RES with id %s to %s", hex.EncodeToString(rpc.SenderId), contact.Address)
}

func CheckError(err error) {
	if err != nil {
		log.Fatal("Error: ", err)
	}
}

func (network *Network) SendFindDataMessage(fileId *KademliaID, contact *Contact) {
	SerialLock.Lock()
	serial := Serial
	Serial++
	SerialLock.Unlock()

	rpc := RPC{
		RpcType:  6,
		Ser:      serial,
		SenderId: MyId.ToBytes(),
		LookupId: fileId.ToBytes(),
	}

	rpcData, err := proto.Marshal(&rpc)
	if err != nil {
		log.Fatal("marshalling error: ", err)
	}
	buf := []byte(rpcData)

	conn, err := net.Dial("udp", contact.Address+":4000")
	CheckError(err)
	defer conn.Close()
	conn.Write(buf)
	fmt.Printf("sending FIND_VALUE_REQ with id %s to %s serial: %d", hex.EncodeToString(rpc.SenderId), contact.Address, serial)
}

func (network *Network) SendFindDataResponseMessage(data []byte, contactList []Contact, contact *Contact, serial int32) {

	rpc := RPC{
		RpcType:  7,
		Ser:      serial,
		SenderId: MyId.ToBytes(),
	}

	if data != nil {
		rpc.Value = data
	} else if contactList != nil {
		rpc.Klist = contactListToRpc(contactList)
	} else {
		log.Fatal("No data or contacts received")
	}

	rpcData, err := proto.Marshal(&rpc)
	if err != nil {
		log.Fatal("marshalling error: ", err)
	}
	buf := []byte(rpcData)

	conn, err := net.Dial("udp", contact.Address+":4000")
	CheckError(err)
	defer conn.Close()
	conn.Write(buf)
	fmt.Printf("sending FIND_VALUE_RES with id %s to %s serial: %d", hex.EncodeToString(rpc.SenderId), contact.Address, serial)
}
