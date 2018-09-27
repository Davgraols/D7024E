package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	kad "D7024E-GR8/d7024e"

	"github.com/golang/protobuf/proto"
)

var (
	mode = flag.String("m", "server", "mode: client or server")
	//port = flag.String("p", "4000", "host: ip:port")
	//port    = "4000"
	network = kad.Network{Port: "4000", BootstrapIP: "127.0.0.1"}
)

// run server with: go run main.go RPC.pb.go
// run client with: go run main.go RPC.pb.go-p 4000 -m client
func main() {

	//flag.Parse()
	//switch *mode {
	//case "server":
	//	RunServer()
	//case "client":
	//	RunClient()
	//}
	kad.NewRoutingTable()
	//id := kad.NewRandomKademliaID()
	//println(id.String(), hex.EncodeToString(id.ToBytes()))

}

func RunServer() {
	serverAddr, err := net.ResolveUDPAddr("udp", ":"+network.Port)
	CheckError(err)

	serverConn, err := net.ListenUDP("udp", serverAddr)
	CheckError(err)
	defer serverConn.Close()

	buf := make([]byte, 1024)

	log.Println("Listening on port " + network.Port)
	for {
		n, addr, err := serverConn.ReadFromUDP(buf)
		rpc := &RPC{}
		err = proto.Unmarshal(buf[0:n], rpc)

		switch rpc.RpcType {
		case 1:
			log.Printf("Received PING from %s\n", hex.EncodeToString(rpc.SenderId))
		default:
			log.Printf("Received something else from addr: %s", addr)
		}

		if err != nil {
			log.Fatal("Error: ", err)
		}
	}
}

func RunClient() {
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

func CreatePingRPC() *RPC {
	ser := []byte{112, 64, 128}
	sender_id := kad.NewRandomKademliaID().ToBytes()
	rpc := RPC{
		RpcType:  1,
		Ser:      ser,
		SenderId: sender_id,
	}
	return &rpc
}

func CheckError(err error) {
	if err != nil {
		log.Fatal("Error: ", err)
	}
}
