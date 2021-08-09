package log

import (
	"bufio"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

var p = NewPrinter(os.Stderr)

func SetPromptEnabled(b bool) {
	p.SetEnabled(b)
}

func Promptf(format string, a ...interface{}) {
	p.Printf(format, a...)
}

func Prompt(a ...interface{}) {
	p.Print(a...)
}

func Promptln(a ...interface{}) {
	p.Println(a...)
}

func prompt(ps string, choices map[string]string, r *bufio.Reader) (string, bool, error) {
	Promptf(ps)

	v, err := r.ReadString('\n')
	if err != nil {
		return "", false, err
	}

	c, ok := choices[strings.ToLower(strings.TrimSpace(v))]
	return c, ok, nil
}

func Confirm(ps string) (string, error) {
	r := bufio.NewReader(os.Stdin)

	choices := map[string]string{"y": "y", "yes": "y", "n": "n", "no": "n"}
	for {
		choice, ok, err := prompt(ps, choices, r)
		if err != nil {
			return "", err
		}
		if ok {
			return choice, nil
		}
	}
}

func PromptSecret(ps string) (string, error) {
	Promptf(ps)

	b, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	Promptln()
	return string(b), nil
}
