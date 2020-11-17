package allparts

import (
	"brickrecon/eventstore"
	"brickrecon/lego"
)

var ProjectionName = "allparts"

type AllPartsProjection struct{}

func (p *AllPartsProjection) Name() string {
	return ProjectionName
}

func (p *AllPartsProjection) CreateState() interface{} {
	return NewAllPartsView()
}

func (p *AllPartsProjection) Project(state interface{}, event eventstore.Event) interface{} {
	view := state.(*AllPartsView)

	switch e := event.(type) {

	case *lego.PartCreated:
		view.KnownParts[e.Key] = true

	case *lego.PartImageAdded:
		key := lego.PartKey(e.AggregateRootID)
		view.HasImage[key] = true
	}

	return view
}
