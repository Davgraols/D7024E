package main

import (
	"flag"

	kad "D7024E-GR8/d7024e"
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
		kad.Init(true)
	case "client":
		kad.Init(false)
	}
}
