package d7024e

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
)

// Network ADSASD
type Network struct {
	Port        string
	BootstrapIP string
}

const messageSize int = 1024

func Listen(ip string, port int) {
	fmt.Println(strconv.Itoa(port))
	serverAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	CheckError(err)

	serverConn, err := net.ListenUDP("udp", serverAddr)
	CheckError(err)
	defer serverConn.Close()

	buf := make([]byte, messageSize)

	log.Println("Listening on port " + strconv.Itoa(port))
	for {
		n, addr, err := serverConn.ReadFromUDP(buf)
		rpc := &RPC{}
		err = proto.Unmarshal(buf[0:n], rpc)

		rpc.SenderIp = addr.String()
		//fmt.Println(rpc.SenderIp)
		Requests <- *rpc
		/*
			switch rpc.RpcType {
			case 1:
				//log.Printf("Received PING from %s\n", hex.EncodeToString(rpc.SenderId))
				Requests <- "PING"
				/*contact := NewContact(IdFromBytes(rpc.SenderId), addr.String())
				fmt.Println(contact.String())
				// add contact to routing table

			default:
				Requests <- "OTHER"
				log.Printf("Received something else from addr: %s", addr)
			}*/

		CheckError(err)
	}
}

func (network *Network) TestPing() {
	contact := NewContact(NewRandomKademliaID(), "127.0.0.1")
	network.SendPingMessage(&contact)
}

// SendPingMessage ASDASD
func (network *Network) SendPingMessage(contact *Contact) {
	remoteAddr, err := net.ResolveUDPAddr("udp", network.BootstrapIP+":"+network.Port)
	CheckError(err)

	localAddr, err := net.ResolveUDPAddr("udp", network.BootstrapIP+":0")
	CheckError(err)

	conn, err := net.DialUDP("udp", localAddr, remoteAddr)
	CheckError(err)

	defer conn.Close()

	i := 1
	for {
		rpc := CreatePingRPC()
		data, err := proto.Marshal(rpc)
		if err != nil {
			log.Fatal("marshalling error: ", err)
		}
		buf := []byte(data)
		_, err = conn.Write(buf)
		fmt.Printf("sending PING with id %s\n", hex.EncodeToString(rpc.SenderId))
		if err != nil {
			log.Println(err)
		}

		i++
		time.Sleep(time.Second * 1)
	}
}

func (network *Network) SendPingResponseMessage(contact *Contact) {

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

func CreatePingRPC() *RPC {
	ser := []byte{112, 64, 128}
	sender_id := NewRandomKademliaID().ToBytes()
	rpc := RPC{
		RpcType:  1,
		Ser:      ser,
		SenderId: sender_id,
	}
	return &rpc
}
