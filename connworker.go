package eventbus

import (
	"bytes"
	"log"
	"net"
)

//--------------------------------------------------
// Worker
//--------------------------------------------------

type ConnWorker interface {
	Run() error
}

type connWorkerImpl struct {
	conn net.Conn
}

//--------------------------------------------------
// Methods
//--------------------------------------------------

func (w connWorkerImpl) Run() error {
	//log
	log.Printf("Running...")
	//read frame
	frame, err := ReadFrame(w.conn)
	if err != nil {
		w.conn.Close()
		return err
	}
	//log
	log.Printf("Running...")
	//verificar o tipo
	switch frame.Type {
	case RegReqFrameType:
		//registration req  (listener)
		//se for listener,
		//  - recuperar o tipo de evento desejado
		//  - armazenar a conexão de acordo com um tipo de evento
		uuid := string(frame.Data)
		rep := GetRepositoryInst()
		rep.Keep(uuid, w.conn)
	case EvtReqFrameType:
		//event req (notifier)
		//se for notifier,
		//  - recuperar o tipo de evento
		//  - recuperar o código de registro do notificador
		//  - recuperar os dados (json)
		//  - armazenar todos os dados
		//  - recuperar o código de registro do evento
		//  - delegar para o notificador local
		defer w.conn.Close()
		event, err := unmarshalEvent(frame.Data)
		if err != nil {
			//resp K - OK, D - Event, P - Ping and E - ERROR
			b := make([]byte, 0)
			buf := bytes.NewBuffer(b)
			buf.WriteByte('E')
			buf.WriteString(err.Error())
			frame := NewFrame(EvtRespFrameType, buf.Bytes())
			WriteFrame(w.conn, &frame) //todo convert to json-data
			//todo notification
		} else {
			sto := GetStorageInst()
			sto.StoreEvent(event)
			//resp K - OK, D - Event, P - Ping and E - ERROR
			WriteFrame(w.conn, &Frame{Type: EvtRespFrameType, Size: 1, Data: []byte{'K'}}) //todo convert to json-data
		}
	case PngReqFrameType:
		//ping req (notifier | listener)
	default:
		w.conn.Close()
	}
	//
	return nil

}

//--------------------------------------------------
// Public Operations
//--------------------------------------------------

func NewWorker(conn net.Conn) ConnWorker {
	//log
	log.Printf("Create work")
	//
	w := new(connWorkerImpl)
	w.conn = conn
	return w
}
