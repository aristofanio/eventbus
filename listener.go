package eventbus

//--------------------------------------------------
// Listener
//--------------------------------------------------

type Callback func(event Event)

type Listener interface {
	On(e EventType, callback Callback) error
}
