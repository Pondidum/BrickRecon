package lego

import (
	"fmt"
	"strconv"
	"strings"
)

type PartKey string

func CreatePartKey(part LDrawPart, colour LDrawColour) PartKey {
	return PartKey(fmt.Sprintf("%v|%v", part, colour))
}

func ParsePartKey(key PartKey) (LDrawPart, LDrawColour) {
	segments := strings.Split(string(key), "|")
	val, _ := strconv.Atoi(segments[1])

	return LDrawPart(segments[0]), LDrawColour(val)
}
