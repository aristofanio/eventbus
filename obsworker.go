package eventbus

//--------------------------------------------------
// Storage Observer
//--------------------------------------------------

type ObserverWorker interface {
	Notify(e Event)
	Run()
}
