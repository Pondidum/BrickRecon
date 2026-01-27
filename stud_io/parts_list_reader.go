package stud_io

import (
	"brickrecon/lego"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	brickLinkID int = iota
	elementID
	ldrawID
	partName
	brickLinkColour
	ldrawColour
	colourName
	colourCategory
	quantity
	weight
)

func ReadPartsList(content io.Reader) ([]*lego.InventoryPart, error) {

	reader := csv.NewReader(content)
	reader.Comma = '\t'

	parts := []*lego.InventoryPart{}

	// read the header
	_, err := reader.Read()
	if err != nil {
		return nil, err
	}

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if isBlank(record) {
			continue
		}

		if isSummaryHeader(record) {
			break
		}

		part, err := parsePart(record)
		if err != nil {
			return nil, fmt.Errorf("%s\nRecord: %s", err.Error(), record)
		}

		parts = append(parts, part)
	}

	return parts, nil
}

func isBlank(fields []string) bool {
	return len(fields) == 0 || strings.TrimSpace(fields[0]) == ""
}

func isSummaryHeader(fields []string) bool {
	return len(fields) > 0 && fields[brickLinkID] == "Total qty"
}

func parsePart(fields []string) (*lego.InventoryPart, error) {
	colorId, err := lego.GetColorId(fields[ldrawColour], "ldraw")
	if err != nil {
		return nil, err
	}

	quantity, err := strconv.Atoi(fields[quantity])
	if err != nil {
		return nil, fmt.Errorf("Unable to convert quantity to '%s'", fields[quantity])
	}

	return &lego.InventoryPart{
		Part: lego.Part{
			Id:   lego.PartId(fields[ldrawID]),
			Name: lego.PartName(fields[partName]),
		},

		ColourId: colorId,
		Quantity: quantity,
	}, nil
}
