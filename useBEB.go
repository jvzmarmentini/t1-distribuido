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

	BEB.NewBEB(addresses, id, false)

	time.Sleep(3 * time.Second)

	// Keep the main program running
	for true{}
}