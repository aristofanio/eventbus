package eventbus

import "sync"

//--------------------------------------------------
// Storage
//--------------------------------------------------

type Storage interface {
	StoreEvent(e *Event) error
	ListEvents(t EventType) ([]*Event, error)
}

type storageInst struct {
	events map[EventType][]*Event
}

//--------------------------------------------------
// Methods
//--------------------------------------------------

//StoreEvent store events in memory
func (s *storageInst) StoreEvent(e *Event) error {
	es := s.events[e.Type]
	es = append(es, e)
	s.events[e.Type] = es
	return nil
}

//ListEvents list all event in memory filtered by type t
func (s *storageInst) ListEvents(t EventType) ([]*Event, error) {
	return s.events[t], nil
}

//--------------------------------------------------
// Singleton Implementation (lazy)
// by: http://marcio.io/2015/07/singleton-pattern-in-go/
//--------------------------------------------------
var stoInst *storageInst
var stoOnce = new(sync.Once)

//GetStorageInst get events store instance
func GetStorageInst() Storage {
	stoOnce.Do(func() {
		stoInst = new(storageInst)
		stoInst.events = make(map[EventType][]*Event)
	})
	return stoInst
}
