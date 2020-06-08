package background

import (
	"mvc/distributor"
	"mvc/eventstore"
	"mvc/lego"
	"os"

	uuid "github.com/satori/go.uuid"
)

type PartsAddedMessage struct {
	distributor.Message

	Parts []lego.Part
}

type PartsAddedMessageHandler struct {
	es *eventstore.EventStore
}

func NewPartsAddedMessageHandler(es *eventstore.EventStore) (*PartsAddedMessageHandler, error) {

	handler := &PartsAddedMessageHandler{es}

	if _, err := handler.loadCache(); err != nil {
		return nil, err
	}

	return handler, nil
}

func (h *PartsAddedMessageHandler) loadCache() (*ImageCache, error) {
	id, _ := uuid.FromString("b83e7c15-24d7-4f18-8de7-34de416eb9de")
	ic := NewImageCache()

	if err := h.es.LoadAggregate(id, ic); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return ic, nil
}

func (h *PartsAddedMessageHandler) Handle(message distributor.Message) {

	m, ok := message.(PartsAddedMessage)

	if !ok {
		return
	}

	ic, err := h.loadCache()

	if err != nil {
		return
	}

	for _, part := range m.Parts {
		ic.AddPart(&part)
	}

	h.es.SaveAggregate(ic)

	ic.Run()

	h.es.SaveAggregate(ic)
}
