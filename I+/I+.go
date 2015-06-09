package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	acme "github.com/acmescripts"
)

func main() {
	checkError := func(err error, msg string) {
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", msg, err.Error())
			os.Exit(-1)
		}
	}

	var numSpaces int = 2
	var err error = nil

	switch len(os.Args) {
	case 1:
	case 2:
		numSpaces64, err := strconv.ParseInt(os.Args[1], 10 /* base */, 0 /* bits */)
		numSpaces = int(numSpaces64)
		if err == nil && numSpaces < 0 {
			err = errors.New("Negative numSpacess are disallowed.")
		}
		checkError(err, fmt.Sprintf("Can't parse command line argument: %s", os.Args[1]))
	default:
		fmt.Fprintln(os.Stderr, "Usage: %s [numSpaces]", os.Args[0])
	}

	win, err := acme.GetCurrentWindow()
	checkError(err, "Can't get current window")
	defer win.CloseFiles()

	data, err := acme.ReadSelection(win)
	checkError(err, "Can't read dot contents")

	start, _, err := acme.GetDotAddr(win)
	checkError(err, "Can't get dot addr")

	// Indents dot, updates selection.
	data = acme.Indent(data, numSpaces)
	win.Write("data", data)
	win.Addr("#%v,#%v", start, start+len(data))
	win.Ctl("dot=addr")
}
