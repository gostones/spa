package log

import (
	"os"
)

var debug Printer = NewPrinter(os.Stderr)

func SetDebugEnabled(b bool) {
	debug.SetEnabled(b)
}

func Debugf(format string, a ...interface{}) {
	debug.Printf(format, a...)
}

func Debug(a ...interface{}) {
	debug.Print(a...)
}

func Debugln(a ...interface{}) {
	debug.Println(a...)
}
