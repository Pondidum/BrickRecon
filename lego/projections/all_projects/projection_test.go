package all_projects

import (
	"brickrecon/eventstore"
	"brickrecon/lego"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddingAuditWithExtraData(t *testing.T) {

	project := &ProjectView{Events: []*EventDescription{}}
	event := &lego.KitAddedToProject{
		EventMeta: eventstore.EventMeta{AggregateRootID: eventstore.NewAggregateID(), Timestamp: time.Now()},
		KitName:   lego.KitName("test kit"),
		KitNumber: lego.KitNumber("1234-2"),
		Parts: map[lego.PartKey]int{
			lego.PartKey("123|23"): 5,
		},
	}

	project.audit(event, "Test message")

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
		{Key: lego.PartKey("1|85"), Quantity: 5},
		{Key: lego.PartKey("1|17"), Quantity: 1},
		{Key: lego.PartKey("5|2"), Quantity: 2},
	}}

	view := apply(
		event,
	)

	assert.Contains(t, view.Kits, lego.KitNumber("134-1"))
	assert.Len(t, view.Kits["134-1"].Parts, 3)
}

func TestAddingProjectParts(t *testing.T) {

	projectID := eventstore.NewAggregateID()
	projectName := lego.ProjectName("test-project")
	projectPart := partFromKey(lego.PartKey("567|72"), 7)
	projectPart.Colour.Aliases.BrickLinkID = lego.BrickLinkColour(85)

	view := apply(
		createProject(projectID, projectName),
		createProjectParts(projectID, projectPart),
	)

	assert.Contains(t, view.Projects, projectName)

	expectedParts := []*ProjectPartView{
		&ProjectPartView{
			Key:       lego.PartKey("567|72"),
			ID:        lego.LDrawPart("567"),
			ColourID:  lego.LDrawColour(72),
			ImagePath: "567-85.png",
			Quantity:  7,
		},
	}

	assert.Equal(t, expectedParts, view.Projects[projectName].Parts)
}

func TestAddingMultipleProjectParts(t *testing.T) {

	projectID := eventstore.NewAggregateID()
	projectName := lego.ProjectName("test-project")
	partOne := partFromKey(lego.PartKey("123|85"), 1)
	partTwo := partFromKey(lego.PartKey("456|10"), 2)
	partThree := partFromKey(lego.PartKey("789|85"), 3)

	view := apply(
		createProject(projectID, projectName),
		createProjectParts(projectID, partOne, partTwo, partThree),
	)

	expectedColours := []*ColourView{
		{ID: lego.LDrawColour(85)},
		{ID: lego.LDrawColour(10)},
	}

	assert.Equal(t, expectedColours, view.Projects[projectName].Colours)
}

func TestWhenKitAddedAfterProject(t *testing.T) {

	kitNumber := lego.KitNumber("134-1")
	kitPart := partFromKey(lego.PartKey("567|85"), 5)

	projectID := eventstore.NewAggregateID()
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

	projectID := eventstore.NewAggregateID()
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

	p := NewProjectsProjection(nil)
	state := p.CreateState()

	for _, e := range events {
		p.Project(state, e)
	}

	view := state.(*AllProjectsView)

	return view
}

func createProject(projectID eventstore.AggregateID, projectName lego.ProjectName) *lego.ProjectCreated {
	event := &lego.ProjectCreated{
		EventMeta: eventstore.EventMeta{AggregateRootID: projectID},
		ID:        projectID,
		Name:      projectName,
	}

	return event
}

func createProjectParts(projectID eventstore.AggregateID, parts ...*lego.Part) *lego.ProjectPartsAdded {
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
		Aliases:  lego.PartAliases{LDrawID: id, BrickLinkID: lego.BrickLinkPart(id)},
		Colour:   lego.Colour{Aliases: lego.ColourAliases{LDrawID: colour}},
		Quantity: quantity,
	}
}
