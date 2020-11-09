package all_kits

import (
	"brickrecon/eventstore"
	"brickrecon/lego"
	"fmt"
)

var ProjectionName string = "kits"

type KitsProjection struct{}

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
		view.Kits[e.KitNumber] = &KitView{
			ID:     e.AggregateRootID,
			Name:   e.KitName,
			Number: e.KitNumber,
			Parts:  toPartView(e.Parts),
		}
	}

	return view
}

func toPartView(parts []*lego.Part) []PartView {

	views := make([]PartView, len(parts))

	for i, part := range parts {
		views[i] = PartView{
			Key:        part.Key,
			ID:         part.Aliases.LDrawID,
			Name:       part.Name,
			ColourID:   part.Colour.ID,
			ColourName: part.Colour.Name,
			ColourHex:  part.Colour.Hex,
			ImagePath:  fmt.Sprintf("%s-%v.png", part.Aliases.LDrawID, part.Colour.ID),
			Quantity:   part.Quantity,
		}
	}

	return views
}
