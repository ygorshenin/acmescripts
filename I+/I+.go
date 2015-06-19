package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	acme "github.com/ygorshenin/acmescripts"
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: %s [num-spaces]", os.Args[0])
	os.Exit(-1)
}

func addPrefix(data []byte, numSpaces int) []byte {
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
	dot = addPrefix(dot, numSpaces)
	acme.Write(win, "data", dot)
	acme.WriteAddr(win, "#%v,#%v", start, start+len(dot))
	acme.Ctl(win, "dot=addr")
}
