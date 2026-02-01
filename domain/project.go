package domain

import (
	"brickrecon/goes"
	"brickrecon/lego"
	"fmt"
	"github.com/google/uuid"
)

func BlankProject() *Project {
	p := &Project{
		AggregateState: goes.NewAggregateState(),
		parts:          map[string]*ProjectPart{},
	}

	goes.Register(p.AggregateState, p.onProjectCreated)
	goes.Register(p.AggregateState, p.onPartsImported)
	goes.Register(p.AggregateState, p.onPartsAdded)
	goes.Register(p.AggregateState, p.onPartsRemoved)
	goes.Register(p.AggregateState, p.onStockAdded)
	goes.Register(p.AggregateState, p.onStockRemoved)

	return p
}

func CreateProject(name string) (*Project, error) {
	project := BlankProject()
	err := goes.Apply(project, ProjectCreated{
		ID:   uuid.New(),
		Name: name,
	})
	return project, err
}

type Project struct {
	*goes.AggregateState

	parts map[string]*ProjectPart
}

type ProjectPart struct {
	Number lego.PartId
	Color  lego.ColorId
	Wanted int
	Stock  map[lego.ColorId]int
}

func newProjectPart(part *lego.InventoryPart) *ProjectPart {
	return &ProjectPart{
		Number: part.Id,
		Color:  part.ColorId,
		Wanted: part.Quantity,
		Stock:  map[lego.ColorId]int{},
	}
}

type ProjectCreated struct {
	ID   uuid.UUID
	Name string
}

func keyFor(p lego.PartId, c lego.ColorId) string {
	return fmt.Sprintf("%s|%s", p, c)
}

func (p *Project) onProjectCreated(e ProjectCreated) {
	goes.SetID(p.AggregateState, e.ID)
}

func (p *Project) ImportPartsList(parts []*lego.InventoryPart, source string) {
	goes.Apply(p.AggregateState, &PartsImported{
		Parts:  parts,
		Source: source,
	})
}

func (p *Project) onPartsImported(e PartsImported) {
	//for now, just wipe out everything and replace
	p.parts = make(map[string]*ProjectPart, len(e.Parts))

	for _, part := range e.Parts {
		p.parts[keyFor(part.Id, part.ColorId)] = newProjectPart(part)
	}
}

type PartsImported struct {
	Parts  []*lego.InventoryPart
	Source string
}

func (p *Project) AddParts(parts []*lego.InventoryPart) error {
	return goes.Apply(p.AggregateState, &PartsAdded{Parts: parts})
}

func (p *Project) onPartsAdded(e PartsAdded) {
	for _, add := range e.Parts {
		key := keyFor(add.Id, add.ColorId)

		if part, found := p.parts[key]; found {
			part.Wanted += add.Quantity
		} else {
			p.parts[key] = newProjectPart(add)
		}
	}
}

type PartsAdded struct {
	Parts []*lego.InventoryPart
}

func (p *Project) RemoveParts(parts []*lego.InventoryPart) error {
	return goes.Apply(p.AggregateState, &PartsRemoved{Parts: parts})
}

func (p *Project) onPartsRemoved(e PartsRemoved) {
	for _, rem := range e.Parts {
		key := keyFor(rem.Id, rem.ColorId)

		if part, found := p.parts[key]; found {
			part.Wanted -= rem.Quantity
			if part.Wanted <= 0 {
				delete(p.parts, key)
			}
		}
	}
}

type PartsRemoved struct {
	Parts []*lego.InventoryPart
}

func (p *Project) AddStock(part lego.PartId, color lego.ColorId, quantity int) error {
	return goes.Apply(p.AggregateState, &StockAdded{
		Part:     part,
		Color:    color,
		Quantity: quantity,
	})
}

func (p *Project) onStockAdded(e StockAdded) {
	part := p.parts[keyFor(e.Part, e.Color)]
	part.Stock[part.Color] = part.Stock[part.Color] + e.Quantity
}

type StockAdded struct {
	Part     lego.PartId
	Color    lego.ColorId
	Quantity int
}

func (p *Project) RemoveStock(part lego.PartId, color lego.ColorId, quantity int) error {
	return goes.Apply(p.AggregateState, &StockRemoved{
		Part:     part,
		Color:    color,
		Quantity: quantity,
	})
}

func (p *Project) onStockRemoved(e StockRemoved) {
	part := p.parts[keyFor(e.Part, e.Color)]
	part.Stock[part.Color] = max(part.Stock[part.Color]-e.Quantity, 0)
}

type StockRemoved struct {
	Part     lego.PartId
	Color    lego.ColorId
	Quantity int
}
