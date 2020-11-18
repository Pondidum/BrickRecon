package all_kits

import (
	"brickrecon/eventstore"
	"brickrecon/lego"
	"context"
)

var ProjectionName string = "kits"

func NewKitsProjection(es eventstore.EventStore) *KitsProjection {
	return &KitsProjection{
		partLoader: func(key lego.PartKey) *lego.PartAggregate {
			part := lego.BlankPart()
			es.LoadAggregate(context.Background(), eventstore.AggregateID(key), part)
			return part
		},
	}
}

type KitsProjection struct {
	partLoader PartLoader
}

func (p *KitsProjection) Name() string {
	return ProjectionName
}

func (p *KitsProjection) CreateState() interface{} {
	return &AllKitsView{
		Kits: map[lego.KitNumber]*KitView{},
	}
}

func (p *KitsProjection) Project(state interface{}, event eventstore.Event) interface{} {
	view := state.(*AllKitsView)

	switch e := event.(type) {

	case *lego.KitCreated:

		parts := []*PartView{}
		for key, quantity := range e.Parts {
			parts = append(parts, newPartView(p.partLoader, key, quantity))
		}

		view.Kits[e.KitNumber] = &KitView{
			ID:     e.AggregateRootID,
			Name:   e.KitName,
			Number: e.KitNumber,
			Parts:  parts,
		}

	}

	return view
}
