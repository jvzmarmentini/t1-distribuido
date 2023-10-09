package main

import (
	"os"
	"time"
	"strconv"
	CRB "SD/CRB"
)

func main() {
	id, _ := strconv.Atoi(os.Args[1])
	addresses := os.Args[2:]

	crb := CRB.NewCORB(addresses, id, true)

	time.Sleep(3 * time.Second)
	msg := CRB.Message{
		Msg:   "Hello from " + strconv.Itoa(id) + "!",
		Timestamp: nil,
	}

	// Keep the main program running
	for {
		crb.BroadcastMessageChannel <- msg
		time.Sleep(1 * time.Second)
	}
}