package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"../tap"
)

var (
	br = flag.String("br", "br0", "bridge to create and drop taps into")
	tp = flag.String("tp", "tap0", "Tap to create for link")
	port = flag.Int("p", 9999, "listening port")
)

func init() {
	flag.Parse()
	if *br == "" || *tp == "" {
		log.Fatal("Invalid parameters")
	}
	if *port >= 0xffff || *port <= 1024 {
		log.Fatal("Invalid listen port")
	}
}

func main() {
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
	conn, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		log.Fatal("Failed to dial server:", err)
	}
	defer conn.Close()

	for {
		c, err := conn.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}
		if err := t.Relay(c); err != nil {
			log.Fatal("Failed to relay tap connection")
		}
		c.Close()
	}
	fmt.Printf("DONE\n")
}
