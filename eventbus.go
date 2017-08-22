package eventbus

import (
	"fmt"
	"log"
	"net"
)

//--------------------------------------------------
// EventBus
//--------------------------------------------------

type EventBus interface {
	Start() error
	Stop() error
}

type eventBusImpl struct {
	host   string
	port   int
	listen net.Listener
}

//--------------------------------------------------
// Methods
//--------------------------------------------------

func (bus *eventBusImpl) Start() error {
	//log
	log.Printf("Staring EventBus")
	//
	for {
		//log
		log.Printf("Wait connection")
		//wait
		con, _ := bus.listen.Accept()
		//log
		log.Printf("Connection done: %v", con.RemoteAddr())
		//run in work
		go (func(conn net.Conn) {
			wkr := NewWorker(con)
			wkr.Run()
		})(con)
	}
}

func (bus *eventBusImpl) Stop() error {
	return bus.listen.Close()
}

//--------------------------------------------------
// Public Operations
//--------------------------------------------------

func NewEventBus(host string, port int) (EventBus, error) {
	//log
	log.Printf("Instance EventBus in %s:%d", host, port)
	//init listen
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}
	//
	bus := new(eventBusImpl)
	bus.listen = ln
	//
	return bus, nil
}
