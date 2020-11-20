package testing

import (
	"brickrecon/eventstore"
	"context"
)

type TestProjectionState struct {
	Names map[string]bool
}

type testProjection struct {
	name    string
	init    eventstore.Initialiser
	project func(ctx context.Context, state interface{}, event eventstore.Event) interface{}
}

func (p *testProjection) Name() string {
	return p.name
}
func (p *testProjection) CreateState() interface{} {
	return p.init()
}
func (p *testProjection) Project(ctx context.Context, state interface{}, event eventstore.Event) interface{} {
	return p.project(ctx, state, event)
}
