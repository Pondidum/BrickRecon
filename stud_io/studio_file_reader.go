package stud_io

import (
	"brickrecon/ldraw"
	"brickrecon/lego"
	"bufio"
	"fmt"
	"io"

	"github.com/yeka/zip"
)

type PartLookup interface {
	GetPartName(id lego.LDrawPart) lego.PartName
	GetColour(colour lego.LDrawColour) lego.BrickLinkColour
}

type StudioPart struct {
	BrickLinkID     lego.BrickLinkPart
	LDrawID         lego.LDrawPart
	Name            lego.PartName
	BrickLinkColour lego.BrickLinkColour
	LDrawColour     lego.LDrawColour
	// ColourName      lego.ColourName
	// ColourCategory  string
	Quantity int
	// Weight          float64
}

func ReadParts(file io.ReaderAt, fileSize int64, lookup PartLookup) ([]*StudioPart, error) {

	reader, err := zip.NewReader(file, fileSize)
	if err != nil {
		return nil, err
	}

	model, found := findModel(reader.File)
	if !found {
		return nil, fmt.Errorf("Unable to find model.ldr in the stud.io file")
	}

	model.SetPassword("soho0909")
	r, err := model.Open()
	if err != nil {
		return nil, err
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	bricks, err := ldraw.CreateBrickList(scanner)
	if err != nil {
		return nil, err
	}

	parts := make([]*StudioPart, len(bricks))

	for i, brick := range bricks {
		parts[i] = toPart(brick, lookup)
	}

	return parts, nil
}

func findModel(all []*zip.File) (*zip.File, bool) {

	for _, f := range all {
		if f.Name == "model.ldr" {
			return f, true
		}
	}

	return nil, false
}

func toPart(brick *ldraw.Brick, lookup PartLookup) *StudioPart {

	id := lego.LDrawPart(brick.LDrawID)
	name := lookup.GetPartName(id)

	ldColour := lego.LDrawColour(brick.Colour)
	blColour := lookup.GetColour(ldColour)

	return &StudioPart{
		LDrawID:         id,
		BrickLinkID:     lego.BrickLinkPart(brick.LDrawID),
		Name:            name,
		LDrawColour:     ldColour,
		BrickLinkColour: blColour,
		Quantity:        brick.Quantity,
	}
}
