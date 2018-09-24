package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
)

var (
	mode = flag.String("m", "server", "mode: client or server")
	port = flag.String("p", "4000", "host: ip:port")
)

// run server with: go run main.go RPC.pb.go
// run client with: go run main.go RPC.pb.go-p 4000 -m client
func main() {
	flag.Parse()

	switch *mode {
	case "server":
		RunServer()
	case "client":
		RunClient()
	}
}

func RunServer() {
	serverAddr, err := net.ResolveUDPAddr("udp", ":"+*port)
	CheckError(err)

	serverConn, err := net.ListenUDP("udp", serverAddr)
	CheckError(err)
	defer serverConn.Close()

	buf := make([]byte, 1024)

	log.Println("Listening on port " + *port)
	for {
		n, addr, err := serverConn.ReadFromUDP(buf)
		rpc := &RPC{}
		err = proto.Unmarshal(buf[0:n], rpc)
		serial := string(rpc.Ser)
		id := string(rpc.SenderId)
		log.Printf("Received %s from %s address %s", serial, id, addr)

		if err != nil {
			log.Fatal("Error: ", err)
		}
	}
}

func RunClient() {
	remoteAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:"+*port)
	CheckError(err)

	localAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	CheckError(err)

	conn, err := net.DialUDP("udp", localAddr, remoteAddr)
	CheckError(err)

	defer conn.Close()

	i := 1
	for {
		rpc := CreateRPC()
		data, err := proto.Marshal(rpc)
		if err != nil {
			log.Fatal("marshalling error: ", err)
		}
		buf := []byte(data)
		_, err = conn.Write(buf)
		if err != nil {
			log.Println(err)
		}

		i++
		time.Sleep(time.Second * 1)
	}
}

func CreateRPC() *RPC {
	ser := []byte{112}
	sender_id := []byte{2, 3, 1, 3, 3, 16}
	rpc := RPC{
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
