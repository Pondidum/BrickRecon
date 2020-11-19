package command

import (
	"brickrecon/app"
	"brickrecon/lego/projections/allprojects"
	"fmt"

	"github.com/honeycombio/beeline-go"
)

type ProjectListCommand struct {
	Meta
}

func (c *ProjectListCommand) Help() string {
	return ""
}

func (c *ProjectListCommand) Synopsis() string {
	return "Lists all Projects"
}

func (c *ProjectListCommand) Name() string {
	return "project list"
}

func (c *ProjectListCommand) Run(args []string) int {

	ctx, send := c.NewPhase(c)
	defer send()

	flags := c.FlagSet(c)

	if err := flags.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if flags.NArg() != 0 {
		c.UI.Error("This command takes no arguments")
		return 1
	}

	store, err := app.NewAppBuilder(ctx).CreateAppStore()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	var view allprojects.AllProjectsView
	if err := store.EventStore.ReadView(ctx, allprojects.ProjectionName, &view); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	rows := []string{
		"Name | Part Count | UUID",
	}
	for _, project := range view.Projects {

		rows = append(rows, fmt.Sprintf("%s | %d | %s",
			string(project.Name),
			len(project.Parts),
			project.ID,
		))
	}

	c.UI.Output(tableOutput(rows))

	beeline.AddField(ctx, "complete", true)

	return 0
}
