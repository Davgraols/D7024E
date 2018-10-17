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
func (network *Network) SendPingMessage(contact *Contact, serialnr int32) {

	rpc := RPC{
		RpcType:  0,
		Ser:      serialnr,
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
	fmt.Printf("sending PING with id %s to %s\n", hex.EncodeToString(rpc.SenderId), contact.Address)

	if err != nil {
		log.Println(err)
	}

}

func (network *Network) SendPingResponseMessage(contact *Contact, serialnr int32) {

	rpc := RPC{
		RpcType:  1,
		Ser:      serialnr,
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
	fmt.Printf("sending PONG with id %s to %s\n", hex.EncodeToString(rpc.SenderId), contact.Address)

}

func (network *Network) sendLookupKmessage(Kcontact Contact, target *KademliaID, serialnr int32) {
	rpc := RPC{
		RpcType:  4,
		Ser:      serialnr,
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
	fmt.Printf("sending FIND_NODE_REQ with id %s to %s serial: %d\n", hex.EncodeToString(rpc.SenderId), Kcontact.Address, rpc.Ser)
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

func (network *Network) sendLookupKresp(target *KademliaID, contact *Contact, serialnr int32) {
	// TODO aquire RT mutex
	RTLock.Lock()
	Kcontact := RT.FindClosestContacts(target, K)
	RTLock.Unlock()

	rpcklist := contactListToRpc(Kcontact)

	rpc := RPC{
		RpcType:  5,
		Ser:      serialnr,
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
	fmt.Printf("sending FIND_NODE_RES with id %s to %s serial: %d\n", hex.EncodeToString(rpc.SenderId), contact.Address, rpc.Ser)
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte, contact *Contact) {

	rpc := RPC{
		RpcType:  2,
		Ser:      NewRandomSerial(),
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
	fmt.Printf("sending STORE_REQ with id %s to %s serial: %d\n", hex.EncodeToString(rpc.SenderId), contact.Address, rpc.Ser)
}

func (network *Network) SendStoreResponseMessage(contact *Contact, serialnr int32) {
	rpc := RPC{
		RpcType:  3,
		Ser:      serialnr,
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
	fmt.Printf("sending STORE_RES with id %s to %s\n", hex.EncodeToString(rpc.SenderId), contact.Address)
}

func CheckError(err error) {
	if err != nil {
		log.Fatal("Error: ", err)
	}
}

func (network *Network) SendFindDataMessage(fileId *KademliaID, contact *Contact, serialnr int32) {

	rpc := RPC{
		RpcType:  6,
		Ser:      serialnr,
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
	fmt.Printf("sending FIND_VALUE_REQ with id %s to %s serial: %d\n", hex.EncodeToString(rpc.SenderId), contact.Address, rpc.Ser)
}

func (network *Network) SendFindDataResponseMessage(data []byte, contactList []Contact, contact *Contact, serialnr int32) {

	rpc := RPC{
		RpcType:  7,
		Ser:      serialnr,
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
	fmt.Printf("sending FIND_VALUE_RES with id %s to %s serial: %d\n", hex.EncodeToString(rpc.SenderId), contact.Address, rpc.Ser)
}
