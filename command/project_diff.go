package command

import (
	"brickrecon/app"
	"brickrecon/lego"
	"brickrecon/stud_io"
	"fmt"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/honeycombio/beeline-go"
)

type ProjectDiffCommand struct {
	Meta
}

func (c *ProjectDiffCommand) Help() string {
	return ""
}

func (c *ProjectDiffCommand) Synopsis() string {
	return "Diff a lego Projects parts"
}

func (c *ProjectDiffCommand) Name() string {
	return "project diff"
}

func (c *ProjectDiffCommand) Run(args []string) int {

	ctx, send := c.NewPhase(c)
	defer send()

	flags := c.FlagSet(c)

	if err := flags.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if flags.NArg() != 2 {
		c.UI.Error("This command takes two arguments, <name> and <path>")
		return 1
	}

	projectName := lego.ProjectName(flags.Arg(0))
	filepath := flags.Arg(1)

	partsFile, err := os.Open(filepath)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	defer partsFile.Close()

	parts, err := stud_io.ReadPartsList(partsFile)
	if err != nil {
		beeline.AddField(ctx, "read_parts_error", err)
		c.UI.Error(err.Error())
		return 1
	}

	store, err := app.NewAppBuilder(ctx).CreateAppStore()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	selected, _ := store.ReadProject(ctx, projectName)
	project := lego.BlankProject()
	if err := store.EventStore.LoadAggregate(ctx, selected.ID, project); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	diff := project.Diff(parts)

	rows := []string{
		"Part Number | Colour | Quantity Change",
	}

	for key, change := range diff {
		part, c := lego.ParsePartKey(key)

		delta := ""
		if change >= 0 {
			delta = color.GreenString("%+d", change)
		} else {
			delta = color.RedString("%+d", change)
		}

		rows = append(rows, fmt.Sprintf(
			"%s | %s | %s",
			part,
			strconv.Itoa(int(c)),
			delta,
		))
	}

	c.UI.Output(tableOutput(rows))

	beeline.AddField(ctx, "complete", true)

	return 0
}
