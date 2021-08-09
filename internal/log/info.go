package log

import (
	"os"
)

var (
	info = NewPrinter(os.Stderr)
	err  = NewPrinter(os.Stderr)
)

func SetInfoEnabled(b bool) {
	info.SetEnabled(b)
}

func SetErrorEnabled(b bool) {
	err.SetEnabled(b)
}

func Infof(format string, a ...interface{}) {
	info.Printf(format, a...)
}

func Info(a ...interface{}) {
	info.Print(a...)
}

func Infoln(a ...interface{}) {
	info.Println(a...)
}

func Errorf(format string, a ...interface{}) {
	err.Printf(format, a...)
}

func Error(a ...interface{}) {
	err.Print(a...)
}

func Errorln(a ...interface{}) {
	err.Println(a...)
}
