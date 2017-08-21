package eventbus

//--------------------------------------------------
// EventBus
//--------------------------------------------------

type EventBus interface {
	Start() error
	Stop() error
}

type eventBusImpl struct {
	host string
	port int
}

//--------------------------------------------------
// Methods
//--------------------------------------------------

func (bus eventBusImpl) Start() error {
	//abrir um socket
	//aguarda uma conexão
	//ao se conectar cria uma tarefa
	//worker
}

func (bus eventBusImpl) Stop() error {
	//encerrar um socket
}

//--------------------------------------------------
// Public Operations
//--------------------------------------------------

func NewEventBus(host string, port int) EventBus {

}
