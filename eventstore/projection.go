package eventstore

import "context"

type Projection interface {
	Name() string
	CreateState() interface{}
	Project(state interface{}, event Event) interface{}
}

func runProjection(ctx context.Context, backend Backend, events []Event, projection Projection) error {
	view := backend.NewView(projection.Name())
	state := projection.CreateState()

	if err := view.ReadView(ctx, state); err != nil {
		return err
	}

	for _, e := range events {
		state = projection.Project(state, e)
	}

	return view.WriteView(ctx, state, 0)

}
