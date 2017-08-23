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
	frameIn, err := ReadFrame(w.conn)
	if err != nil {
		w.conn.Close()
		return err
	}
	//log
	log.Printf("Convert to frame with size %d and type 0x%x: %s", frameIn.Size, frameIn.Type, frameIn.Data)
	//verificar o tipo
	switch frameIn.Type {
	case RegReqFrameType:
		//registration req  (listener)
		//se for listener,
		//  - recuperar o tipo de evento desejado
		//  - armazenar a conexão de acordo com um tipo de evento
		ln, err := toListener(w.conn, frameIn.Data)
		if err != nil {
			return err
		}
		l := ln.(*listenerImpl)
		//log
		log.Printf("Get repository and keep listener %s", l.UUID)
		//repository
		rep := GetRepositoryInst()
		rep.Keep(l.Type, l)
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
		e, err := toEvent(frameIn.Data)
		if err != nil {
			//log
			log.Printf("Error on receving event type")
			//create response data
			byt := make([]byte, 0)
			buf := bytes.NewBuffer(byt)
			buf.WriteByte('E')
			buf.WriteString(err.Error())
			//log
			log.Printf("Error: %s", buf.Bytes())
			//create frame
			frameOut := NewFrame(EvtRespFrameType, buf.Bytes())
			//write in conn (error)
			WriteFrame(w.conn, frameOut) //todo convert to json-data
		} else {
			//log
			log.Printf("Receiving event type %s", e.Type.Name)
			//get event store instance
			sto := GetStorageInst()
			sto.StoreEvent(e)
			//log
			log.Printf("Storing event type")
			//create frame
			frame := NewFrame(EvtRespFrameType, []byte{'K'})
			//log
			log.Printf("Notifing...")
			//write in conn (ok)
			WriteFrame(w.conn, frame)
			//create observe worker
			rep := GetRepositoryInst()
			obs := NewObserverWorker(rep)
			obs.Notify(e)
		}
	case PngReqFrameType:
		//log
		log.Printf("FrameType: PngReqFrameType")
		//ping req (notifier | listener)
		//todo ping
	default:
		//log
		log.Printf("FrameType: unknown")
		//close connection
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
