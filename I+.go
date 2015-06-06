package main

import (
	"errors"
	"fmt"
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

func Indent(data []byte, shift int) []byte {
	spaces := []byte(strings.Repeat(" ", shift))
	result := make([]byte, 0)
	insertSpaces := true
	for _, b := range data {
		if insertSpaces && b != '\n' {
			result = append(result, spaces...)
			insertSpaces = false
		}
		result = append(result, b)
		if b == '\n' {
			insertSpaces = true
		}
	}

	return result
}

func main() {
	checkError := func(err error, msg string) {
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", msg, err.Error())
			os.Exit(-1)
		}
	}

	var shift int = 2
	var err error = nil

	switch len(os.Args) {
	case 1:
	case 2:
		shift64, err := strconv.ParseInt(os.Args[1], 10 /* base */, 0 /* bits */)
		shift = int(shift64)
		if err == nil && shift < 0 {
			err = errors.New("Negative shifts are disallowed.")
		}
		checkError(err, fmt.Sprintf("Can't parse command line argument: %s", os.Args[1]))
	default:
		fmt.Fprintln(os.Stderr, "Usage: %s [shift]", os.Args[0])
	}

	win, err := GetCurrentWindow()
	checkError(err, "Can't get current window")
	defer win.CloseFiles()

	// Following code warms up addr file. See discussion here:
	// https://groups.google.com/forum/#!topic/comp.os.plan9/z4N7eEIW4iw
	_, _, err = win.ReadAddr()
	checkError(err, "Can't warm up addr file")

	data, err := ReadSelection(win)
	checkError(err, "Can't read dot contents")

	start, _, err := GetDotAddr(win)
	checkError(err, "Can't get dot addr")

	// Indents dot, updates selection.
	data = Indent(data, shift)
	win.Write("data", data)
	win.Addr("#%v,#%v", start, start+len(data))
	win.Ctl("dot=addr")
}
