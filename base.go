// This files contains checked wrappers around acme interface.
package acmescripts

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"9fans.net/go/acme"
)

func GetCurrentWindow() *acme.Win {
	env := os.Getenv("winid")
	if len(env) == 0 {
		CheckError(errors.New("can't get winid env."), "GetCurrentWindow")
	}

	id, err := strconv.ParseInt(env, 10 /* base */, 0 /* bits */)
	CheckError(err, "GetCurrentWindow: can't parse winid", env)

	win, err := acme.Open(int(id), nil)
	CheckError(err, "GetCurrentWindow: can't open current window")

	// Following code warms up addr file. See discussion here:
	// https://groups.google.com/forum/#!topic/comp.os.plan9/z4N7eEIW4iw
	_, _, err = win.ReadAddr()
	CheckError(err, "GetCurrentWindow: can't warm-up addr file")

	return win
}

func ReadAddr(win *acme.Win) (int, int) {
	Ctl(win, "addr=dot")
	q0, q1, err := win.ReadAddr()
	CheckError(err, "ReadAddr")
	return q0, q1
}

func WriteAddr(win *acme.Win, format string, args ...interface{}) {
	addr := fmt.Sprintf(format, args...)
	CheckError(win.Addr(addr), "WriteAddr", addr)
}

func Ctl(win *acme.Win, format string, args ...interface{}) {
	ctl := fmt.Sprintf(format, args...)
	CheckError(win.Ctl(ctl), "Ctl", ctl)
}

func ReadFileName(win *acme.Win) []byte {
	tag := Read(win, "tag")
	if tag == nil {
		return nil
	}
	fields := bytes.Fields(tag)
	if len(fields) == 0 {
		return nil
	}
	return fields[0]
}

func Read(win *acme.Win, file string) []byte {
	content, err := win.ReadAll(file)
	CheckError(err, "Read", file)
	return content
}

func Write(win *acme.Win, file string, content []byte) {
	n, err := win.Write(file, content)
	if err == nil && n != len(content) {
		err = errors.New(fmt.Sprintf("bytes written mismatch: expected: %v, actual: %v", len(content), n))
	}
	CheckError(err, "Write", file)
}

func ReadDot(win *acme.Win) []byte {
	q0, q1 := ReadAddr(win)
	dot := make([]byte, q1-q0)
	if len(dot) == 0 {
		return dot
	}
	n, err := win.Read("data", dot)
	if err == nil && n != len(dot) {
		err = errors.New(fmt.Sprintf("bytes read mismatch: expected: %v, actual: %v", len(dot), n))
	}
	CheckError(err, "ReadDot")
	return dot
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

func CheckError(err error, msg ...interface{}) {
	if err != nil {
		fmt.Fprint(os.Stderr, msg...)
		fmt.Fprint(os.Stderr, ":", err.Error)
		os.Exit(-1)
	}
}
