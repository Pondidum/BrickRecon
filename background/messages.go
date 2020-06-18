package background

import (
	"brickrecon/distributor"
	"brickrecon/eventstore"
	"brickrecon/lego"
	"context"

	"github.com/honeycombio/beeline-go"
)

type PartsAddedMessage struct {
	distributor.Message

	Parts []lego.Part
}

func ImageCacheHandler(es eventstore.EventStore) func(context.Context, distributor.Message) {
	return func(ctx context.Context, message distributor.Message) {
		handler(es, ctx, message)
	}
}

func handler(es eventstore.EventStore, ctx context.Context, message distributor.Message) {

	m, ok := message.(*PartsAddedMessage)

	if !ok {
		return
	}

	ic, err := NewImageCache(es, "./app/static/img/parts", ctx)

	if err != nil {
		beeline.AddField(ctx, "error_loading_cache", err)
		return
	}

	for _, part := range m.Parts {
		ic.AddPart(part)
	}

	if err := es.SaveAggregate(ctx, ic); err != nil {
		beeline.AddField(ctx, "error_saving_cache", err)
		return
	}

	// later, move this to a separate process to run periodically

	ic.Run(ctx)

	if err := es.SaveAggregate(ctx, ic); err != nil {
		beeline.AddField(ctx, "error_saving_processed_cache", err)
	}
}
