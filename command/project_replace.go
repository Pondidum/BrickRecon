package command

import (
	"brickrecon/app"
	"brickrecon/lego"
	"brickrecon/stud_io"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/honeycombio/beeline-go"
	"github.com/posener/complete"
)

type ProjectReplaceCommand struct {
	Meta
}

func (c *ProjectReplaceCommand) Help() string {
	return ""
}

func (c *ProjectReplaceCommand) Synopsis() string {
	return "Replaces the parts in a lego Projects"
}

func (c *ProjectReplaceCommand) Name() string {
	return "project replace"
}

func (c *ProjectReplaceCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{
		"--dry": complete.PredictNothing,
	}
}

func (c *ProjectReplaceCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ProjectReplaceCommand) Run(args []string) int {

	ctx, send := c.NewPhase(c)
	defer send()

	var dryRun bool

	flags := c.FlagSet(c)
	flags.BoolVar(&dryRun, "dry", false, "")

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

	project, err := store.ReadProject(ctx, projectName)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if dryRun {
		diff := project.Diff(parts)
		c.outputTable(diff)
	} else {
		diff := project.ReplaceParts(parts)
		c.outputTable(diff)

		c.UI.Info("Saving project...")

		if err := store.Save(ctx, project); err != nil {
			c.UI.Error(err.Error())
			return 1
		}

		c.UI.Output("Done.")
	}

	beeline.AddField(ctx, "complete", true)

	return 0
}

func (c *ProjectReplaceCommand) outputTable(diff map[lego.PartKey]int) {
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
			"%s | %v | %s",
			part,
			c,
			delta,
		))
	}

	c.UI.Output(tableOutput(rows))
}
