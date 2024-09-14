package tests

import (
	"fmt"
	"gopkg.in/zeromq/goczmq.v4"
	"log"
	"testing"
)

func TestRev(t *testing.T) {
	dealer, err := goczmq.NewRep("tcp://localhost:5556")
	if err != nil {
		log.Fatal(err)
	}
	defer dealer.Destroy()

	for {
		msg, err := dealer.RecvMessage()
		if err != nil {
			log.Println("Error receiving message:", err)
			continue
		}
		dealer.SendMessage(msg)
		fmt.Println("Consumer received:", string(msg[0]))
	}
}
