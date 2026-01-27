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
	GetPartName(id lego.PartId) lego.PartName
}

func ReadParts(file io.ReaderAt, fileSize int64, lookup PartLookup) ([]*lego.InventoryPart, error) {

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

	parts := make([]*lego.InventoryPart, len(bricks))

	for i, brick := range bricks {
		part, err := toPart(brick, lookup)
		if err != nil {
			return nil, err
		}
		parts[i] = part
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

func toPart(brick *ldraw.Brick, lookup PartLookup) (*lego.InventoryPart, error) {

	id := lego.PartId(brick.LDrawID)
	name := lookup.GetPartName(id)

	color, err := lego.GetColorId(brick.ColorId, "ldraw")
	if err != nil {
		return nil, err
	}

	return &lego.InventoryPart{
		Part: lego.Part{
			Id:   id,
			Name: name,
		},
		ColorId:  color,
		Quantity: brick.Quantity,
	}, nil
}
