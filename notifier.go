package eventbus

import (
	"errors"
	"net"
)

//--------------------------------------------------
// Notifier
//--------------------------------------------------

type Notifier interface {
	Notify(e *Event) error
}

type notifierImpl struct {
	host string
	port int
}

//--------------------------------------------------
// Methods
//--------------------------------------------------

func (n *notifierImpl) Notify(e *Event) error {
	//setup conn
	addr := net.TCPAddr{Port: n.port, IP: net.ParseIP(n.host)}
	conn, err := net.Dial("tcp", addr.String())
	if err != nil {
		return err
	}
	//convert event in bytes
	b, err := fromEvent(e)
	if err != nil {
		return err
	}
	//create frame
	frameOut := NewFrame(EvtReqFrameType, b)
	//write frame in conn
	err = WriteFrame(conn, frameOut)
	if err != nil {
		return err
	}
	//read response (error ou ok)
	frameIn, err := ReadFrame(conn)
	if err != nil {
		return err
	}
	//check if exists error
	if frameIn.Type == ErrRespFrameType {
		return errors.New(string(frameIn.Data))
	}
	//result
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
