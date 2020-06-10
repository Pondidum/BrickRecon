package background

import (
	"context"
	"mvc/distributor"
	"mvc/eventstore"
	"mvc/lego"

	"github.com/honeycombio/beeline-go"
	uuid "github.com/satori/go.uuid"
)

type PartsAddedMessage struct {
	distributor.Message

	Parts []lego.Part
}

func AttachImageCacheListener(bus *distributor.Distributor, es *eventstore.EventStore) error {

	ic, err := loadCache(es)
	if err != nil {
		return err
	}

	if err := ic.ReadFromCache(); err != nil {
		return err
	}

	if err := es.SaveAggregate(ic); err != nil {
		return err
	}

	bus.RegisterFor(&PartsAddedMessage{}, func(ctx context.Context, message distributor.Message) {
		handler(es, ctx, message)
	})

	return nil
}

var cacheID uuid.UUID = uuid.Must(uuid.FromString("b83e7c15-24d7-4f18-8de7-34de416eb9de"))

func loadCache(es *eventstore.EventStore) (*ImageCache, error) {

	ic := blankImageCache("./app/static/img/parts")

	err := es.LoadAggregate(cacheID, ic)

	if err == nil {
		return ic, nil
	}

	if !eventstore.IsAggregateNotFound(err) {
		return nil, err
	}

	ic = NewImageCache(cacheID, "./app/static/img/parts")

	if err = es.SaveAggregate(ic); err != nil {
		return nil, err
	}

	return ic, nil
}

func handler(es *eventstore.EventStore, ctx context.Context, message distributor.Message) {

	m, ok := message.(*PartsAddedMessage)

	if !ok {
		return
	}

	ic, err := loadCache(es)

	if err != nil {
		beeline.AddField(ctx, "error_loading_cache", err)
		return
	}

	for _, part := range m.Parts {
		ic.AddPart(part)
	}

	if err := es.SaveAggregate(ic); err != nil {
		beeline.AddField(ctx, "error_saving_cache", err)
		return
	}

	// later, move this to a separate process to run periodically

	ic.Run()

	if err := es.SaveAggregate(ic); err != nil {
		beeline.AddField(ctx, "error_saving_processed_cache", err)
	}
}
