package printer

import (
	"fmt"
	"io"
	"strings"

	"github.com/MarcinKonowalczyk/trace"
)

// helper to write trace output with indentation
type writer struct {
	io.Writer
}

func (w *writer) W(symbol rune, indent int, name string) {
	f := "%s%c %s\n"
	msg := fmt.Sprintf(f, strings.Repeat(".", indent*2), symbol, name)
	w.Write([]byte(msg))
}

func NewTracePrinter(
	out io.Writer,
	print_messages bool,
) func(trace.Node) error {
	if out == nil {
		// no output. nothing to do.
		return func(n trace.Node) error { return nil }
	}
	indent := 0
	w := writer{out}
	return func(n trace.Node) error {
		switch n := n.(type) {
		case *trace.Enter:
			if n.Name() == trace.START_NODE {
				return nil
			}
			w.W('>', indent, n.Name())
			indent++
		case *trace.Exit:
			if n.Name() == trace.END_NODE {
				return nil
			}
			indent--
			w.W('<', indent, n.Name())
		case *trace.Message:
			if print_messages {
				w.W('@', indent, n.Message)
			}
		default:
			panic(fmt.Sprintf("unknown node type: %T", n))
		}
		return nil
	}
}

func SprintTrace(
	tracer trace.Tracer,
	print_messages bool,
) string {
	tracer.Done()
	walkable, err := tracer.ToWalkable()
	if err != nil {
		panic(err)
	}
	var out strings.Builder
	_ = walkable.Walk(NewTracePrinter(&out, print_messages))
	// i.stdout.Write([]byte(out.String()))
	return out.String()
}

func PrintTrace(
	tracer trace.Tracer,
	print_messages bool,
	writer io.Writer,
) {
	writer.Write([]byte(SprintTrace(tracer, print_messages)))
}
