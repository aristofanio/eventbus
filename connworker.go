package eventbus

import (
	"bytes"
	"encoding/json"
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
		lsn := new(listenerImpl)
		buf := bytes.NewBuffer(frame.Data)
		dec := json.NewDecoder(buf)
		err := dec.Decode(lsn)
		if err != nil {
			return err
		}
		lsn.Conn = w.conn
		//log
		log.Printf("Get repository and keep listener")
		//
		rep := GetRepositoryInst()
		rep.Keep(lsn.Type, *lsn)
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
		//log
		log.Printf("Receving event type %s", event.Type.Name)
		//
		if err != nil {
			//log
			log.Printf("Error on receving event type")
			//resp K - OK, D - Event, P - Ping and E - ERROR
			b := make([]byte, 0)
			buf := bytes.NewBuffer(b)
			buf.WriteByte('E')
			buf.WriteString(err.Error())
			frame := NewFrame(EvtRespFrameType, buf.Bytes())
			WriteFrame(w.conn, &frame) //todo convert to json-data
		} else {
			//log
			log.Printf("Storing event type")
			//
			sto := GetStorageInst()
			sto.StoreEvent(event)
			//
			log.Printf("Notifing...")
			//
			frame := NewFrame(EvtRespFrameType, []byte{'K'})
			WriteFrame(w.conn, &frame) //todo convert to json-data
			//
			rep := GetRepositoryInst()
			obsWorker := NewObserverWorker(rep)
			obsWorker.Notify(*event)
		}
	case PngReqFrameType:
		//ping req (notifier | listener)
		//todo ping
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
