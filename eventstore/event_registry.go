package eventstore

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/honeycombio/beeline-go"
)

type EventRegistry struct {
	registry map[string]Initialiser
}

func NewRegistry() *EventRegistry {
	er := &EventRegistry{
		registry: map[string]Initialiser{},
	}

	return er
}

func (er *EventRegistry) Register(ctx context.Context, creator Initialiser) error {

	event := creator()
	v := reflect.ValueOf(event)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("Event initialiser must return a pointer to a struct")
	}

	er.registry[eventName(event)] = creator

	return nil
}

func (er *EventRegistry) RegisterMany(ctx context.Context, events []Initialiser) error {
	for _, ei := range events {
		if err := er.Register(ctx, ei); err != nil {
			return err
		}
	}

	return nil
}

func (er *EventRegistry) CreateInstance(ctx context.Context, eventType string) (interface{}, error) {

	creator, found := er.registry[eventType]

	if !found {
		beeline.AddField(ctx, "es.event_type_lookup_failed", true)
		return nil, fmt.Errorf("Unable to find an event of type %s", eventType)
	}

	event := creator()

	return event, nil
}
