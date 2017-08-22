package eventbus

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
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
	UUID string    `json:"uuid"`
	Type EventType `json:"type"`
	Conn net.Conn  `json:"-"`
}

//--------------------------------------------------
// Methods
//--------------------------------------------------

func (ln *listenerImpl) On(t EventType, callback Callback) error {
	//
	ln.Type = t
	//
	buf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buf)
	err0 := encoder.Encode(ln)
	if err0 != nil {
		return err0
	}
	//log
	log.Printf("Registering listener %s on %s", ln.UUID, ln.Type.Name)
	//
	frameOut := NewFrame(RegReqFrameType, buf.Bytes())
	err1 := WriteFrame(ln.Conn, &frameOut)
	if err1 != nil {
		return err1
	}
	//
	for {
		//
		log.Printf("Waiting event %s", ln.Type.Name)
		//
		frameIn, err := ReadFrame(ln.Conn)
		if err != nil {
			return err
		}
		//
		switch frameIn.Type {
		case RegRespFrameType:
			//
			databuf := bytes.NewBuffer(frameIn.Data)
			decoder := json.NewDecoder(databuf)
			eventdt := new(Event)
			err := decoder.Decode(eventdt)
			if err != nil {
				return err
			}
			//
			log.Printf("Callback...")
			//
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
