package main

import (
	"flag"

	kad "github.com/davgraols/D7024E-GR8/d7024e"
)

var (
	mode     = flag.String("m", "server", "mode: client or server")
	network  = kad.Network{Port: "4000", BootstrapIP: "127.0.0.1"}
	Requests = make(chan string)
)

// run server with: go run main.go RPC.pb.go
// run client with: go run main.go RPC.pb.go -m client
func main() {

	flag.Parse()
	switch *mode {
	case "server":
		go kad.Listen("127.0.0.1", 4000)
	case "client":
		go network.TestPing()
		go network.TestPing()
		go kad.Listen("127.0.0.1", 4001)
	}

	kad.Init()

}
