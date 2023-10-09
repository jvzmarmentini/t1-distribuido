package CORB

import (
	"fmt"
	BEB "SD/BEB" // Import the BEB module
)

type CORB_Module struct {
	ID          int
	BEBModule   *BEB.BEB_Module // Best Effort Broadcast module
	VectorClock []int           // Vector clock for causal ordering
	Buffer      []BEB.BroadcastMessage
	Acknowledgments map[int]int // Acknowledgments for reliability
}

type Message struct {
	Msg   string
	Timestamp []int
}

func NewCORB(_addresses []string, _id int, _dbg bool) *CORB_Module {
	if _dbg {
		fmt.Println("Initializing CORB")
	}

	beb := BEB.NewBEB(_addresses, _id, _dbg)

	corb := &CORB_Module{
		ID:              _id,
		BEBModule:       beb,
		VectorClock:     make([]int, len(_addresses)),
		Buffer:          make([]BEB.BroadcastMessage, 0),
		Acknowledgments: make(map[int]int),
	}

	corb.Start()

	return corb
}

func (corb *CORB_Module) Start() {
	// Function to handle incoming BEB broadcast messages
	go func() {
		for {
			message := <-corb.BEBModule.BroadcastChannel
			corb.Deliver(message)
		}
	}()

	// Function to broadcast messages
	go func() {
		for {
			message := <-corb.BEBModule.BroadcastMessageChannel
			corb.Broadcast(message)
		}
	}()
}

func (corb *CORB_Module) Broadcast(message string) {
	// Increment the local vector clock
	corb.VectorClock[corb.ID]++
	localTimestamp := make([]int, len(corb.VectorClock))
	copy(localTimestamp, corb.VectorClock)

	// Create a broadcast message with the message content and local vector clock
	broadcastMessage := BEB.BroadcastMessage{
		SenderID:   corb.ID,
		Timestamp:  localTimestamp,
		Message:    message,
	}

	// Send the broadcast message using the BEB module
	corb.BEBModule.BroadcastChannel <- broadcastMessage
}

func (corb *CORB_Module) Deliver(message BEB.BroadcastMessage) {
	// Check if the message can be delivered (causal order)
	canDeliver := true
	for i := range message.Timestamp {
		if i != message.SenderID {
			if message.Timestamp[i] > corb.VectorClock[i]+1 {
				canDeliver = false
				break
			}
		}
	}

	if canDeliver {
		// Deliver the message
		fmt.Printf("Process %d received from %d: %s\n", corb.ID, message.SenderID, message.Message)

		// Update the local vector clock
		for i := range message.Timestamp {
			corb.VectorClock[i] = max(corb.VectorClock[i], message.Timestamp[i])
		}

		// Check if buffered messages can now be delivered
		deliverableMessages := make([]BEB.BroadcastMessage, 0)
		for _, bufferedMsg := range corb.Buffer {
			bufferedMsgCanDeliver := true
			for i := range bufferedMsg.Timestamp {
				if i != bufferedMsg.SenderID {
					if bufferedMsg.Timestamp[i] > corb.VectorClock[i]+1 {
						bufferedMsgCanDeliver = false
						break
					}
				}
			}
			if bufferedMsgCanDeliver {
				// Deliver the buffered message
				fmt.Printf("Process %d received from %d (Buffered): %s\n", corb.ID, bufferedMsg.SenderID, bufferedMsg.Message)
				// Update the local vector clock
				for i := range bufferedMsg.Timestamp {
					corb.VectorClock[i] = max(corb.VectorClock[i], bufferedMsg.Timestamp[i])
				}
			} else {
				// Add the non-deliverable message back to the buffer
				deliverableMessages = append(deliverableMessages, bufferedMsg)
			}
		}

		// Update the buffer with the remaining undelivered messages
		corb.Buffer = deliverableMessages
	} else {
		// Buffer the message for later delivery
		corb.Buffer = append(corb.Buffer, message)
	}
}


func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
