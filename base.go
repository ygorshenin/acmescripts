package acmescripts

import (
	"bytes"
	"errors"
	"os"
	"strconv"
	"strings"

	"9fans.net/go/acme"
)

// Returns descriptor of the current Acme window.
func GetCurrentWindow() (win *acme.Win, err error) {
	env := os.Getenv("winid")
	if len(env) == 0 {
		err = errors.New("Can't get winid env.")
		return
	}

	id, err := strconv.ParseInt(env, 10 /* base */, 0 /* bits */)
	if err != nil {
		return
	}

	win, err = acme.Open(int(id), nil)

	if err != nil {
		return
	}

	// Following code warms up addr file. See discussion here:
	// https://groups.google.com/forum/#!topic/comp.os.plan9/z4N7eEIW4iw
	_, _, err = win.ReadAddr()
	if err != nil {
		return
	}

	return
}

// Returns addr of the dot.
func GetDotAddr(win *acme.Win) (q0, q1 int, err error) {
	err = win.Ctl("addr=dot")
	if err != nil {
		return
	}
	return win.ReadAddr()
}

// Returns bytes of the dot.
func ReadSelection(win *acme.Win) (content []byte, err error) {
	q0, q1, err := GetDotAddr(win)
	if err != nil {
		return
	}

	content = make([]byte, q1-q0)
	if len(content) == 0 {
		return
	}

	bytesRead, err := win.Read("data", content)
	if err != nil {
		return
	}
	content = content[0:bytesRead]
	return
}

func Indent(data []byte, numSpaces int) []byte {
	if len(data) == 0 {
		return data
	}

	var prefix []byte = []byte(strings.Repeat(" ", numSpaces))
	var sep []byte = []byte("\n")

	var result []byte = make([]byte, 0)
	for _, line := range bytes.Split(data, sep) {
		if len(line) > 0 {
			result = append(result, prefix...)
			result = append(result, line...)
		}
		result = append(result, '\n')
	}
	return result[0 : len(result)-1]
}
