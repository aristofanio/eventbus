package eventbus

import (
	"net"
)

//--------------------------------------------------
// Worker
//--------------------------------------------------

type ConnWorker interface {
	Run()
}

type connWorkerImpl struct {
	conn net.Conn
}

//--------------------------------------------------
// Methods
//--------------------------------------------------

func (w connWorkerImpl) Run() {
	//ler os dados
	//verificar o tipo (0 - listener/1 - notifier)
	//se for listener,
	//  - recuperar o tipo de evento desejado
	//  - armazenar a conexão de acordo com um tipo de evento
	//  - devolver o código do registro
	//se for notifier,
	//  - recuperar o tipo de evento
	//  - recuperar o código de registro do notificador
	//  - recuperar os dados (json)
	//  - armazenar todos os dados
	//  - recuperar o código de registro do evento
	//  - delegar para o notificador local
}

//--------------------------------------------------
// Public Operations
//--------------------------------------------------

func NewWorker(conn net.Conn) ConnWorker {
	w := new(connWorkerImpl)
	w.conn = conn
	return w
}
