package lego

import (
	"encoding/json"
	"fmt"
	"strings"

	_ "embed"
)

type ColorId string

//go:embed colors.json
var colorJson []byte

var brickowlIndex map[string]*colorDto
var bricklinkIndex map[string]*colorDto
var ldrawIndex map[string]*colorDto
var officialIndex map[ColorId]*colorDto

func init() {

	lookup := map[string]*colorDto{}
	if err := json.Unmarshal(colorJson, &lookup); err != nil {
		panic(err)
	}

	brickowlIndex = make(map[string]*colorDto, len(lookup)*2)
	bricklinkIndex = make(map[string]*colorDto, len(lookup)*2)
	ldrawIndex = make(map[string]*colorDto, len(lookup)*2)
	officialIndex = make(map[ColorId]*colorDto, len(lookup)*2)

	for owl, color := range lookup {

		for _, id := range color.BrickLinkIds {
			bricklinkIndex[id] = color
		}

		for _, official := range color.OfficialData {
			officialIndex[official.LegoId] = color
		}

		for _, ldraw := range color.LDrawIds {
			ldrawIndex[ldraw] = color
		}

		brickowlIndex[owl] = color
	}
}

func GetColorId(id string, source string) (ColorId, error) {
	switch strings.ToLower(source) {
	case "brickowl":
		color := brickowlIndex[id]
		if len(color.OfficialData) == 0 {
			return "", fmt.Errorf("no official data for %s", id)
		}

		return color.OfficialData[0].LegoId, nil

	case "bricklink":
		color := bricklinkIndex[id]
		if len(color.OfficialData) == 0 {
			return "", fmt.Errorf("no official data for %s", id)
		}

		return color.OfficialData[0].LegoId, nil

	case "ldraw":
		color := ldrawIndex[id]
		if len(color.OfficialData) == 0 {
			return "", fmt.Errorf("no official data for %s", id)
		}

		return color.OfficialData[0].LegoId, nil

	default:
		return "", fmt.Errorf("unknown source %s", source)
	}
}

func GetColor(id ColorId, target string) (string, error) {
	color, found := officialIndex[id]
	if !found {
		return "", fmt.Errorf("unknown official color %s", id)
	}

	switch strings.ToLower(target) {
	case "bricklink":
		return color.BrickLinkIds[0], nil

	default:
		return "", fmt.Errorf("unknown target %s", target)
	}
}

func GetColourHex(id ColorId) string {

	color, found := officialIndex[id]
	if !found {
		return ""
	}

	return color.Hex
}

type colorDto struct {
	Id   string
	Name string
	Hex  string

	PeeronNames []string `json:"peeron_names"`

	LDrawIds       []string `json:"ldraw_ids"`
	BrickLinkIds   []string `json:"bl_ids"`
	BrickLinkNames []string `json:"bl_names"`

	OfficialData []officialColour `json:"lego_colors"`
}

type officialColour struct {
	LegoId  ColorId `json:"lego_id"`
	RawName string  `json:"raw_name"`
}
