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

type Callback func(event *Event, err error)

type Listener interface {
	On(t EventType, callback Callback)
}

type listenerImpl struct {
	UUID string    `json:"uuid"`
	Type EventType `json:"type"`
	Conn net.Conn  `json:"-"`
}

//--------------------------------------------------
// Methods
//--------------------------------------------------

func (l *listenerImpl) On(t EventType, callback Callback) {
	//convert in bytes
	b, err := fromListener(l, t)
	if err != nil {
		callback(nil, err)
	}
	//log
	log.Printf("Registering listener %s on %s", l.UUID, l.Type.Name)
	//create frame
	frameOut := NewFrame(RegReqFrameType, b)
	//write frame in conn (registering)
	err = WriteFrame(l.Conn, frameOut)
	if err != nil {
		callback(nil, err)
	}
	//
	go (func() {
		//loop
		for {
			//log
			log.Printf("Waiting event: %s", l.Type.Name)
			//read conn
			frameIn, err := ReadFrame(l.Conn)
			if err != nil {
				callback(nil, err)
			}
			//switch action by frame type
			switch frameIn.Type {
			case RegRespFrameType:
				//convert data in frame to event
				e, err := toEvent(frameIn.Data)
				if err != nil {
					callback(nil, err)
				}
				//log
				log.Printf("Execute callback")
				//exec callback
				callback(e, nil)
			case PngReqFrameType:
				//create frame
				frameResp := NewFrame(PngRespFrameType, []byte{})
				err := WriteFrame(l.Conn, frameResp)
				if err != nil {
					callback(nil, err)
				}
			case ErrRespFrameType:
				callback(nil, errors.New(string(frameIn.Data)))
			default:
				callback(nil, errors.New("Fame type unknown"))
			}
		}
	})()
}

//--------------------------------------------------
// Private Operations
//--------------------------------------------------

func fromListener(l *listenerImpl, etyp EventType) ([]byte, error) {
	//
	l.Type = etyp
	//
	buf := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(buf)
	err := enc.Encode(l)
	if err != nil {
		return nil, err
	}
	//
	return buf.Bytes(), nil
}

func toListener(conn net.Conn, b []byte) (Listener, error) {
	//create new listener
	l := new(listenerImpl)
	//convert bytes to listener
	buf := bytes.NewBuffer(b)
	dec := json.NewDecoder(buf)
	err := dec.Decode(l)
	if err != nil {
		return nil, err
	}
	//set connection
	l.Conn = conn
	//result
	return l, nil
}

//--------------------------------------------------
// Public Operations
//--------------------------------------------------

func NewListener(host string, port int, uuid string) (Listener, error) {
	//setup conn
	addr := net.TCPAddr{Port: port, IP: net.ParseIP(host)}
	conn, err := net.Dial("tcp", addr.String())
	if err != nil {
		return nil, err
	}
	//instance listener
	ln := new(listenerImpl)
	ln.UUID = uuid
	ln.Conn = conn
	//result
	return ln, nil
}
