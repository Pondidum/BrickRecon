package ldraw

import (
	"bufio"
	"fmt"
	"strings"
)

type Model struct {
	Index  int
	Name   string
	Parts  []*Part
	Models []*ModelReference
}

type Part struct {
	ColorId string
	File    string
}

type ModelReference struct {
	PrimaryColorId string
	Name           string
}

const metaHeaderName = "0 Name:"
const metaHeaderModelEnd = "0 NOFILE"
const partPrefix = "1 "

func CreateBrickList(reader *bufio.Scanner) ([]*Brick, error) {

	models, err := parseFile(reader)
	if err != nil {
		return nil, err
	}

	return collectBricks(models)
}

func parseFile(reader *bufio.Scanner) (map[string]*Model, error) {

	models := map[string]*Model{}
	var model *Model

	for reader.Scan() {
		if reader.Err() != nil {
			return nil, reader.Err()
		}

		line := reader.Text()

		if strings.HasPrefix(line, metaHeaderName) {
			model = &Model{
				Index:  len(models),
				Name:   strings.ToLower(strings.TrimSpace(strings.TrimPrefix(line, metaHeaderName))),
				Models: []*ModelReference{},
				Parts:  []*Part{},
			}
		}

		if strings.HasPrefix(line, metaHeaderModelEnd) {
			models[model.Name] = model
		}

		if strings.HasPrefix(line, partPrefix) {

			if err := parsePartLine(model, line); err != nil {
				return nil, err
			}
		}
	}

	// single model files don't end with `metaHeaderModelEnd`
	if model != nil {
		models[model.Name] = model
	}

	return models, nil
}

const (
	partLineType = iota
	partColour
	partX
	partY
	partZ
	partMatrixA
	partMatrixB
	partMatrixC
	partMatrixD
	partMatrixE
	partMatrixF
	partMatrixG
	partMatrixH
	partMatrixI
	partFile
)

func parsePartLine(model *Model, line string) error {

	cells := strings.Split(strings.TrimSpace(line), " ")

	colorId := cells[partColour]

	isPart := strings.HasSuffix(line, ".dat")
	filename := strings.TrimSuffix(strings.Join(cells[partFile:], " "), ".dat")

	if isPart {
		model.Parts = append(model.Parts, &Part{
			ColorId: colorId,
			File:    filename,
		})
	} else {
		model.Models = append(model.Models, &ModelReference{
			PrimaryColorId: colorId,
			Name:           filename,
		})
	}

	return nil
}

type Brick struct {
	LDrawID  string
	ColorId  string
	Quantity int
}

func addBrick(all map[string]*Brick, part *Part) {

	key := fmt.Sprintf("%s|%s", part.File, part.ColorId)

	if b, found := all[key]; found {
		b.Quantity++
	} else {
		all[key] = &Brick{
			LDrawID:  part.File,
			ColorId:  part.ColorId,
			Quantity: 1,
		}
	}
}

func collectBricks(models map[string]*Model) ([]*Brick, error) {

	var primary *Model
	for _, m := range models {
		if m.Index == 0 {
			primary = m
			break
		}
	}

	bricks := map[string]*Brick{}

	if err := addModelBricks(bricks, models, primary); err != nil {
		return nil, err
	}

	list := []*Brick{}

	for _, b := range bricks {
		list = append(list, b)
	}

	return list, nil
}

func addModelBricks(bricks map[string]*Brick, models map[string]*Model, model *Model) error {

	for _, p := range model.Parts {
		addBrick(bricks, p)
	}

	for _, mr := range model.Models {

		nextModel, found := models[mr.Name]

		if !found {
			return fmt.Errorf("No model called %s found", mr.Name)
		}

		if err := addModelBricks(bricks, models, nextModel); err != nil {
			return err
		}
	}

	return nil
}
