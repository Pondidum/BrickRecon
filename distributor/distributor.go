package distributor

import (
	"reflect"
	"sync"
)

type MessageHandlerFunc func(Message)

type Distributor struct {
	topics map[string][]MessageHandlerFunc
}

func NewDistributor() *Distributor {
	return &Distributor{
		topics: map[string][]MessageHandlerFunc{},
	}
}

func (d *Distributor) RegisterFor(messageType Message, handler MessageHandlerFunc) {

	name := messageName(messageType)

	if _, found := d.topics[name]; !found {
		d.topics[name] = []MessageHandlerFunc{}
	}

	d.topics[name] = append(d.topics[name], handler)
}

func (d *Distributor) Dispatch(message Message) func() {

	name := messageName(message)
	listeners, found := d.topics[name]

	if !found {
		return func() {}
	}

	wg := sync.WaitGroup{}

	for _, handler := range listeners {
		wg.Add(1)

		go func(h MessageHandlerFunc) {
			defer wg.Done()
			h(message)
		}(handler)
	}

	return wg.Wait
}

func messageName(event interface{}) string {
	t := reflect.TypeOf(event)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name()
}
