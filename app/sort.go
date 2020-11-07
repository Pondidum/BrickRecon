package app

import (
	"brickrecon/preen"
	"sort"
)

type sorter struct {
	kit string

	sortKey string
}

func createSorter(pc *preen.PreenContext) *sorter {
	return &sorter{
		kit:     pc.QueryValue("kit"),
		sortKey: pc.QueryValue("sort"),
	}
}

func (s *sorter) Parts(parts []*PartWithKitPart) []*PartWithKitPart {

	if s.kit != "" {
		return s.byKitAddition(parts)
	}

	if s.sortKey == "name" {
		return s.byName(parts)
	}

	if s.sortKey == "colour" {
		return s.byColour(parts)
	}

	return s.byName(parts)

}

func (s *sorter) clone(parts []*PartWithKitPart) []*PartWithKitPart {
	result := make([]*PartWithKitPart, len(parts))

	for i, p := range parts {
		result[i] = p
	}

	return result
}

func (s *sorter) byName(parts []*PartWithKitPart) []*PartWithKitPart {
	result := s.clone(parts)

	sort.Slice(result, func(x, y int) bool {
		l := result[x]
		r := result[y]

		if l.Name == r.Name {
			return l.ColourID < r.ColourID
		}

		return l.Name < r.Name
	})

	return result
}

func (s *sorter) byColour(parts []*PartWithKitPart) []*PartWithKitPart {
	result := s.clone(parts)

	sort.Slice(result, func(x, y int) bool {
		l := result[x]
		r := result[y]

		if l.ColourID == r.ColourID {
			return l.Name < r.Name
		}

		return l.ColourID < r.ColourID

	})

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
