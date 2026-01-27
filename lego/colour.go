package lego

import (
	"encoding/json"
	"fmt"
	"strings"

	_ "embed"
)

type ColorId string

//go:embed colours.json
var colourJson []byte

var brickowlIndex map[string]*colourDto
var bricklinkIndex map[string]*colourDto
var ldrawIndex map[string]*colourDto
var officialIndex map[ColorId]*colourDto

func init() {

	lookup := map[string]*colourDto{}
	if err := json.Unmarshal(colourJson, &lookup); err != nil {
		panic(err)
	}

	brickowlIndex = make(map[string]*colourDto, len(lookup)*2)
	bricklinkIndex = make(map[string]*colourDto, len(lookup)*2)
	ldrawIndex = make(map[string]*colourDto, len(lookup)*2)
	officialIndex = make(map[ColorId]*colourDto, len(lookup)*2)

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

func GetColourHex(id ColorId) string {

	colour, found := officialIndex[id]
	if !found {
		return ""
	}

	return colour.Hex
}

type colourDto struct {
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
