package eventbus

//--------------------------------------------------
// Connector
//--------------------------------------------------

type Connector interface {
	ConnListener(uuid string) Listener
	ConnNotifier(uuid string) Notifier
}
