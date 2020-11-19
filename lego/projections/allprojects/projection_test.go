package allprojects

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

	event := &lego.KitCreated{KitName: "test", KitNumber: "134-1", Parts: map[lego.PartKey]int{
		lego.PartKey("1|85"): 5,
		lego.PartKey("1|17"): 1,
		lego.PartKey("5|2"):  2,
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

	view := apply(
		createProject(projectID, projectName),
		createProjectParts(projectID, map[lego.PartKey]int{
			lego.PartKey("567|72"): 7,
		}),
	)

	assert.Contains(t, view.Projects, projectName)

	expectedParts := []*ProjectPartView{
		&ProjectPartView{
			Key:       lego.PartKey("567|72"),
			ID:        lego.LDrawPart("567"),
			ColourID:  lego.LDrawColour(72),
			ColourHex: "595D60",
			Quantity:  7,
		},
	}

	assert.Equal(t, expectedParts, view.Projects[projectName].Parts)
}

func TestAddingMultipleProjectParts(t *testing.T) {

	projectID := eventstore.NewAggregateID()
	projectName := lego.ProjectName("test-project")

	view := apply(
		createProject(projectID, projectName),
		createProjectParts(projectID, map[lego.PartKey]int{
			lego.PartKey("123|85"): 1,
			lego.PartKey("456|10"): 2,
			lego.PartKey("789|85"): 3,
		}),
	)

	expectedColours := []*ColourView{
		{ID: lego.LDrawColour(85), Hex: lego.GetColourHex(lego.LDrawColour(85))},
		{ID: lego.LDrawColour(10), Hex: lego.GetColourHex(lego.LDrawColour(10))},
	}

	assert.Equal(t, expectedColours, view.Projects[projectName].Colours)
}

func TestWhenKitAddedAfterProject(t *testing.T) {

	kitNumber := lego.KitNumber("134-1")

	projectID := eventstore.NewAggregateID()
	projectName := lego.ProjectName("test-project")

	view := apply(
		createProject(projectID, projectName),
		createProjectParts(projectID, map[lego.PartKey]int{
			lego.PartKey("567|85"): 7,
		}),
		createKit(kitNumber, map[lego.PartKey]int{
			lego.PartKey("567|85"): 5,
		}),
	)

	assert.Contains(t, view.Projects, projectName)
	assert.Contains(t, view.Projects[projectName].Kits, kitNumber)
	assert.Contains(t, view.Projects[projectName].Kits[kitNumber].Parts, lego.PartKey("567|85"))

	assert.Equal(t, view.Projects[projectName].Kits[kitNumber].Parts[lego.PartKey("567|85")], 5)
}

func TestWhenProjectAddedAfterKit(t *testing.T) {

	kitNumber := lego.KitNumber("134-1")

	projectID := eventstore.NewAggregateID()
	projectName := lego.ProjectName("test-project")

	view := apply(
		createKit(kitNumber, map[lego.PartKey]int{
			lego.PartKey("567|85"): 5,
		}),
		createProject(projectID, projectName),
		createProjectParts(projectID, map[lego.PartKey]int{
			lego.PartKey("567|85"): 7,
		}),
	)

	assert.Contains(t, view.Projects, projectName)
	assert.Contains(t, view.Projects[projectName].Kits, kitNumber)
	assert.Contains(t, view.Projects[projectName].Kits[kitNumber].Parts, lego.PartKey("567|85"))

	assert.Equal(t, view.Projects[projectName].Kits[kitNumber].Parts[lego.PartKey("567|85")], 5)
}

func apply(events ...eventstore.Event) *AllProjectsView {

	p := NewProjectsProjection(nil)
	p.partLoader = func(k lego.PartKey) *lego.Part {
		part := lego.BlankPart()

		n, c := lego.ParsePartKey(k)
		part.PartID = n
		part.ColourID = c
		part.ColourHex = lego.GetColourHex(c)

		return part
	}

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

func createProjectParts(projectID eventstore.AggregateID, parts map[lego.PartKey]int) *lego.ProjectPartsAdded {
	event := &lego.ProjectPartsAdded{
		EventMeta: eventstore.EventMeta{AggregateRootID: projectID},
		Parts:     parts,
	}

	return event
}

func createKit(kn lego.KitNumber, kitParts map[lego.PartKey]int) *lego.KitCreated {
	return &lego.KitCreated{
		KitName:   "test",
		KitNumber: kn,
		Parts:     kitParts,
	}
}
