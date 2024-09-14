package pool

import (
	"errors"
	"fmt"
	uuid2 "github.com/google/uuid"
	"go-mq/common"
	zmq "gopkg.in/zeromq/goczmq.v4"
	"log"
	"sync"
)

var PoolInstance *ResourcePool

func init() {
	PoolInstance = NewResourcePool(10)
}

type MqResource struct {
	RevPort      int
	SendPort     int
	Id           string
	RouterSocket *zmq.Sock
	DealerSocket *zmq.Sock
	Closed       bool
}

func (this *MqResource) Init(repServerList []string) {
	var err error
	uuid, err := uuid2.NewRandom()
	if err != nil {
		fmt.Println("Creating new resource err: ", err.Error())
		return
	}
	this.Id = uuid.String()
	this.RevPort = common.FindNextFreePort(5555)
	this.RouterSocket, err = zmq.NewRouter(fmt.Sprintf("tcp://*:%d", this.RevPort))
	if err != nil {
		fmt.Println(err)
		return
	}
	this.SendPort = common.FindNextFreePort(5555)

	dealer := zmq.NewSock(zmq.Dealer)
	for _, server := range repServerList {
		err = dealer.Connect(server)
		if err != nil {
			log.Fatalf("cannot link REP server %s: %v", server, err)
		}
		fmt.Printf("link REP server: %s\n", server)
	}
	this.DealerSocket, err = zmq.NewDealer(fmt.Sprintf("tcp://*:%d", this.SendPort))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (this *MqResource) Run() {
	go func() {
		for {
			msg, err := this.RouterSocket.RecvMessage()
			if err != nil {
				fmt.Println("Error receiving message from producer:", err)
				continue
			}

			fmt.Println("Broker received from producer:", string(msg[0]), string(msg[2]))

			//time.Sleep(time.Millisecond * 500)

			err = this.DealerSocket.SendMessage(msg)
			if err != nil {
				fmt.Println("Error forwarding message to consumer:", err)
			}
		}
	}()

	go func() {
		for {
			msg, err := this.DealerSocket.RecvMessage()
			if err != nil {
				fmt.Println("Error receiving message from backend: %v", err)
				break
			}
			fmt.Printf("Received message at backend: %v\n", msg)

			err = this.RouterSocket.SendMessage(msg)
			if err != nil {
				log.Printf("Error sending message to frontend: %v", err)
				break
			}
		}
	}()
}

func (this *MqResource) Close() {
	this.Closed = true
}

type ResourcePool struct {
	pool   chan *MqResource
	mu     sync.Mutex
	closed bool
}

func NewResourcePool(maxSize int) *ResourcePool {
	return &ResourcePool{
		pool: make(chan *MqResource, maxSize),
	}
}

func (p *ResourcePool) Get(ipList []string) (*MqResource, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return nil, errors.New("resource pool is closed")
	}

	select {
	case res := <-p.pool:
		fmt.Println("Reusing existing resource:", res.Id)
		return res, nil
	default:
		newResource := new(MqResource)
		newResource.Init(ipList)
		newResource.Run()
		fmt.Println("Creating new resource:", newResource.Id)
		return newResource, nil
	}
}

func (p *ResourcePool) Put(res *MqResource) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return errors.New("resource pool is closed")
	}

	select {
	case p.pool <- res:
		fmt.Println("MqResource returned to pool:", res.Id)
		return nil
	default:
		fmt.Println("MqResource pool is full, discarding resource:", res.Id)
		return nil
	}
}

func (p *ResourcePool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.closed {
		close(p.pool)
		p.closed = true
	}
}
