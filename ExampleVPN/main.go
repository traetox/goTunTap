package main

import (
	"fmt"
	"os"
	"flag"
	"net"
	tap "github.com/traetox/goTunTap"
)

var (
	br	=	flag.String("br", "br0", "Bridge in charge of taps")
	brIP	=	flag.String("br-ip", "172.19.0.1", "IP to set bridge")
	gw	=	flag.String("gw", "172.19.0.1", "Gateway for network")
	cidr	=	flag.String("cidr", "172.19.0.0/24", "CIDR for network virtual network")
	listen	=	flag.String("s", ":5150", "IP Port to serve clients on")
	authString string
	nwk *net.IPNet
	tapNames chan string
)

func init() {
	var err error
	flag.Parse()
	//Check all the values to ensure they are legit
	_, nwk, err = net.ParseCIDR(*cidr)
	if(err != nil) {
		fmt.Printf("%s is an invalid CIDR\n", *cidr)
		os.Exit(-1)
	}
	brip := net.ParseIP(*brIP)
	if(brip == nil) {
		fmt.Printf("%s is an invalid IP for the bridge\n", *brIP)
		os.Exit(-1)
	}
	gwip := net.ParseIP(*gw)
	if(gwip == nil) {
		fmt.Printf("%s is an invalid IP for the gateway\n", *gw)
		os.Exit(-1)
	}
	if(!nwk.Contains(gwip)) {
		fmt.Printf("%s is not part of the %s subnet\n", *gw, *cidr)
		os.Exit(-1)
	}
	if(!nwk.Contains(brip)) {
		fmt.Printf("WARNING: bridge IP %s is not part of the network %s\n", *brIP, cidr)

	}
	_, _, err = net.SplitHostPort(*listen)
	if(err != nil) {
		fmt.Printf("%s is an invalid listen parameter\n", *listen)
		os.Exit(-1)
	}
}

func main() {
	fmt.Printf("Enter authorization string: ")
	fmt.Scanf("%s", authString)
	err := tap.CreateBridge(*br)
	if(err != nil) {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(-1)
	}
	tapNames = make(chan string, 2)
	go tapNamer(tapNames)
	err = ListenAndServe(*listen, *br, authString, handleClient)
	if(err != nil) {
		fmt.Printf("ERROR starting up server: %s\n", err)
		os.Exit(-1)
	}
}

func handleClient(bridge string, conn net.Conn) {
	var tapName string
	if tap.CheckBridge(bridge) != nil {
		fmt.Printf("Bridge is down or could not be created")
		return
	}
	tapName = <-tapNames
	tuntap, err := tap.CreateTap(tapName)
	if err != nil {
		fmt.Printf("Failed to create tap: %s\n", err)
		return
	}

	err = tuntap.Start()
	if err != nil {
		fmt.Printf("Failed to create tap: %s\n", err)
		return
	}

	err = tap.AddTapToBridge(bridge, tapName)
	if err != nil {
		fmt.Printf("Failed to add %s to bridge %s\n", tapName, bridge)
		err = tuntap.Stop()
		if err != nil {
			fmt.Printf("Failed to destroy tap: %s\n", err)
			return
		}
	}

	for {
		//do some reading and writing and shit until connection breaks down
	}

	tap.RemoveTapFromBridge(bridge, tapName)
	if err != nil {
		fmt.Printf("Failed to remove %s from bridge %s\n", tapName, bridge)
		return
	}

	err = tuntap.Stop()
	if err != nil {
		fmt.Printf("Failed to destroy tap: %s\n", err)
		return
	}
}

func tapExists(tapname string) bool {
	fi, err := os.Stat(fmt.Sprintf("/sys/class/net/%s/", tapname))
	if(err != nil) {
		return false
	}
	if(fi.IsDir()) {
		return false
	}
	return true
}

func tapNamer(nameChan chan string) {
	var tapName string
	i := uint32(0)
	for {
		for {
			tapName = fmt.Sprintf("tSrv%x", i)
			i++
			if(tapExists(tapName)) {
				continue
			}
			break
		}
		nameChan <- tapName
	}
}
