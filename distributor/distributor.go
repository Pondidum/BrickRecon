package distributor

import (
	"context"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/honeycombio/beeline-go"
)

type MessageHandlerFunc func(ctx context.Context, message Message)

type messageHandler struct {
	action MessageHandlerFunc
	name   string
}

type Distributor struct {
	topics map[string][]messageHandler
}

func NewDistributor() *Distributor {
	return &Distributor{
		topics: map[string][]messageHandler{},
	}
}

func (d *Distributor) RegisterFor(messageType Message, handler MessageHandlerFunc) {

	name := messageName(messageType)

	if _, found := d.topics[name]; !found {
		d.topics[name] = []messageHandler{}
	}

	d.topics[name] = append(d.topics[name], messageHandler{
		action: handler,
		name:   runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name(),
	})
}

func (d *Distributor) Dispatch(ctx context.Context, message Message) func() {

	name := messageName(message)
	key := "bus." + slither(name)

	listeners, found := d.topics[name]

	beeline.AddField(ctx, key+"_handlers_count", len(listeners))

	if !found {
		return func() {}
	}

	wg := sync.WaitGroup{}

	for _, handler := range listeners {
		wg.Add(1)

		go func(h messageHandler) {
			defer wg.Done()
			c, _ := beeline.StartSpan(ctx, h.name)
			h.action(c, message)
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

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func slither(input string) string {

	snake := matchFirstCap.ReplaceAllString(input, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
