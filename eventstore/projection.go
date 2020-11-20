package eventstore

import (
	"context"
)

type Projection interface {
	Name() string
	CreateState() interface{}
	Project(ctx context.Context, state interface{}, event Event) interface{}
}

type StatelessProjection interface {
	Name() string
	Project(event Event)
}

type statelessProjection struct {
	name    string
	project func(event Event)
}

func (p statelessProjection) Name() string        { return p.name }
func (p statelessProjection) Project(event Event) { p.project(event) }

func NewStatelessProjection(name string, project func(event Event)) StatelessProjection {
	return &statelessProjection{name: name, project: project}
}

type Projector struct {
	projections map[string]Projection
	stateless   map[string]StatelessProjection
	backend     Backend
}

func NewProjector(backend Backend) *Projector {
	return &Projector{
		projections: map[string]Projection{},
		stateless:   map[string]StatelessProjection{},
		backend:     backend,
	}
}

func (p *Projector) registerProjection(ctx context.Context, projection Projection) {
	p.projections[projection.Name()] = projection
}

func (p *Projector) registerStatelessProjection(ctx context.Context, projection StatelessProjection) {
	p.stateless[projection.Name()] = projection
}

func (p *Projector) runAllProjections(ctx context.Context, events []Event) error {
	var err error

	for _, projection := range p.projections {
		err = p.runStatefulProjection(ctx, projection, events)
	}

	for _, projection := range p.stateless {
		p.runStatelessProjection(ctx, projection, events)
	}

	return err
}

func (p *Projector) runStatefulProjection(ctx context.Context, projection Projection, events []Event) error {

	view := p.backend.NewView(projection.Name())
	state := projection.CreateState()

	if err := view.ReadView(ctx, state); err != nil {
		return err
	}

	for _, e := range events {
		state = projection.Project(ctx, state, e)
	}

	return view.WriteView(ctx, state, 0)
}

func (p *Projector) runStatelessProjection(ctx context.Context, projection StatelessProjection, events []Event) {
	for _, e := range events {
		projection.Project(e)
	}
}
