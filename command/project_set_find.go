package command

import (
	"brickrecon/config"
	"brickrecon/domain"
	"brickrecon/lego"
	"brickrecon/storage"
	"brickrecon/tracing"
	"brickrecon/util"
	"context"
	"fmt"
	"sort"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func NewProjectSetFindCommand() *ProjectSetFindCommand {
	return &ProjectSetFindCommand{
		tr: otel.Tracer("command.project.set.find"),
	}
}

type ProjectSetFindCommand struct {
	tr trace.Tracer
}

func (c *ProjectSetFindCommand) Name() string {
	return "project set find"
}

func (c *ProjectSetFindCommand) Synopsis() string {
	return "find sets to build a project"
}

func (c *ProjectSetFindCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("project set find", pflag.ContinueOnError)
	return flags
}

func (c *ProjectSetFindCommand) Execute(ctx context.Context, config *config.Config, args []string) error {
	ctx, span := c.tr.Start(ctx, "execute")
	defer span.End()

	if len(args) != 1 {
		return tracing.Errorf(span, "this command takes exactly 1 argument: project_name")
	}

	store, err := storage.NewClient(ctx, config.DatabaseFile)
	if err != nil {
		return tracing.Error(span, err)
	}

	project, err := storage.GetProjectView(ctx, store, storage.WithName(args[0]))
	if err != nil {
		return tracing.Error(span, err)
	}

	missing := []*domain.ProjectPart{}

	for _, part := range project.Parts {
		stock := domain.GetStock(project.Stock, part.Number, part.Color)
		if stock >= part.Wanted {
			continue
		}

		missing = append(missing, part)
	}

	fmt.Printf("Finding sets for %d unique unstocked parts...", len(missing))

	setSummaries := make(map[string][]*storage.SetSummary, len(missing))

	for _, part := range missing {
		sets, err := storage.GetLegoSetsForPart(ctx, store, part.Number, part.Color)
		if err != nil {
			return err
		}

		setSummaries[part.Key()] = sets
	}

	fmt.Println("Done")

	noSetParts := map[string]any{}
	singleSetParts := map[string]*storage.SetSummary{}

	setGroups := map[lego.SetNumber][]*storage.SetSummary{}

	for key, sets := range setSummaries {
		switch len(sets) {
		case 0:
			noSetParts[key] = true
		case 1:
			singleSetParts[key] = sets[0]
		default:
			for _, summary := range sets {

				grouped := []*storage.SetSummary{}
				if group, found := setGroups[summary.SetNumber]; found {
					grouped = group
				}

				setGroups[summary.SetNumber] = append(grouped, summary)
			}
		}
	}

	if len(noSetParts) > 0 {
		fmt.Printf("%d parts are not in any sets:\n", len(noSetParts))
		for key := range noSetParts {
			fmt.Println("* ", key)
		}
	}

	if len(singleSetParts) > 0 {
		fmt.Printf("%d parts are 1 set only:\n", len(singleSetParts))
		for key, summary := range singleSetParts {
			fmt.Printf("* %s: %s, %s (%d)\n", key, summary.SetNumber, summary.SetName, summary.SetYear)
		}
	}

	if len(setGroups) > 0 {

		fmt.Printf("%d sets provide parts:\n", len(setGroups))

		groups := make([][]*storage.SetSummary, 0, len(setGroups))
		for _, group := range setGroups {
			groups = append(groups, group)
		}

		sort.Slice(groups, func(i, j int) bool {
			return len(groups[i]) > len(groups[j])
		})

		lines := make([]string, 0, len(groups)+2)
		lines = append(lines, "Set Number | Set Name | Year | Total Parts | Stock | Stock Percent")

		for _, group := range groups {
			lines = append(lines, fmt.Sprintf("%s | %s | %d | %d | %d | %.2f",
				group[0].SetNumber,
				group[0].SetName,
				group[0].SetYear,
				group[0].TotalParts,
				len(group),
				float64(len(group))/float64(group[0].TotalParts)*100,
			))
		}

		fmt.Println(util.TableOutput(lines))
	}
	return nil
}
