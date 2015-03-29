package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	tap "github.com/traetox/goTunTap"
)

var (
	br = flag.String("br", "br0", "bridge to create and drop taps into")
	tp = flag.String("tp", "tap0", "Tap to create for link")
	remote = flag.String("s", "10.0.0.1:9999", "String for server")
)

func init() {
	flag.Parse()
	if *br == "" || *tp == "" || *remote == "" {
		log.Fatal("Invalid parameters")
	}
}

func main() {
	conn, err := net.Dial("tcp", *remote)
	if err != nil {
		log.Fatal("Failed to dial server:", err)
	}
	defer conn.Close()
	if err := tap.CreateBridge(*br); err != nil {
		log.Fatal("Failed to create bridge", err)
	}
	t, err := tap.CreateTap(*tp)
	if err != nil {
		log.Fatal("Failed to create tap manager")
	}
	if err := t.AddToBridge(*br); err != nil {
		log.Fatal("Failed to add tap to bridge", err)
	}
	if err := t.Relay(conn); err != nil {
		log.Fatal("Failed to relay tap connection")
	}
	fmt.Printf("DONE\n")
}
