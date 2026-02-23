package command

import (
	"brickrecon/bricklink"
	"brickrecon/config"
	"brickrecon/domain"
	"brickrecon/lego"
	"brickrecon/storage"
	"brickrecon/tracing"
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func NewProjectPartsImportCommand() *ProjectPartsImportCommand {
	return &ProjectPartsImportCommand{
		tr: otel.Tracer("command.project.parts.import"),
	}
}

type ProjectPartsImportCommand struct {
	tr     trace.Tracer
	dryrun bool
}

func (c *ProjectPartsImportCommand) Name() string {
	return "project parts import"
}

func (c *ProjectPartsImportCommand) Synopsis() string {
	return "import a parts list to a project "
}

func (c *ProjectPartsImportCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("project parts import", pflag.ContinueOnError)
	flags.BoolVar(&c.dryrun, "dry-run", false, "")
	return flags
}

func (c *ProjectPartsImportCommand) Execute(ctx context.Context, config *config.Config, args []string) error {
	ctx, span := c.tr.Start(ctx, "execute")
	defer span.End()

	if len(args) != 2 {
		return tracing.Errorf(span, "this command takes exactly 2 arguments: name and parts-path")
	}

	store, err := storage.NewClient(ctx, config.DatabaseFile)
	if err != nil {
		return tracing.Error(span, err)
	}

	name := args[0]
	wantedList := args[1]

	project, err := storage.GetProjectByName(ctx, store, name)
	if err != nil {
		return tracing.Error(span, err)
	}

	fmt.Println(project.Name, "currently has", len(project.Parts), "parts")

	content, err := os.Open(wantedList)
	if err != nil {
		return tracing.Error(span, err)
	}
	defer content.Close()

	tx, err := store.BeginTx(ctx)
	if err != nil {
		return tracing.Error(span, err)
	}
	defer tx.Rollback()

	getPart := func(partId lego.PartId) (*lego.Part, error) {
		parts, err := GetPart(ctx, tx, partId)
		if err != nil {
			return nil, err
		}
		switch len(parts) {
		case 0:
			return nil, fmt.Errorf("unable to find any bricks matching %s", partId)
		case 1:
			return parts[0], nil
		default:
			// we could query the user here to pick the right one
			return parts[0], nil
		}
	}

	parts, stock, err := bricklink.ParseWantedList(ctx, getPart, content)
	if err != nil {
		return tracing.Error(span, err)
	}

	fmt.Println("Wanted List has", len(parts))

	if c.dryrun {
		fmt.Println("Would import", len(parts))

		for _, part := range parts {
			fmt.Println(part.Id, part.Name, lego.GetColorName(part.ColorId), part.Quantity, domain.GetStock(stock, part.Id, part.ColorId))
		}

		return nil
	}

	if err := project.AddParts(parts); err != nil {
		return tracing.Error(span, err)
	}

	if err := project.AddStock(stock); err != nil {
		return tracing.Error(span, err)
	}

	if err := store.SaveAggregate(ctx, project); err != nil {
		return tracing.Error(span, err)
	}

	if err := tx.Commit(); err != nil {
		return tracing.Error(span, err)
	}

	fmt.Println("Done")

	return nil
}

func GetPart(ctx context.Context, tx *sql.Tx, partId lego.PartId) ([]*lego.Part, error) {
	stmt := `
	
		select rp.part_num, rp.name, rpc.name, 0 "priority"
		from rebrickable_parts rp 
		join rebrickable_part_categories rpc on rpc.id = rp.category_id 
		where rp.part_num = @part_num 

		union

		select rp.part_num, rp.name, rpc.name, 1 "priority"
		from rebrickable_parts rp 
		join rebrickable_part_categories rpc on rpc.id = rp.category_id 
		join ldraw_alternate_ids lai on rp.part_num = lai.part_num
		where lai.alt_part_num = @part_num 

		union
	
		select rp.part_num, rp.name, rp.name, 2 "priority"
		from rebrickable_parts rp 
		join (
			select lai.alt_part_num "part_num"
			from ldraw_alternate_ids lai 
			where lai.part_num in (
				select lai.part_num
				from ldraw_alternate_ids lai 
				where lai.alt_part_num = @part_num
			)) alt on alt.part_num = rp.part_num

		union

		select rp.part_num, rp.name, rpc.name, 3 "priority"
		from rebrickable_parts rp 
		join rebrickable_part_categories rpc on rpc.id = rp.category_id 
		join ldraw_moved_parts lmp on lmp.new_part_num  = rp.part_num
		where lmp.old_part_num = @part_num 

		union

		select rp.part_num , rp.name, rpc.name, 4 "priority"
		from rebrickable_parts rp
		join rebrickable_part_categories rpc on rpc.id = rp.category_id
		where rp.part_num like concat(@part_num, '%')

		order by priority, rp.part_num
 	`

	rows, err := tx.QueryContext(ctx, stmt, sql.Named("part_num", partId))
	if err != nil {
		return nil, err
	}

	parts := []*lego.Part{}
	for rows.Next() {

		part := &lego.Part{}
		priority := 0
		if err := rows.Scan(&part.Id, &part.Name, &part.Category, &priority); err != nil {
			return nil, err
		}

		parts = append(parts, part)
	}

	return parts, nil

}
