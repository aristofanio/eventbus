package eventbus

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
)

//--------------------------------------------------
// Listener
//--------------------------------------------------

type Callback func(event Event)

type Listener interface {
	On(e EventType, callback Callback) error
}

type listenerImpl struct {
	UUID string
	Conn net.Conn
}

//--------------------------------------------------
// Methods
//--------------------------------------------------

func (ln *listenerImpl) On(e EventType, callback Callback) error {
	//
	frameOut := NewFrame(RegReqFrameType, []byte(ln.UUID))
	err := WriteFrame(ln.Conn, &frameOut)
	if err != nil {
		return err
	}
	//
	for {
		//
		frameIn, err := ReadFrame(ln.Conn)
		if err != nil {
			return err
		}
		//
		switch frameIn.Type {
		case RegRespFrameType:
			databuf := bytes.NewBuffer(frameIn.Data)
			decoder := json.NewDecoder(databuf)
			eventdt := new(Event)
			err := decoder.Decode(eventdt)
			if err != nil {
				return err
			}
			callback(*eventdt)
		case PngReqFrameType:
			frameResp := NewFrame(PngRespFrameType, []byte{})
			err := WriteFrame(ln.Conn, &frameResp)
			if err != nil {
				return err
			}
		case ErrRespFrameType:
			return errors.New(string(frameIn.Data))
		default:
			return errors.New("Fame type unknown")
		}
	}
}

//--------------------------------------------------
// Public Operations
//--------------------------------------------------

func NewListener(uuid, host string, port int) (Listener, error) {
	addr := net.TCPAddr{Port: port, IP: net.ParseIP(host)}
	conn, err := net.Dial("tcp", addr.String())
	if err != nil {
		return nil, err
	}
	ln := new(listenerImpl)
	ln.UUID = uuid
	ln.Conn = conn
	return ln, nil
}
