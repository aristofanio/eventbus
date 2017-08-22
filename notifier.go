package eventbus

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
)

//--------------------------------------------------
// Notifier
//--------------------------------------------------

type Notifier interface {
	Notify(event Event) error
}

type notifierImpl struct {
	host string
	port int
}

//--------------------------------------------------
// Methods
//--------------------------------------------------

func (n *notifierImpl) Notify(event Event) error {
	//
	addr := net.TCPAddr{Port: n.port, IP: net.ParseIP(n.host)}
	conn, err := net.Dial("tcp", addr.String())
	if err != nil {
		return err
	}
	//
	buf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buf)
	err = encoder.Encode(event)
	if err != nil {
		return err
	}
	//
	frameOut := NewFrame(EvtReqFrameType, buf.Bytes())
	err = WriteFrame(conn, &frameOut)
	if err != nil {
		return err
	}
	//
	frameIn, err := ReadFrame(conn)
	if err != nil {
		return err
	}
	//
	if frameIn.Type == ErrRespFrameType {
		return errors.New(string(frameIn.Data))
	}
	//
	return nil
}

//--------------------------------------------------
// Public Operations
//--------------------------------------------------

func NewNotifier(host string, port int) Notifier {
	n := new(notifierImpl)
	n.host = host
	n.port = port
	return n
}
