package all_projects

import (
	"brickrecon/eventstore"
	"brickrecon/lego"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddingKits(t *testing.T) {

	event := &lego.KitCreated{KitName: "test", KitNumber: "134-1", Parts: []lego.Part{
		lego.Part{ID: lego.LDrawPart("1"), Colour: lego.Colour{ID: lego.BrickLinkColour(85)}, Quantity: 5},
		lego.Part{ID: lego.LDrawPart("1"), Colour: lego.Colour{ID: lego.BrickLinkColour(17)}, Quantity: 1},
		lego.Part{ID: lego.LDrawPart("5"), Colour: lego.Colour{ID: lego.BrickLinkColour(2)}, Quantity: 2},
	}}

	view := apply(
		event,
	)

	assert.Contains(t, view.Kits, lego.KitNumber("134-1"))
	assert.Len(t, view.Kits["134-1"].Parts, 3)
}

func TestAddingProjectParts(t *testing.T) {

	projectID := uuid.NewV4()
	projectName := lego.ProjectName("test-project")
	projectPart := lego.Part{ID: lego.LDrawPart("567"), Colour: lego.Colour{ID: lego.BrickLinkColour(85)}, Quantity: 7}

	view := apply(
		createProject(projectID, projectName),
		createProjectParts(projectID, projectPart),
	)

	assert.Contains(t, view.Projects, projectName)

	expected := []*ProjectPartView{
		&ProjectPartView{
			ID:       lego.LDrawPart("567"),
			ColourID: lego.BrickLinkColour(85),
			Key:      PartKey("567|85"),
			Quantity: 7,
		},
	}

	assert.Equal(t, expected, view.Projects[projectName].Parts)
}

func TestWhenKitAddedAfterProject(t *testing.T) {

	kitNumber := lego.KitNumber("134-1")
	kitPart := lego.Part{ID: lego.LDrawPart("567"), Colour: lego.Colour{ID: lego.BrickLinkColour(85)}, Quantity: 5}

	projectID := uuid.NewV4()
	projectName := lego.ProjectName("test-project")
	projectPart := lego.Part{ID: lego.LDrawPart("567"), Colour: lego.Colour{ID: lego.BrickLinkColour(85)}, Quantity: 7}

	view := apply(
		createProject(projectID, projectName),
		createProjectParts(projectID, projectPart),
		createKit(kitNumber, kitPart),
	)

	assert.Contains(t, view.Projects, projectName)
	assert.Contains(t, view.Projects[projectName].Kits, kitNumber)
	assert.Contains(t, view.Projects[projectName].Kits[kitNumber].Parts, PartKey("567|85"))

	assert.Equal(t, view.Projects[projectName].Kits[kitNumber].Parts[PartKey("567|85")], 5)
}

func TestWhenProjectAddedAfterKit(t *testing.T) {

	kitNumber := lego.KitNumber("134-1")
	kitPart := lego.Part{ID: lego.LDrawPart("567"), Colour: lego.Colour{ID: lego.BrickLinkColour(85)}, Quantity: 5}

	projectID := uuid.NewV4()
	projectName := lego.ProjectName("test-project")
	projectPart := lego.Part{ID: lego.LDrawPart("567"), Colour: lego.Colour{ID: lego.BrickLinkColour(85)}, Quantity: 7}

	view := apply(
		createKit(kitNumber, kitPart),
		createProject(projectID, projectName),
		createProjectParts(projectID, projectPart),
	)

	assert.Contains(t, view.Projects, projectName)
	assert.Contains(t, view.Projects[projectName].Kits, kitNumber)
	assert.Contains(t, view.Projects[projectName].Kits[kitNumber].Parts, PartKey("567|85"))

	assert.Equal(t, view.Projects[projectName].Kits[kitNumber].Parts[PartKey("567|85")], 5)
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

func createProject(projectID uuid.UUID, projectName lego.ProjectName) *lego.ProjectCreated {
	event := &lego.ProjectCreated{
		EventMeta: eventstore.EventMeta{AggregateRootID: projectID},
		ID:        projectID,
		Name:      projectName,
	}

	return event
}

func createProjectParts(projectID uuid.UUID, parts ...lego.Part) *lego.ProjectPartsAdded {
	event := &lego.ProjectPartsAdded{
		EventMeta: eventstore.EventMeta{AggregateRootID: projectID},

		Parts: parts,
	}

	return event
}

func createKit(kn lego.KitNumber, kitParts ...lego.Part) *lego.KitCreated {
	return &lego.KitCreated{
		KitName:   "test",
		KitNumber: kn,
		Parts:     kitParts,
	}
}
