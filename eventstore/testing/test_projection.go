package testing

import "brickrecon/eventstore"

type TestProjectionState struct {
	Names map[string]bool
}

type testProjection struct {
	name    string
	init    eventstore.Initialiser
	project eventstore.Projector
}

func (p *testProjection) Name() string {
	return p.name
}
func (p *testProjection) CreateState() interface{} {
	return p.init()
}
func (p *testProjection) Project(state interface{}, event eventstore.Event) interface{} {
	return p.project(state, event)
}
