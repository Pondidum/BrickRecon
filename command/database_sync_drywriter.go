package command

import (
	"context"
	"fmt"
)

var _ writer = &drywriter{}

type drywriter struct {
	debug    bool
	finished bool
}

func (dw *drywriter) CreateTables(ctx context.Context, source string) error {
	fmt.Println("Dry run starting...")
	return nil
}

func (dw *drywriter) ClearTables(ctx context.Context, source string) error {
	return nil
}

func (dw *drywriter) Cancel() error {
	if !dw.finished {
		fmt.Println("Dry run cancelled")
	}
	return nil
}

func (dw *drywriter) Prepare(ctx context.Context, tableName string, header []string) (func(record []string) error, func() error, error) {
	counter := 0
	if dw.debug {
		fmt.Println("header", header)
	}
	return func(record []string) error {
			if dw.debug && counter < 5 {
				fmt.Println("insert", record)
			}
			counter++
			return nil
		}, func() error {
			fmt.Println("Would have inserted", counter, tableName, "records")
			return nil
		}, nil
}

func (dw *drywriter) Finish(ctx context.Context) error {
	dw.finished = true
	fmt.Println("Dry run finished")
	return nil
}
