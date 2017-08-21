package eventbus

//--------------------------------------------------
// Event
//--------------------------------------------------

type Event struct{
	Origin string
	Type EventType
	Data byte[]
}

type EventType struct{
	Name string
}
