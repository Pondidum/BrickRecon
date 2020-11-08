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

func ReadParts(file io.ReaderAt, fileSize int64, lookup PartLookup) ([]lego.Part, error) {

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

	parts := make([]lego.Part, len(bricks))

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

func toPart(brick *ldraw.Brick, lookup PartLookup) lego.Part {

	id := lego.LDrawPart(brick.LDrawID)
	name := lookup.GetPartName(id)

	ldColour := lego.LDrawColour(brick.Colour)
	blColour := lookup.GetColour(ldColour)

	return lego.Part{
		Key:  lego.CreatePartKey(id, blColour),
		Name: name,
		Colour: lego.Colour{
			ID: blColour,
			Aliases: lego.ColourAliases{
				LDrawID:     ldColour,
				BrickLinkID: blColour,
			}},
		Quantity: brick.Quantity,
	}
}
