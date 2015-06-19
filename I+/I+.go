package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	acme "github.com/ygorshenin/acmescripts"
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: %s [num-spaces]", os.Args[0])
	os.Exit(-1)
}

func main() {
	var numSpaces int = 2

	switch len(os.Args) {
	case 1:
	case 2:
		numSpaces64, err := strconv.ParseInt(os.Args[1], 10 /* base */, 0 /* bits */)
		numSpaces = int(numSpaces64)
		if err == nil && numSpaces < 0 {
			err = errors.New("Negative numSpaces are disallowed.")
		}
		acme.CheckError(err, "Can't parse command line argument:", os.Args[1])
	default:
		usage()
	}

	win := acme.GetCurrentWindow()
	defer win.CloseFiles()

	dot := acme.ReadDot(win)
	start, _ := acme.ReadAddr(win)

	// Indents dot, updates selection.
	dot = acme.Indent(dot, numSpaces)
	acme.Write(win, "data", dot)
	acme.WriteAddr(win, "#%v,#%v", start, start+len(dot))
	acme.Ctl(win, "dot=addr")
}
