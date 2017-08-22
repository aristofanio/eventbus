# EventBus
EventBus is an simple messaging system in publisher-subscriber style.

## Fases

[![Codeship](https://img.shields.io/codeship/d6c1ddd0-16a3-0132-5f85-2e35c05e22b1.svg?style=plastic)]()


## Use

Download source code
```
$ go get -u github.com/aristofanio/eventbus
```

## Example of use

Start server in host=localhost and port=9090:
```
package main

import "github.com/aristofanio/eventbus"

func main() {
	//configuring server
	bus, err := eventbus.NewEventBus("localhost", 9090)
	if err != nil {
		println(err.Error())
	}
    //start service
	bus.Start()
}

```

An listener:
```
package main

import "github.com/aristofanio/eventbus"

func main() {
	//
	listener, err := eventbus.NewListener("ln-019101", "localhost", 9090)
	if err != nil {
		println(err.Error())
	}
	//
	listener.On(eventbus.EventType{"onTest"}, func(e eventbus.Event) {
		println(e.Type.Name)
	})
	//
	println("opa")
}

```

An Notifier:
```
package main

import "github.com/aristofanio/eventbus"

func main() {
	//
	e := eventbus.Event{
		Type:   eventbus.EventType{"onTest"},
		Origin: "nt-000001",
		Data:   []byte(`{"name": "teste"}`),
	}
	//
	notifier := eventbus.NewNotifier("localhost", 9090)
	err := notifier.Notify(e)
	if err != nil {
		println(err.Error())
	}
}

```
