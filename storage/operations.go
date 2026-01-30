package storage

import (
	"brickrecon/lego"
	"brickrecon/tracing"
	"context"
	"database/sql"
	"strings"
)

func InsertLegoParts(ctx context.Context, tx *sql.Tx, parts []*lego.InventoryPart) error {

	stmt := `insert into parts(id, name) values `
	values := []any{}

	for _, part := range parts {
		stmt += "(?, ?),"
		values = append(values, string(part.Id), string(part.Name))
	}

	stmt = strings.TrimRight(stmt, ",")
	stmt = stmt + ` on conflict(id) do nothing`

	_, err := tx.ExecContext(ctx, stmt, values...)
	return err
}

func linkParts(ctx context.Context, tx *sql.Tx, setNumber lego.SetNumber, parts []*lego.InventoryPart) error {

	stmt := `insert into sets_parts(set_id, part_id, color, quantity) values `
	values := []any{}

	for _, part := range parts {
		stmt += "(?, ?, ?, ?),"
		values = append(values, string(setNumber), string(part.Id), string(part.ColorId), part.Quantity)
	}

	stmt = strings.TrimRight(stmt, ",")

	_, err := tx.ExecContext(ctx, stmt, values...)
	return err
}

func InsertLegoSet(ctx context.Context, tx *sql.Tx, legoSet *lego.Set) error {
	_, err := tx.ExecContext(ctx,
		`insert into sets(id, name) values (@id, @name)`,
		sql.Named("id", legoSet.Number),
		sql.Named("name", legoSet.Name))
	if err != nil {
		return err
	}

	if err := InsertLegoParts(ctx, tx, legoSet.Parts); err != nil {
		return err
	}

	if err := linkParts(ctx, tx, legoSet.Number, legoSet.Parts); err != nil {
		return err
	}

	return nil
}

func GetLegoSetByNumber(ctx context.Context, tx *sql.Tx, setNumber lego.SetNumber) (*lego.Set, error) {
	ctx, span := tr.Start(ctx, "get_set_by_number")
	defer span.End()

	setStmt := `select id, name from sets where id = @set_number`
	row := tx.QueryRowContext(ctx, setStmt, sql.Named("set_number", string(setNumber)))

	ls := &lego.Set{}
	if err := row.Scan(&ls.Number, ls.Name); err != nil {
		return nil, tracing.Error(span, err)
	}

	partsStmt := `
select
	p.id,
	p.name,
	sp.color,
	sp.quantity
from sets_parts sp
join parts p on sp.part_id = p.id
where sp.set_id = @set_number
`
	rows, err := tx.QueryContext(ctx, partsStmt, sql.Named("set_number", string(setNumber)))

	if err != nil {
		return nil, tracing.Error(span, err)
	}

	for rows.Next() {
		part := &lego.InventoryPart{
			Part: lego.Part{},
		}
		if rows.Scan(part.Id, part.Name, part.ColorId, part.Quantity); err != nil {
			return nil, tracing.Error(span, err)
		}

		ls.Parts = append(ls.Parts, part)
	}

	return ls, nil
}
