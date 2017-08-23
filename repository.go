package eventbus

import (
	"log"
	"sync"
)

//--------------------------------------------------
// Listener Repository
//--------------------------------------------------

type Repository interface {
	Keep(typ EventType, l Listener)
	Find(typ EventType) []Listener
	Remove(uuid string)
}

type repositoryInst struct {
	listeners map[EventType][]Listener
}

//--------------------------------------------------
// Methods
//--------------------------------------------------

func (r *repositoryInst) Keep(typ EventType, ln Listener) {
	//listener impl
	l := ln.(*listenerImpl)
	//log
	log.Printf("Persisting listener %s on %s", l.UUID, l.Type.Name)
	//add listeners
	r.listeners[typ] = append(r.listeners[typ], ln)
}

func (r *repositoryInst) Find(typ EventType) []Listener {
	//log
	log.Printf("Finding listener for %s", typ.Name)
	//
	return r.listeners[typ]
}

func (r *repositoryInst) Remove(uuid string) {
	//log
	log.Printf("Removing listener %s", uuid)
	//seek in listeners
	for k, la := range r.listeners {
		//create new listener arrays
		lanew := make([]Listener, 0)
		//seek of 0 to len(vs)
		for i := 0; i < len(la); i++ {
			l := la[i].(*listenerImpl)
			if l.UUID == uuid {
				l.Conn.Close()
				continue
			} else {
				lanew = append(lanew, l)
			}
		}
		r.listeners[k] = lanew
	}
}

//--------------------------------------------------
// Singleton Implementation (lazy)
// by: http://marcio.io/2015/07/singleton-pattern-in-go/
//--------------------------------------------------
var repInst *repositoryInst
var repOnce = new(sync.Once)

func GetRepositoryInst() Repository {
	repOnce.Do(func() {
		repInst = new(repositoryInst)
		repInst.listeners = make(map[EventType][]Listener)
	})
	return repInst
}
