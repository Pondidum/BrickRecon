package all_projects

import (
	"brickrecon/eventstore"
	"brickrecon/lego"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddingAuditWithExtraData(t *testing.T) {

	project := &ProjectView{Events: []*EventDescription{}}
	event := &lego.KitAddedToProject{
		EventMeta: eventstore.EventMeta{AggregateRootID: uuid.NewV4(), Timestamp: time.Now()},
		KitName:   lego.KitName("test kit"),
		KitNumber: lego.KitNumber("1234-2"),
		Parts: []lego.PartQuantity{
			{PartID: lego.LDrawPart("123"), ColourID: lego.BrickLinkColour(23), Quantity: 5},
		},
	}

	audit(project, event, "Test message")

	expected := map[string]interface{}{
		"KitName":   event.KitName,
		"KitNumber": event.KitNumber,
		"Parts":     event.Parts,
	}

	assert.Equal(t, event.KitName, project.Events[0].Additional["KitName"])
	assert.Equal(t, expected, project.Events[0].Additional)
}

func TestAddingKits(t *testing.T) {

	event := &lego.KitCreated{KitName: "test", KitNumber: "134-1", Parts: []*lego.Part{
		{Key: lego.CreatePartKey(lego.LDrawPart("1"), lego.BrickLinkColour(85)), Quantity: 5},
		{Key: lego.CreatePartKey(lego.LDrawPart("1"), lego.BrickLinkColour(17)), Quantity: 1},
		{Key: lego.CreatePartKey(lego.LDrawPart("5"), lego.BrickLinkColour(2)), Quantity: 2},
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
	projectPart := partFromKey(lego.PartKey("567|85"), 7)

	view := apply(
		createProject(projectID, projectName),
		createProjectParts(projectID, projectPart),
	)

	assert.Contains(t, view.Projects, projectName)

	expectedParts := []*ProjectPartView{
		&ProjectPartView{
			Key:      lego.PartKey("567|85"),
			ID:       lego.LDrawPart("567"),
			ColourID: lego.BrickLinkColour(85),
			Quantity: 7,
		},
	}

	assert.Equal(t, expectedParts, view.Projects[projectName].Parts)
}

func TestAddingMultipleProjectParts(t *testing.T) {

	projectID := uuid.NewV4()
	projectName := lego.ProjectName("test-project")
	partOne := partFromKey(lego.PartKey("123|85"), 1)
	partTwo := partFromKey(lego.PartKey("456|10"), 2)
	partThree := partFromKey(lego.PartKey("789|85"), 3)

	view := apply(
		createProject(projectID, projectName),
		createProjectParts(projectID, partOne, partTwo, partThree),
	)

	expectedColours := []*ColourView{
		{ID: lego.BrickLinkColour(85)},
		{ID: lego.BrickLinkColour(10)},
	}

	assert.Equal(t, expectedColours, view.Projects[projectName].Colours)
}

func TestWhenKitAddedAfterProject(t *testing.T) {

	kitNumber := lego.KitNumber("134-1")
	kitPart := partFromKey(lego.PartKey("567|85"), 5)

	projectID := uuid.NewV4()
	projectName := lego.ProjectName("test-project")
	projectPart := partFromKey(lego.PartKey("567|85"), 7)

	view := apply(
		createProject(projectID, projectName),
		createProjectParts(projectID, projectPart),
		createKit(kitNumber, kitPart),
	)

	assert.Contains(t, view.Projects, projectName)
	assert.Contains(t, view.Projects[projectName].Kits, kitNumber)
	assert.Contains(t, view.Projects[projectName].Kits[kitNumber].Parts, lego.PartKey("567|85"))

	assert.Equal(t, view.Projects[projectName].Kits[kitNumber].Parts[lego.PartKey("567|85")], 5)
}

func TestWhenProjectAddedAfterKit(t *testing.T) {

	kitNumber := lego.KitNumber("134-1")
	kitPart := partFromKey(lego.PartKey("567|85"), 5)

	projectID := uuid.NewV4()
	projectName := lego.ProjectName("test-project")
	projectPart := partFromKey(lego.PartKey("567|85"), 7)

	view := apply(
		createKit(kitNumber, kitPart),
		createProject(projectID, projectName),
		createProjectParts(projectID, projectPart),
	)

	assert.Contains(t, view.Projects, projectName)
	assert.Contains(t, view.Projects[projectName].Kits, kitNumber)
	assert.Contains(t, view.Projects[projectName].Kits[kitNumber].Parts, lego.PartKey("567|85"))

	assert.Equal(t, view.Projects[projectName].Kits[kitNumber].Parts[lego.PartKey("567|85")], 5)
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

func createProjectParts(projectID uuid.UUID, parts ...*lego.Part) *lego.ProjectPartsAdded {
	event := &lego.ProjectPartsAdded{
		EventMeta: eventstore.EventMeta{AggregateRootID: projectID},

		Parts: parts,
	}

	return event
}

func createKit(kn lego.KitNumber, kitParts ...*lego.Part) *lego.KitCreated {
	return &lego.KitCreated{
		KitName:   "test",
		KitNumber: kn,
		Parts:     kitParts,
	}
}

func partFromKey(key lego.PartKey, quantity int) *lego.Part {
	id, colour := lego.ParsePartKey(key)

	return &lego.Part{
		Key:      key,
		Aliases:  lego.PartAliases{LDrawID: id},
		Colour:   lego.Colour{ID: colour},
		Quantity: quantity,
	}
}
