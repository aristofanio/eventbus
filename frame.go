package eventbus

import (
	"errors"
	"io"
	"log"
	"net"
)

//--------------------------------------------------
// Variables and Constants
//--------------------------------------------------

var _op_type = []byte{0xF0, 0x00} // 1111 0000 0000 0000
var _op_size = []byte{0x0F, 0xFF} // 0000 1111 1111 1111 (4kb = 4095)

const (
	RegReqFrameType  = 0x1000 //registration request 	- 0001 0000 0000 0000
	RegRespFrameType = 0x9000 //registration response 	- 1001 0000 0000 0000
	EvtReqFrameType  = 0x3000 //event request			- 0011 0000 0000 0000
	EvtRespFrameType = 0xb000 //event request			- 1011 0000 0000 0000
	PngReqFrameType  = 0x7000 //ping request			- 0111 0000 0000 0000
	PngRespFrameType = 0xF000 //ping response			- 1111 0000 0000 0000
	ErrRespFrameType = 0x8000 //error response			- 1000 0000 0000 0000
)

//--------------------------------------------------
// Type
//--------------------------------------------------

type Frame struct {
	Type int
	Size int
	Data []byte
}

//--------------------------------------------------
// Private Operations
//--------------------------------------------------

func rtype(b0, b1 byte) int {
	p0 := b0 & _op_type[0]
	p1 := b1 & _op_type[1]
	return int(p0)*256 + int(p1)
}

func rsize(b0, b1 byte) int {
	p0 := b0 & _op_size[0]
	p1 := b1 & _op_size[1]
	return int(p0)*256 + int(p1)
}

func wtypesize(t int, s int) []byte {
	b := (s & 0x0FFF) | t
	//int to byte array
	p0 := (b / 256)
	p1 := b - (p0 * 256)
	//
	pt := make([]byte, 2)
	pt[0] = byte(p0)
	pt[1] = byte(p1)
	//
	return pt
}

func _read(r io.Reader, s int) ([]byte, error) {
	//slice
	result := make([]byte, 0, s)
	//count received
	count := 0
	for {
		//
		input := make([]byte, s)
		c, err := r.Read(input)
		if err != nil {
			return nil, err
		}
		if c == 0 {
			return nil, errors.New("Closed connection, EOF stream or badly formatted data")
		}
		//
		result = append(result, input[:c]...)
		count += c
		//
		if count >= s {
			break
		}
	}
	//result
	return result, nil
}

func _readFrame(r io.Reader) (*Frame, error) {
	//read 2 bytes
	bs, err := _read(r, 2)
	if err != nil {
		return nil, err
	}
	//translate
	t := rtype(bs[0], bs[1])
	s := rsize(bs[0], bs[1])
	//log
	log.Printf("Reading initial 2 bytes: 0x%x", bs)
	//data
	dt, err := _read(r, s)
	if err != nil {
		return nil, err
	}
	//log
	log.Printf("Reading data with %d bytes: %s", s, dt)
	//
	fr := new(Frame)
	fr.Type = t
	fr.Size = s
	fr.Data = dt
	//result
	return fr, nil
}

func _writeFrame(w io.Writer, f *Frame) error {
	//write 2 bytes
	bs := wtypesize(f.Type, f.Size)
	w.Write(bs)
	w.Write(f.Data)
	return nil
}

//--------------------------------------------------
// Public Operations
//--------------------------------------------------

func NewFrame(t int, b []byte) Frame {
	f := Frame{Type: t, Size: len(b), Data: b}
	return f
}

func ReadFrame(conn net.Conn) (*Frame, error) {
	return _readFrame(conn)
}

func WriteFrame(conn net.Conn, frame *Frame) error {
	return _writeFrame(conn, frame)
}
