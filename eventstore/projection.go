package eventstore

import (
	"context"
)

type Projection interface {
	Name() string
	CreateState() interface{}
	Project(state interface{}, event Event) interface{}
}

type Projector struct {
	projections map[string]Projection
	backend     Backend
}

func NewProjector(backend Backend) *Projector {
	return &Projector{
		projections: map[string]Projection{},
		backend:     backend,
	}
}

func (p *Projector) registerProjection(ctx context.Context, projection Projection) {
	p.projections[projection.Name()] = projection
}

func (p *Projector) runAllProjections(ctx context.Context, events []Event) error {
	var err error

	for _, projection := range p.projections {
		err = p.runProjection(ctx, projection, events)
	}

	return err
}

func (p *Projector) runProjection(ctx context.Context, projection Projection, events []Event) error {

	view := p.backend.NewView(projection.Name())
	state := projection.CreateState()

	if err := view.ReadView(ctx, state); err != nil {
		return err
	}

	for _, e := range events {
		state = projection.Project(state, e)
	}

	return view.WriteView(ctx, state, 0)
}
