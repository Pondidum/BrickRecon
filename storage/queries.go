package storage

import (
	"brickrecon/domain"
	"brickrecon/goes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

var ErrViewNotFound = errors.New("no matching view found")

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

func GetProjectViewByName(ctx context.Context, client *Client, name string) (*domain.ProjectView, error) {
	row := client.db.QueryRowContext(
		ctx,
		`select aggregate_id, view from auto_projections where aggregate_type = 'Project' and view ->> '$.Name' == @name`,
		sql.Named("name", name),
	)

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

func GetProjectViewsAll(ctx context.Context, client *Client) ([]*domain.ProjectView, error) {

	row, err := client.db.QueryContext(ctx, `select aggregate_id, view from auto_projections where aggregate_type = 'Project'`)
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
