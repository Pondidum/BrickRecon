package lego

type Project struct {
	Name string

	parts []*Part
}

func NewProject(name string, parts []Part) Project {
	m := Project{
		Name:  name,
		parts: make([]*Part, len(parts)),
	}

	for i, p := range parts {
		m.parts[i] = &p
	}

	return m
}

func (m *Project) AddPart(part Part) {

	id := part.BrickLinkID
	colour := part.Colour.BrickLinkID

	existing, found := m.partByTypeAndColour(id, colour)

	if found {
		existing.Quantity += part.Quantity
		return
	}

	m.parts = append(m.parts, &part)
}

func (m *Project) partByTypeAndColour(brickLinkID string, colourID int) (*Part, bool) {

	for _, p := range m.parts {

		if p.BrickLinkID == brickLinkID && p.Colour.BrickLinkID == colourID {
			return p, true
		}
	}

	return nil, false
}
