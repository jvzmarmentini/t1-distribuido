package main

import (
	"os"
	"time"
	"strconv"
	BEB "SD/BEB"
	CRB "SD/CRB"
)

func main() {
	id, _ := strconv.Atoi(os.Args[1])
	addresses := os.Args[2:]

	crb := CRB.NewCORB(addresses, id, true)

	time.Sleep(3 * time.Second)
	msg := BEB.BroadcastMessage{
		SenderID:   id,
		Message:   "Hello from " + strconv.Itoa(id) + "!",
		Timestamp: nil,
	}

	for {
		crb.BEBModule.BroadcastMessageChannel <- msg
		time.Sleep(1 * time.Second)
	}
}