package colours

import (
	"brickrecon/eventstore"
	"brickrecon/lego"
)

var ProjectionName = "colours"

type ColoursProjection struct{}

func (p *ColoursProjection) Name() string {
	return ProjectionName
}

func (p *ColoursProjection) CreateState() interface{} {
	return &ColoursView{
		ByBrickLink: map[lego.BrickLinkColour]*ColourView{},
		ByLDraw:     map[lego.LDrawColour]*ColourView{},
	}
}

func (p *ColoursProjection) Project(state interface{}, event eventstore.Event) interface{} {
	view := state.(*ColoursView)

	switch e := event.(type) {

	case *lego.ProjectPartsAdded:
		for _, p := range e.Parts {
			view.addPart(p)
		}
	}

	return view
}

func (v *ColoursView) addPart(p *lego.Part) {
	colour := &ColourView{
		Name:        p.Colour.Name,
		LDrawID:     p.Colour.Aliases.LDrawID,
		BrickLinkID: p.Colour.Aliases.BrickLinkID,
		Category:    p.Colour.Category,
		Hex:         p.Colour.Hex,
	}

	v.ByBrickLink[colour.BrickLinkID] = colour
	v.ByLDraw[colour.LDrawID] = colour
}
