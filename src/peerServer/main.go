package main

import (
	"peerServer/server"
)

func main() {
	InitLog()

	defer func() {
	}()

	server.Start()
}
