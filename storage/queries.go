package storage

import (
	"brickrecon/domain"
	"brickrecon/goes"
	"brickrecon/lego"
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

var ErrViewNotFound = errors.New("no matching view found")

type ViewOptions struct {
	includeArchived bool
	name            string
}

func (vo *ViewOptions) apply(funcs []ViewOption) {
	for _, fn := range funcs {
		fn(vo)
	}
}

type ViewOption func(o *ViewOptions)

func IncludeArchived() ViewOption {
	return func(o *ViewOptions) {
		o.includeArchived = true
	}
}

func WithName(name string) ViewOption {
	return func(o *ViewOptions) {
		o.name = name
	}
}

func GetProjectByName(ctx context.Context, client *Client, name string) (*domain.Project, error) {

	row := client.db.QueryRowContext(
		ctx,
		`select aggregate_id from auto_projections where aggregate_type = 'Project' and view ->> '$.Name' == @name`,
		sql.Named("name", name),
	)

	var id uuid.UUID
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil, goes.ErrNotFound
		}
		return nil, err
	}

	project := domain.BlankProject()
	if err := client.LoadAggregate(ctx, id, project); err != nil {
		return nil, err
	}

	return project, nil
}

func GetProjectView(ctx context.Context, client *Client, options ...ViewOption) (*domain.ProjectView, error) {

	opt := &ViewOptions{}
	opt.apply(options)

	stmt := `select aggregate_id, view from auto_projections where aggregate_type = 'Project'`
	params := []any{}
	if !opt.includeArchived {
		stmt = stmt + ` and (view ->> '$.Archived' is null or view ->> '$.Archived' = false)`
	}
	if opt.name != "" {
		stmt = stmt + ` and view ->> '$.Name' == @name`
		params = append(params, sql.Named("name", opt.name))
	}

	row := client.db.QueryRowContext(ctx, stmt, params...)

	var aggregateId uuid.UUID
	var viewJson []byte
	if err := row.Scan(&aggregateId, &viewJson); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrViewNotFound
		} else {
			return nil, err
		}
	}

	view := &domain.ProjectView{}
	if err := json.Unmarshal(viewJson, view); err != nil {
		return nil, err
	}
	view.AggregateID = aggregateId

	return view, nil
}
func GetProjectViews(ctx context.Context, client *Client, options ...ViewOption) ([]*domain.ProjectView, error) {

	opt := &ViewOptions{}
	opt.apply(options)

	stmt := `select aggregate_id, view from auto_projections where aggregate_type = 'Project'`
	if !opt.includeArchived {
		stmt = stmt + ` and (view ->> '$.Archived' is null or view ->> '$.Archived' = false)`
	}

	row, err := client.db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}

	projects := []*domain.ProjectView{}
	for row.Next() {

		var aggregateId uuid.UUID
		var viewJson []byte
		if err := row.Scan(&aggregateId, &viewJson); err != nil {
			return nil, err
		}

		view := &domain.ProjectView{}
		if err := json.Unmarshal(viewJson, view); err != nil {
			return nil, err
		}
		view.AggregateID = aggregateId

		projects = append(projects, view)
	}

	return projects, nil
}

func FindMatchingParts(ctx context.Context, client *Client, partId lego.PartId) ([]*lego.Part, error) {
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

	rows, err := client.db.QueryContext(ctx, stmt, sql.Named("part_num", partId))
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
