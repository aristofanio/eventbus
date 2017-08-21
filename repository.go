package eventbus

import (
	"net"
)

//--------------------------------------------------
// Listener Repository
//--------------------------------------------------

type Repository interface {
	Keep(uuid string, conn net.Conn)
	Find(uuid string) *net.Conn
	Remove(uuid string)
}
