package background

import (
	"mvc/distributor"
	"mvc/eventstore"
	"mvc/lego"

	uuid "github.com/satori/go.uuid"
)

type PartsAddedMessage struct {
	distributor.Message

	Parts []lego.Part
}

func AttachImageCacheListener(bus *distributor.Distributor, es *eventstore.EventStore) error {

	if _, err := loadCache(es); err != nil {
		return err
	}

	bus.RegisterFor(&PartsAddedMessage{}, func(message distributor.Message) {
		handler(es, message)
	})

	return nil
}

var cacheID uuid.UUID = uuid.Must(uuid.FromString("b83e7c15-24d7-4f18-8de7-34de416eb9de"))

func loadCache(es *eventstore.EventStore) (*ImageCache, error) {

	ic := blankImageCache()

	err := es.LoadAggregate(cacheID, ic)

	if err == nil {
		return ic, nil
	}

	if !eventstore.IsAggregateNotFound(err) {
		return nil, err
	}

	ic = NewImageCache(cacheID)

	if err = es.SaveAggregate(ic); err != nil {
		return nil, err
	}

	return ic, nil
}

func handler(es *eventstore.EventStore, message distributor.Message) {

	m, ok := message.(*PartsAddedMessage)

	if !ok {
		return
	}

	ic, err := loadCache(es)

	if err != nil {
		return
	}

	for _, part := range m.Parts {
		ic.AddPart(part)
	}

	es.SaveAggregate(ic)

	ic.Run()

	es.SaveAggregate(ic)
}
