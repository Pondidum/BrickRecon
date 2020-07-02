package lego

import (
	"brickrecon/eventstore"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddingKits(t *testing.T) {

	event := &KitCreated{KitName: "test", KitNumber: "134-1", Parts: []Part{
		Part{ID: LDrawPart("1"), Colour: Colour{ID: BrickLinkColour(85)}, Quantity: 5},
		Part{ID: LDrawPart("1"), Colour: Colour{ID: BrickLinkColour(17)}, Quantity: 1},
		Part{ID: LDrawPart("5"), Colour: Colour{ID: BrickLinkColour(2)}, Quantity: 2},
	}}

	view := apply(
		event,
	)

	assert.Contains(t, view.Kits, KitNumber("134-1"))
	assert.Len(t, view.Kits["134-1"], 3)
}

func TestAddingProjectParts(t *testing.T) {

	projectID := uuid.NewV4()
	projectPart := Part{ID: LDrawPart("567"), Colour: Colour{ID: BrickLinkColour(85)}, Quantity: 7}

	view := apply(
		createProjectParts(projectID, projectPart),
	)

	assert.Contains(t, view.Projects, projectID)

	assert.Equal(t, view.Projects[projectID].Parts, []PartRequirement{
		PartRequirement{PartID: LDrawPart("567"), Colour: BrickLinkColour(85), Key: PartKey("567|85"), Quantity: 7},
	})
}

func TestWhenKitAddedAfterProject(t *testing.T) {

	kitNumber := KitNumber("134-1")
	kitPart := Part{ID: LDrawPart("567"), Colour: Colour{ID: BrickLinkColour(85)}, Quantity: 5}

	projectID := uuid.NewV4()
	projectPart := Part{ID: LDrawPart("567"), Colour: Colour{ID: BrickLinkColour(85)}, Quantity: 7}

	view := apply(
		createProjectParts(projectID, projectPart),
		createKit(kitNumber, kitPart),
	)

	assert.Contains(t, view.Projects, projectID)
	assert.Contains(t, view.Projects[projectID].Kits, kitNumber)
	assert.Contains(t, view.Projects[projectID].Kits[kitNumber], PartKey("567|85"))

	assert.Equal(t, view.Projects[projectID].Kits[kitNumber][PartKey("567|85")], 5)
}

func TestWhenProjectAddedAfterKit(t *testing.T) {

	kitNumber := KitNumber("134-1")
	kitPart := Part{ID: LDrawPart("567"), Colour: Colour{ID: BrickLinkColour(85)}, Quantity: 5}

	projectID := uuid.NewV4()
	projectPart := Part{ID: LDrawPart("567"), Colour: Colour{ID: BrickLinkColour(85)}, Quantity: 7}

	view := apply(
		createKit(kitNumber, kitPart),
		createProjectParts(projectID, projectPart),
	)

	assert.Contains(t, view.Projects, projectID)
	assert.Contains(t, view.Projects[projectID].Kits, kitNumber)
	assert.Contains(t, view.Projects[projectID].Kits[kitNumber], PartKey("567|85"))

	assert.Equal(t, view.Projects[projectID].Kits[kitNumber][PartKey("567|85")], 5)
}

func apply(events ...eventstore.Event) *ProjectKitsView {

	p := &ProjectKitsProjection{}
	state := p.CreateState()

	for _, e := range events {
		state = p.Project(state, e)
	}

	view := state.(*ProjectKitsView)
	return view
}

func createProjectParts(projectID uuid.UUID, parts ...Part) *ProjectPartsAdded {
	event := &ProjectPartsAdded{
		EventMeta: eventstore.EventMeta{AggregateRootID: projectID},
		Parts:     parts,
	}

	return event
}

func createKit(kn KitNumber, kitParts ...Part) *KitCreated {
	return &KitCreated{
		KitName:   "test",
		KitNumber: kn,
		Parts:     kitParts,
	}
}
