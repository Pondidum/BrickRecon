package allparts

import (
	"brickrecon/eventstore"
	"brickrecon/lego"
	"context"
)

var ProjectionName = "allparts"

type AllPartsProjection struct{}

func (p *AllPartsProjection) Name() string {
	return ProjectionName
}

func (p *AllPartsProjection) CreateState() interface{} {
	return NewAllPartsView()
}

func (p *AllPartsProjection) Project(ctx context.Context, state interface{}, event eventstore.Event) interface{} {
	view := state.(*AllPartsView)

	switch e := event.(type) {

	case *lego.PartCreated:
		view.KnownParts[e.Key] = true

	case *lego.PartImageAdded:
		key := lego.PartKey(e.AggregateRootID)
		view.HasImage[key] = true

	case *lego.PartNamesAdded:
		key := lego.PartKey(e.AggregateRootID)
		partID, _ := lego.ParsePartKey(key)

		view.Names[partID] = e.PartName
	}

	return view
}
