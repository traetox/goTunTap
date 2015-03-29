package main

import (
	"fmt"
	"net"
	"math/rand"
	"time"
	"crypto/sha512"
)

var (
	rnd *rand.Rand
)

func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
}

type ClientHandler func(bridge string, sock net.Conn)

func ListenAndServe(ipPort, auth, bridge string, ch ClientHandler) error {
	listener, err := net.Listen("tcp", ipPort)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept: %s\n", err)
			continue
		}
		go handleNewClient(conn, bridge, auth, ch)
	}
}

func genSalt() []byte {
	salt := make([]byte, 16)
	var i int
	for i = 0; i < 15; i++ {
		salt[i] = byte((rnd.Uint32()%93)+0x21)
	}
	salt[15] = 0
	return salt
}

func handleNewClient(conn net.Conn, bridge, auth string, ch ClientHandler) {
	var buffer []byte

	salt := genSalt()
	/* do some shit to verify the client */
	//send salt
	conn.Write(salt)
	//receive response
	b, err := conn.Read(buffer)
	if err != nil && b != 64 {
		conn.Close()
		return
	}

	//sha512 hsh (salt + auth) and compare	
	if(!compareHash(genHash(salt, []byte(auth)), buffer)) {
		conn.Close()
		return
	}
	//inform client if its a go
	conn.Write([]byte("NinerNiner"))
	
	/* tell the client what their IP and subnet should be */
	
	
	/* launch the handler */
	//ch(bridge, conn)
	conn.Close()
}

func compareHash(a, b []byte) bool {
	if(len(a) != len(b)) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if(a[i] != b[i]) {
			return false
		}
	}
	return true
}

func genHash(salt, auth []byte) []byte {
	hasher := sha512.New()
	hasher.Write(salt)
	hasher.Write(auth)
	return hasher.Sum(nil)
}
