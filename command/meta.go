package command

import (
	"bufio"
	"context"
	"io"

	"github.com/honeycombio/beeline-go"
	"github.com/mitchellh/cli"
	"github.com/spf13/pflag"
)

type NamedCommand interface {
	Name() string
	Help() string
}

type Meta struct {
	UI cli.Ui
}

func (m *Meta) FlagSet(cmd NamedCommand) *pflag.FlagSet {
	f := pflag.NewFlagSet(cmd.Name(), pflag.ContinueOnError)
	f.Usage = func() { m.UI.Output(cmd.Help()) }

	// Create an io.Writer that writes to our UI properly for errors.
	// This is kind of a hack, but it does the job. Basically: create
	// a pipe, use a scanner to break it into lines, and output each line
	// to the UI. Do this forever.
	errR, errW := io.Pipe()
	errScanner := bufio.NewScanner(errR)
	go func() {
		for errScanner.Scan() {
			m.UI.Error(errScanner.Text())
		}
	}()
	f.SetOutput(errW)

	return f
}

func (m *Meta) NewPhase(c NamedCommand) (context.Context, func()) {

	ctx, span := beeline.StartSpan(context.Background(), c.Name())
	return ctx, func() { span.Send() }
}
