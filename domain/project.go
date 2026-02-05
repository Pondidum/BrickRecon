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
		Parts:          map[string]*ProjectPart{},
		Stock:          Stock{},
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

	Name  string
	Parts map[string]*ProjectPart
	Stock Stock
}

type ProjectView struct {
	Name  string
	Parts map[string]*ProjectPart
	Stock Stock
}

func (pv *ProjectView) UniqueParts() int {
	return len(pv.Parts)
}
func (pv *ProjectView) TotalParts() int {
	sum := 0
	for _, part := range pv.Parts {
		sum += part.Wanted
	}
	return sum
}
func (pv *ProjectView) OwnedParts() int {
	sum := 0
	for _, colorStock := range pv.Stock {
		for _, stock := range colorStock {
			sum += stock
		}
	}
	return sum
}

type ProjectPart struct {
	Number lego.PartId
	Color  lego.ColorId
	Wanted int
}

func newProjectPart(part *lego.InventoryPart) *ProjectPart {
	return &ProjectPart{
		Number: part.Id,
		Color:  part.ColorId,
		Wanted: part.Quantity,
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
	p.Name = e.Name
}

func (p *Project) ImportPartsList(parts []*lego.InventoryPart, source string) {
	goes.Apply(p.AggregateState, &PartsImported{
		Parts:  parts,
		Source: source,
	})
}

func (p *Project) onPartsImported(e PartsImported) {
	//for now, just wipe out everything and replace
	p.Parts = make(map[string]*ProjectPart, len(e.Parts))

	for _, part := range e.Parts {
		p.Parts[keyFor(part.Id, part.ColorId)] = newProjectPart(part)
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

		if part, found := p.Parts[key]; found {
			part.Wanted += add.Quantity
		} else {
			p.Parts[key] = newProjectPart(add)
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

		if part, found := p.Parts[key]; found {
			part.Wanted -= rem.Quantity
			if part.Wanted <= 0 {
				delete(p.Parts, key)
			}
		}
	}
}

type PartsRemoved struct {
	Parts []*lego.InventoryPart
}

func (p *Project) AddStock(stock Stock) error {
	return goes.Apply(p.AggregateState, &StockAdded{
		Added: stock,
	})
}

func (p *Project) onStockAdded(e StockAdded) {
	for part, colors := range e.Added {
		for color, quantity := range colors {
			AddStock(p.Stock, part, color, quantity)
		}
	}
}

type StockAdded struct {
	Added Stock
}

func (p *Project) RemoveStock(stock Stock) error {
	return goes.Apply(p.AggregateState, &StockRemoved{
		Removed: stock,
	})
}

func (p *Project) onStockRemoved(e StockRemoved) {

	for part, colors := range e.Removed {
		for color, quantity := range colors {
			RemoveStock(p.Stock, part, color, quantity)
		}
	}
}

type StockRemoved struct {
	Removed Stock
}
