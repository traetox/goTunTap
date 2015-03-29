# goTunTap
Golang libary for creating, deleting, reading, and writing linux Taps.

The library also allows for creating and managing bridges.  THe goal is
to provide a nice wrapper for creating bridges.

For example, a layer two tunnel can be created as simply as:

## Server side
```go
if err := tap.CreateBridge(bridge_name); err != nil {
	log.Fatal("Failed to create bridge", err)
}
t, err := tap.CreateTap(tap_name)
if err != nil {
	log.Fatal("Failed to create tap manager")
}
if err := t.AddToBridge(bridge_name); err != nil {
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
```


## Client side
```go
conn, err := net.Dial("tcp", remote_server)
if err != nil {
	log.Fatal("Failed to dial server:", err)
}
defer conn.Close()
if err := tap.CreateBridge(bridge_name); err != nil {
	log.Fatal("Failed to create bridge", err)
}
t, err := tap.CreateTap(tap_name)
if err != nil {
	log.Fatal("Failed to create tap manager")
}
if err := t.AddToBridge(bridge_name); err != nil {
	log.Fatal("Failed to add tap to bridge", err)
}
if err := t.Relay(conn); err != nil {
	log.Fatal("Failed to relay tap connection")
}

```
