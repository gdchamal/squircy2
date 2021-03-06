package script

import (
	"errors"
	"time"

	"log"

	"github.com/robertkrimen/otto"
	"github.com/tyler-sommer/squircy2/squircy/event"
)

const maxExecutionTime = 5 // in seconds
var Halt = errors.New("Execution limit exceeded")
var UnknownScriptType = errors.New("Unknown script type")

type scriptDriver interface {
	Handle(e event.Event, fnName string)
	RunUnsafe(code string) (interface{}, error)
	String() string
}

type javascriptDriver struct {
	vm *jsVm
	*log.Logger
}

func (d javascriptDriver) Handle(e event.Event, fnName string) {
	d.vm.Interrupt = make(chan func(), 1)
	data, err := d.vm.ToValue(e.Data)
	if err != nil {
		d.Println("An error occurred while creating event data", err)
		return
	}
	_, err = d.vm.Call(fnName, otto.NullValue(), data)
	if err != nil {
		d.Println("An error occurred while executing the Javascript handler", err)
	}
}

func (d javascriptDriver) RunUnsafe(unsafe string) (val interface{}, err error) {
	start := time.Now()
	defer func() {
		if e := recover(); e != nil {
			if e == Halt {
				duration := time.Since(start)
				d.Println("Some Javascript code took too long! Stopping after: ", duration)
			}
			err = e.(error)
		}
		d.vm.Interrupt = make(chan func(), 1)
	}()

	d.vm.Interrupt = make(chan func(), 1)

	go func() {
		time.Sleep(maxExecutionTime * time.Second)
		d.vm.Interrupt <- func() {
			panic(Halt)
		}
	}()

	v, err := d.vm.Run(unsafe)
	val, _ = v.Export()

	return
}

func (d javascriptDriver) String() string {
	return "js"
}
