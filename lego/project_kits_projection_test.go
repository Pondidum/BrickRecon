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
	projectName := ProjectName("test-project")
	projectPart := Part{ID: LDrawPart("567"), Colour: Colour{ID: BrickLinkColour(85)}, Quantity: 7}

	view := apply(
		createProject(projectID, projectName),
		createProjectParts(projectID, projectPart),
	)

	assert.Contains(t, view.Projects, projectName)

	expected := []*ProjectPartView{
		&ProjectPartView{
			ID:       LDrawPart("567"),
			ColourID: BrickLinkColour(85),
			Key:      PartKey("567|85"),
			Quantity: 7,
		},
	}

	assert.Equal(t, expected, view.Projects[projectName].Parts)
}

func TestWhenKitAddedAfterProject(t *testing.T) {

	kitNumber := KitNumber("134-1")
	kitPart := Part{ID: LDrawPart("567"), Colour: Colour{ID: BrickLinkColour(85)}, Quantity: 5}

	projectID := uuid.NewV4()
	projectName := ProjectName("test-project")
	projectPart := Part{ID: LDrawPart("567"), Colour: Colour{ID: BrickLinkColour(85)}, Quantity: 7}

	view := apply(
		createProject(projectID, projectName),
		createProjectParts(projectID, projectPart),
		createKit(kitNumber, kitPart),
	)

	assert.Contains(t, view.Projects, projectName)
	assert.Contains(t, view.Projects[projectName].Kits, kitNumber)
	assert.Contains(t, view.Projects[projectName].Kits[kitNumber], PartKey("567|85"))

	assert.Equal(t, view.Projects[projectName].Kits[kitNumber][PartKey("567|85")], 5)
}

func TestWhenProjectAddedAfterKit(t *testing.T) {

	kitNumber := KitNumber("134-1")
	kitPart := Part{ID: LDrawPart("567"), Colour: Colour{ID: BrickLinkColour(85)}, Quantity: 5}

	projectID := uuid.NewV4()
	projectName := ProjectName("test-project")
	projectPart := Part{ID: LDrawPart("567"), Colour: Colour{ID: BrickLinkColour(85)}, Quantity: 7}

	view := apply(
		createKit(kitNumber, kitPart),
		createProject(projectID, projectName),
		createProjectParts(projectID, projectPart),
	)

	assert.Contains(t, view.Projects, projectName)
	assert.Contains(t, view.Projects[projectName].Kits, kitNumber)
	assert.Contains(t, view.Projects[projectName].Kits[kitNumber], PartKey("567|85"))

	assert.Equal(t, view.Projects[projectName].Kits[kitNumber][PartKey("567|85")], 5)
}

func apply(events ...eventstore.Event) *AllProjectsView {

	p := &ProjectsProjection{}
	state := p.CreateState()

	for _, e := range events {
		p.Project(state, e)
	}

	view := state.(*AllProjectsView)

	return view
}

func createProject(projectID uuid.UUID, projectName ProjectName) *ProjectCreated {
	event := &ProjectCreated{
		EventMeta: eventstore.EventMeta{AggregateRootID: projectID},
		ID:        projectID,
		Name:      projectName,
	}

	return event
}

func createProjectParts(projectID uuid.UUID, parts ...Part) *ProjectPartsAdded {
	event := &ProjectPartsAdded{
		EventMeta: eventstore.EventMeta{AggregateRootID: projectID},

		Parts: parts,
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
