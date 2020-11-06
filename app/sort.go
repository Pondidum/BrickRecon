package app

import (
	"brickrecon/preen"
	"sort"
)

type sorter struct {
	kit string
}

func createSorter(pc *preen.PreenContext) *sorter {
	return &sorter{
		kit: pc.QueryValue("kit"),
	}
}

func (s *sorter) Parts(parts []*PartWithKitPart) []*PartWithKitPart {

	if s.kit != "" {
		return s.byKitAddition(parts)
	}

	return s.byPartThenColour(parts)
}

func (s *sorter) clone(parts []*PartWithKitPart) []*PartWithKitPart {
	result := make([]*PartWithKitPart, len(parts))

	for i, p := range parts {
		result[i] = p
	}

	return result
}

func (s *sorter) byKitAddition(parts []*PartWithKitPart) []*PartWithKitPart {

	result := s.clone(parts)

	sort.Slice(result, func(x int, y int) bool {
		l := result[x]
		r := result[y]

		if (l.KitQuantity > 0) == (r.KitQuantity > 0) {

			if l.ID == r.ID {
				return l.ColourID < r.ColourID
			}

			return l.ID < r.ID
		}

		return l.KitQuantity > 0
	})

	return result
}

func (s *sorter) byPartThenColour(parts []*PartWithKitPart) []*PartWithKitPart {

	result := s.clone(parts)

	sort.Slice(result, func(x int, y int) bool {
		l := result[x]
		r := result[y]

		if l.ID == r.ID {
			return l.ColourID < r.ColourID
		}

		return l.ID < r.ID
	})

	return result
}
