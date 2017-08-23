package eventbus

import (
	"log"
)

//--------------------------------------------------
// Storage Observer
//--------------------------------------------------

type ObserverWorker interface {
	Notify(e *Event)
}

type observerWorkerImpl struct {
	rep Repository
}

//--------------------------------------------------
// Methods
//--------------------------------------------------

func (o *observerWorkerImpl) Notify(e *Event) {
	la := o.rep.Find(e.Type)
	for i := 0; i < len(la); i++ {
		//convert event to bytes
		b, err := fromEvent(e)
		if err != nil {
			log.Printf("Error on notify listener about %s event", e.Type.Name)
		}
		//create frame
		frameOut := NewFrame(RegRespFrameType, b)
		//write in listener
		l := la[i].(*listenerImpl)
		err = WriteFrame(l.Conn, frameOut)
		if err != nil {
			log.Printf("Error on write notification %s", e.Type.Name)
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
