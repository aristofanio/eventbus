package eventbus

import (
	"bytes"
	"encoding/json"
)

//--------------------------------------------------
// Event
//--------------------------------------------------

type Event struct {
	Origin string    `json:"origin"`
	Type   EventType `json:"type"`
	Data   []byte    `json:"data"`
}

type EventType struct {
	Name string `json:"name"`
}

//--------------------------------------------------
// Private Operations
//--------------------------------------------------

func toEvent(buf []byte) (*Event, error) {
	buffer := bytes.NewBuffer(buf)
	dec := json.NewDecoder(buffer)
	evn := new(Event)
	err := dec.Decode(evn)
	if err != nil {
		return nil, err
	}
	return evn, nil
}

func fromEvent(e *Event) ([]byte, error) {
	byt := make([]byte, 0)
	buf := bytes.NewBuffer(byt)
	enc := json.NewEncoder(buf)
	err := enc.Encode(e)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
