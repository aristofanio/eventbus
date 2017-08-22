package eventbus

import (
	"net"
	"sync"
)

//--------------------------------------------------
// Listener Repository
//--------------------------------------------------

type Repository interface {
	Keep(uuid string, conn net.Conn)
	Find(uuid string) net.Conn
	Remove(uuid string)
}

type repositoryInst struct {
	lns map[string]net.Conn
}

//--------------------------------------------------
// Methods
//--------------------------------------------------

func (r *repositoryInst) Keep(uuid string, conn net.Conn) {
	//
	if r.lns[uuid] != nil {
		c := r.lns[uuid]
		c.Close()
	}
	//
	r.lns[uuid] = conn
}

func (r *repositoryInst) Find(uuid string) net.Conn {
	return r.lns[uuid]
}

func (r *repositoryInst) Remove(uuid string) {
	//
	if r.lns[uuid] != nil {
		c := r.lns[uuid]
		c.Close()
	}
	//
	r.lns[uuid] = nil
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
		repInst.lns = make(map[string]net.Conn)
	})
	return repInst
}
