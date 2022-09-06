package debugger

import (
	"context"
	"fmt"

	"github.com/kr/pretty"
	"github.com/spf13/cobra"
)

type contextKeyType string

var contextKey contextKeyType = "debug"

type Debugger struct {
	debug *bool
}

func WithDebugger(ctx context.Context, debug *bool) context.Context {
	dbg := Debugger{debug}
	return context.WithValue(ctx, contextKey, &dbg)
}

func (d *Debugger) Debugf(format string, args ...interface{}) {
	if *d.debug {
		fmt.Printf(fmt.Sprintf("DEBUG: %s\n", format), args...)
	}
}

func (d *Debugger) Flags(cmd *cobra.Command) {
	if *d.debug {
		cmd.DebugFlags()
	}
}

func (d *Debugger) Prettyln(arg ...interface{}) {
	if *d.debug {
		pretty.Println(arg...)
	}
}

func Get(ctx context.Context) *Debugger {
	dbg := ctx.Value(contextKey).(*Debugger)
	return dbg
}
