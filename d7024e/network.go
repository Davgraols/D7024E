package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strconv"

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
		rpc.SenderIp = addr.String()

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
func (network *Network) sendLookupKresp(target *KademliaID, kademlia *Kademlia, contact *Contact) { //[]Contact
	fmt.Printf("im in sendLookupKresp")
	Kcontact := kademlia.routingTB.FindClosestContacts(target, 20)

	var rpcklist []*RPCKnearest

	for i := 0; i < len(Kcontact); i++ {
		rpcnearest := RPCKnearest{
			Id: Kcontact[i].ID.ToBytes(),
			Ip: []byte(Kcontact[i].Address),
		}
		rpcklist = append(rpcklist, &rpcnearest)
	}

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

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

func CheckError(err error) {
	if err != nil {
		log.Fatal("Error: ", err)
	}
}
