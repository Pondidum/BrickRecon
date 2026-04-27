package command

import (
	"brickrecon/config"
	"brickrecon/domain"
	"brickrecon/lego"
	"brickrecon/storage"
	"brickrecon/tracing"
	"context"
	"fmt"
	"os"
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

	renderer string
}

func (c *ProjectSetFindCommand) Name() string {
	return "project set find"
}

func (c *ProjectSetFindCommand) Synopsis() string {
	return "find sets to build a project"
}

func (c *ProjectSetFindCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("project set find", pflag.ContinueOnError)
	flags.StringVar(&c.renderer, "format", "table", "")
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

	for partKey, sets := range setSummaries {
		switch len(sets) {
		case 0:
			noSetParts[partKey] = true
		case 1:
			singleSetParts[partKey] = sets[0]
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
		fmt.Printf("%d parts are 1 set only:\n\n", len(singleSetParts))

		sets := make([]setSummary, 0, len(singleSetParts))

		for _, s := range singleSetParts {
			stockPercent := float64(s.PartQuantity) / float64(s.TotalParts) * 100

			sets = append(sets, setSummary{
				SetNumber:        s.SetNumber,
				SetName:          s.SetName,
				Year:             s.SetYear,
				TotalParts:       s.TotalParts,
				Stock:            s.PartQuantity,
				StockPercent:     fmt.Sprintf("%.2f", stockPercent),
				stockPercentSort: stockPercent,
			})
		}

		if err := Render(c.renderer, os.Stdout, sets); err != nil {
			return tracing.Error(span, err)
		}

		fmt.Println()
	}

	if len(setGroups) > 0 {

		fmt.Printf("%d sets provide parts:\n\n", len(setGroups))
		sets := make([]setSummary, 0, len(setGroups))

		for _, group := range setGroups {

			totalAdding := 0
			for _, s := range group {
				totalAdding += s.PartQuantity
			}

			stockPercent := float64(totalAdding) / float64(group[0].TotalParts) * 100

			g := group[0]
			sets = append(sets, setSummary{
				SetNumber:        g.SetNumber,
				SetName:          g.SetName,
				Year:             g.SetYear,
				TotalParts:       g.TotalParts,
				Stock:            totalAdding,
				StockPercent:     fmt.Sprintf("%.2f", stockPercent),
				stockPercentSort: stockPercent,
			})
		}

		sort.Slice(sets, func(i, j int) bool {
			return sets[i].StockPercent > sets[j].StockPercent
		})

		if err := Render(c.renderer, os.Stdout, sets); err != nil {
			return tracing.Error(span, err)
		}
	}
	return nil
}

type setSummary struct {
	SetNumber    lego.SetNumber
	SetName      lego.SetName
	Year         int
	TotalParts   int
	Stock        int
	StockPercent string

	stockPercentSort float64
}
