package lego

type Model struct {
	Name string

	parts []*Part
}

func NewModel(name string, parts []Part) Model {
	m := Model{
		Name:  name,
		parts: make([]*Part, len(parts)),
	}

	for i, p := range parts {
		m.parts[i] = &p
	}

	return m
}

func (m *Model) AddPart(part Part) {

	id := part.BrickLinkID
	colour := part.Colour.BrickLinkID

	existing, found := m.partByTypeAndColour(id, colour)

	if found {
		existing.Quantity += part.Quantity
		return
	}

	m.parts = append(m.parts, &part)
}

func (m *Model) partByTypeAndColour(brickLinkID string, colourID int) (*Part, bool) {

	for _, p := range m.parts {

		if p.BrickLinkID == brickLinkID && p.Colour.BrickLinkID == colourID {
			return p, true
		}
	}

	return nil, false
}
