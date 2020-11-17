package allparts

import "brickrecon/lego"

type AllPartsView struct {
	KnownParts map[lego.PartKey]bool
	HasImage   map[lego.PartKey]bool
}

func NewAllPartsView() *AllPartsView {
	return &AllPartsView{
		KnownParts: map[lego.PartKey]bool{},
		HasImage:   map[lego.PartKey]bool{},
	}
}
