//go:generate genny -in=$GOFILE -out=../common/int_observable.go gen "ValType=int pkg=common"
//go:generate genny -in=$GOFILE -out=../common/string_observable.go gen "ValType=string pkg=common"
package pkg

import "github.com/cheekybits/genny/generic"

type ValType generic.Type
type pkg generic.Type

type ValTypeEventEmitter struct {
	subscribers []ValTypeSubscriber
}

type ValTypeSubscriber interface {
	EventFired(ValType)
}

func NewValTypeEventEmitter() *ValTypeEventEmitter {
	return &ValTypeEventEmitter{make([]ValTypeSubscriber, 0, 1)}
}

func (emitter *ValTypeEventEmitter) Subscribe(subscriber ValTypeSubscriber) {
	emitter.subscribers = append(emitter.subscribers, subscriber)
}

func (emitter *ValTypeEventEmitter) FireEvent(event ValType) {
	for _, sub := range emitter.subscribers {
		sub.EventFired(event)
	}
}
