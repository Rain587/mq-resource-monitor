package tests

import (
	"fmt"
	"gopkg.in/zeromq/goczmq.v4"
	"log"
	"testing"
)

func TestServer(t *testing.T) {
	req, err := goczmq.NewReq("tcp://localhost:5555")
	if err != nil {
		log.Fatal(err)
	}
	defer req.Destroy()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("Message from producer #%d", i)
		err = req.SendMessage([][]byte{[]byte(msg)})
		if err != nil {
			log.Println("Error sending message:", err)
		} else {
			fmt.Println("Sent:", msg)
		}
		revMsg, err := req.RecvMessage()
		if err != nil {
			log.Println("Error receive message:", err)
		} else {
			fmt.Println("Receive:", revMsg)
		}

		//time.Sleep(time.Second)
	}
}
