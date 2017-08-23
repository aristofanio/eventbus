# EventBus
EventBus is an simple messaging system in publisher-subscriber style.

## Fases

[![Build Status](http://img.shields.io/travis/badges/badgerbadgerbadger.svg?style=flat-square)]()


## Use

Download source code
```
$ go get -u github.com/aristofanio/eventbus
```

## Example of use

Start server in host=localhost and port=9090:
```Go
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
```Go
package main

import (
	"github.com/aristofanio/eventbus"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	//start listener
	l, err := eventbus.NewListener("localhost", 9090, "ln-019101")
	if err != nil {
		println(err.Error())
	}
	//registry (non-blocking)
	l.On(eventbus.EventType{"onTest"}, func(e *eventbus.Event, err error) {
		log.Printf("---> result: [type: %s, data: %s]", e.Type, string(e.Data))
	})
	//for avoid inexpected end
	ioutil.ReadAll(os.Stdin)
}
```

An Notifier:
```Go
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
