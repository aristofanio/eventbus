package eventbus

//--------------------------------------------------
// Storage
//--------------------------------------------------

type Storage interface {
	StoreListener(l Listener) error
	StoreEvent(e Event) error
}
