package log

import (
	"fmt"
	"io"
)

// Printer represents a simple logger interface
type Printer interface {
	SetEnabled(bool)
	Printf(string, ...interface{})
	Print(...interface{})
	Println(...interface{})
}

// NewPrinter creates a logger with the specified writer.
func NewPrinter(w io.Writer) Printer {
	return &printer{
		out: w,
		on:  true,
	}
}

type printer struct {
	out io.Writer
	on  bool
}

func (r *printer) SetEnabled(b bool) {
	r.on = b
}

func (r *printer) Printf(format string, a ...interface{}) {
	if r.on {
		fmt.Fprintf(r.out, format, a...)
	}
}

func (r *printer) Print(a ...interface{}) {
	if r.on {
		fmt.Fprint(r.out, a...)
	}
}

func (r *printer) Println(a ...interface{}) {
	if r.on {
		fmt.Fprintln(r.out, a...)
	}
}
