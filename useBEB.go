package main

import (
	"os"
	"time"
	"strconv"
	BEB "SD/BEB"
)

func main() {
	id, _ := strconv.Atoi(os.Args[1])
	addresses := os.Args[2:]

	var beb = BEB.NewBEB(addresses, id, false)

	time.Sleep(3 * time.Second)

	msg := BEB.Message{
		Msg:   "Hello from " + strconv.Itoa(id) + "!",
		Timestamp: nil,
	}

	// Keep the main program running
	for {
		beb.BroadcastMessageChannel <- msg
		time.Sleep(1 * time.Second)
	}
}