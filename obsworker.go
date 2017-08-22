package eventbus

import (
	"bytes"
	"encoding/json"
	"log"
)

//--------------------------------------------------
// Storage Observer
//--------------------------------------------------

type ObserverWorker interface {
	Notify(e Event)
}

type observerWorkerImpl struct {
	rep Repository
}

//--------------------------------------------------
// Methods
//--------------------------------------------------

func (o *observerWorkerImpl) Notify(e Event) {
	listeners := o.rep.Find(e.Type)
	for i := 0; i < len(listeners); i++ {
		//
		buf := bytes.NewBuffer([]byte{})
		enc := json.NewEncoder(buf)
		err := enc.Encode(e)
		if err != nil {
			log.Printf("Error on notify %s", e.Type.Name)
			continue
		}
		//
		l := listeners[i]
		if l.Conn == nil {
			log.Printf("Listener = nil (%d/%d)", i, len(listeners))
		}
		n, err := l.Conn.Write(buf.Bytes())
		if err != nil {
			log.Printf("Error on write notification %s", e.Type.Name)
		}
		if n == 0 {
			o.rep.Remove(l.UUID)
		}
	}
}

//--------------------------------------------------
// Public Operations
//--------------------------------------------------

func NewObserverWorker(r Repository) ObserverWorker {
	obs := new(observerWorkerImpl)
	obs.rep = r
	return obs
}
