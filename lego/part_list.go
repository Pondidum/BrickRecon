package lego

type PartList struct {
	parts []*Part
}

func NewPartsList(parts []Part) *PartList {
	list := PartList{
		parts: make([]*Part, len(parts)),
	}

	for i, p := range parts {
		list.parts[i] = &p
	}

	return &list
}

func (m *PartList) Add(part Part) {

	id := part.ID
	colour := part.Colour.BrickLinkID

	existing, found := m.byTypeAndColour(id, colour)

	if found {
		existing.Quantity += part.Quantity
		return
	}

	m.parts = append(m.parts, &part)
}

func (m *PartList) byTypeAndColour(brickLinkID string, colourID int) (*Part, bool) {

	for _, p := range m.parts {

		if p.ID == brickLinkID && p.Colour.BrickLinkID == colourID {
			return p, true
		}
	}

	return nil, false
}
