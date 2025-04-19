package main

import (
	"fmt"
	"os"
)

func main() {
	addr := os.Args[1]
	method := os.Args[2]

	switch method {
	case "client":
		handleClient(addr)
	case "server":
		handleServer(addr)
	default:
		fmt.Printf("Invalid method %s", method)
	}
}

func handleServer(addr string) {
	filename := os.Args[3]

	fmt.Printf("Sending file %s to addr: %s\n", filename, addr)
	err := Send(filename, addr); if err != nil {
		fmt.Printf("Cannot send file %s to addr: %s\n%v\n", filename, addr, err)
		return;
	}
}

func handleClient(addr string) {
	err := Receive(addr)
	if err != nil {
		fmt.Printf("Cannot receive file from addr: %s\n%v\n", addr, err)
		return
	}
}