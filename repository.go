package eventbus

import (
	"log"
	"sync"
)

//--------------------------------------------------
// Listener Repository
//--------------------------------------------------

type Repository interface {
	Keep(typ EventType, ln listenerImpl)
	Find(typ EventType) []listenerImpl
	Remove(uuid string)
}

type repositoryInst struct {
	listeners map[EventType][]listenerImpl
}

//--------------------------------------------------
// Methods
//--------------------------------------------------

func (r *repositoryInst) Keep(typ EventType, ln listenerImpl) {
	//log
	log.Printf("Persisting listener %s on %s", ln.UUID, ln.Type.Name)
	//
	r.listeners[typ] = append(r.listeners[typ], ln)
}

func (r *repositoryInst) Find(typ EventType) []listenerImpl {
	//log
	log.Printf("Finding listener for %s", typ.Name)
	//
	return r.listeners[typ]
}

func (r *repositoryInst) Remove(uuid string) {
	//log
	log.Printf("Removing listener %s", uuid)
	//
	for k, vs := range r.listeners {
		vsnew := make([]listenerImpl, 0)
		for i := 0; i < len(vs); i++ {
			v := vs[i]
			if v.UUID == uuid {
				v.Conn.Close()
				continue
			} else {
				vsnew = append(vsnew, v)
			}
		}
		r.listeners[k] = vsnew
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
		repInst.listeners = make(map[EventType][]listenerImpl)
	})
	return repInst
}
