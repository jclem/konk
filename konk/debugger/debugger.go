package debugger

import (
	"context"
	"fmt"
	"os"

	"github.com/kr/pretty"
	"github.com/spf13/cobra"
)

type contextKey struct{}

type Debugger struct {
	debug *bool
}

func WithDebugger(ctx context.Context, debug *bool) context.Context {
	dbg := Debugger{debug}
	return context.WithValue(ctx, contextKey{}, &dbg)
}

func (d *Debugger) Debugf(format string, args ...interface{}) {
	if *d.debug {
		fmt.Fprintf(os.Stdout, fmt.Sprintf("DEBUG: %s\n", format), args...)
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
	dbg, ok := ctx.Value(contextKey{}).(*Debugger)
	if !ok {
		panic("no debugger in context")
	}

	return dbg
}
