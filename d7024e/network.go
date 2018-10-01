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

		Requests <- *rpc
		CheckError(err)
	}
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
	_, err = conn.Write(buf)
	fmt.Printf("sending PING with id %s\n", hex.EncodeToString(rpc.SenderId))
	if err != nil {
		log.Println(err)
	}

}

func (network *Network) SendPingResponseMessage(contact *Contact) {
	remoteAddr, err := net.ResolveUDPAddr("udp", network.BootstrapIP+":"+network.Port)
	CheckError(err)

	localAddr, err := net.ResolveUDPAddr("udp", network.BootstrapIP+":0")
	CheckError(err)

	conn, err := net.DialUDP("udp", localAddr, remoteAddr)
	CheckError(err)

	defer conn.Close()

	rpc := RPC{
		RpcType:  2,
		Ser:      1337,
		SenderId: MyId.ToBytes(),
	}

	data, err := proto.Marshal(&rpc)
	if err != nil {
		log.Fatal("marshalling error: ", err)
	}
	buf := []byte(data)
	_, err = conn.Write(buf)
	fmt.Printf("sending PING with id %s\n", hex.EncodeToString(rpc.SenderId))
	if err != nil {
		log.Println(err)
	}

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
