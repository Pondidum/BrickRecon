package lego

import "encoding/json"

type Colour struct {
	ID      int
	Aliases ColourAliases

	Name     string
	Category string
}

type ColourAliases struct {
	BrickLinkID int
	LDrawID     int
	Boid        int
}

type PartID struct {
	id string
}

func NewPartID(value string) PartID {
	return PartID{id: value}
}

func (p PartID) String() string {
	return p.id
}

func (p *PartID) UnmarshalJSON(s []byte) (err error) {
	return json.Unmarshal(s, &p.id)
}

func (p PartID) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.id)
}

func (p *PartID) UnmarshalText(text []byte) (err error) {
	p.id = string(text)
	return nil
}

type Part struct {
	ID      PartID
	Aliases PartAliases

	Name   string
	Colour Colour

	Quantity int
	Weight   float64
}

type PartAliases struct {
	BrickLinkID string
	ElementID   int
	LDrawID     string
	Boid        string
}
