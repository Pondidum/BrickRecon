package distributor

import "reflect"

type MessageMeta struct {
}

func (m *MessageMeta) meta() *MessageMeta { return m }

type Message interface {
	meta() *MessageMeta
}

// -------------

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

func (d *Distributor) Dispatch(message Message) {

	name := messageName(message)
	listeners, found := d.topics[name]

	if !found {
		return
	}

	for _, handler := range listeners {
		go handler(message)
	}

	return
}

func messageName(event interface{}) string {
	t := reflect.TypeOf(event)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name()
}
