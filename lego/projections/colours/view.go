package colours

import "brickrecon/lego"

type ColoursView struct {
	ByBrickLink map[lego.BrickLinkColour]*ColourView
	ByLDraw     map[lego.LDrawColour]*ColourView
}

type ColourView struct {
	BrickLinkID lego.BrickLinkColour
	LDrawID     lego.LDrawColour
	Name        lego.ColourName
	Hex         lego.HexColour
	Category    string
}
