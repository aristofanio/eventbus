package eventbus

//--------------------------------------------------
// Notifier
//--------------------------------------------------

type Notifier interface {
	Notify(event Event) error
}
