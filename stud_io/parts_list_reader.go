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

type ListPart struct {
	BrickLinkID     lego.BrickLinkPart
	ElementID       string
	LDrawID         lego.LDrawPart
	Name            lego.PartName
	BrickLinkColour lego.BrickLinkColour
	LDrawColour     lego.LDrawColour
	ColourName      lego.ColourName
	ColourCategory  string
	Quantity        int
	Weight          float64
}

func ReadPartsList(content io.Reader) ([]*ListPart, error) {

	reader := csv.NewReader(content)
	reader.Comma = '\t'

	parts := []*ListPart{}

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

func parsePart(fields []string) (*ListPart, error) {

	var err error

	part := &ListPart{
		Name: lego.PartName(fields[partName]),
	}

	parsePartAliases(part, fields)

	if err = parseColour(part, fields); err != nil {
		return nil, err
	}

	if part.Quantity, err = strconv.Atoi(fields[quantity]); err != nil {
		return nil, convertError("part.Quantity", fields[quantity])
	}

	if part.Weight, err = strconv.ParseFloat(fields[weight], 64); err != nil {
		return nil, convertError("part.Weight", fields[weight])
	}

	return part, err
}

func parsePartAliases(part *ListPart, fields []string) {
	part.BrickLinkID = lego.BrickLinkPart(fields[brickLinkID])
	part.LDrawID = lego.LDrawPart(fields[ldrawID])

}

func parseColour(part *ListPart, fields []string) error {

	var err error

	var bricklink, ldraw int

	if bricklink, err = strconv.Atoi(fields[brickLinkColour]); err != nil {
		return convertError("colour.BrickLinkID", fields[brickLinkColour])
	}

	if ldraw, err = strconv.Atoi(fields[ldrawColour]); err != nil {
		return convertError("colour.LDrawID", fields[ldrawColour])
	}

	part.BrickLinkColour = lego.BrickLinkColour(bricklink)
	part.LDrawColour = lego.LDrawColour(ldraw)
	part.ColourName = lego.ColourName(fields[colourName])
	part.ColourCategory = fields[colourCategory]

	return nil
}

func convertError(key string, value string) error {
	return fmt.Errorf("Unable to convert '%s' to %s", value, key)
}
