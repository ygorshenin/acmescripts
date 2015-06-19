package main

import (
	"bytes"
	"fmt"
	"os/exec"

	acme "github.com/ygorshenin/acmescripts"
)

func main() {
	win := acme.GetCurrentWindow()
	defer win.CloseFiles()

	var args []string = make([]string, 0)

	fileName := acme.ReadFileName(win)
	if fileName != nil {
		args = append(args, fmt.Sprintf("-assume-filename=%v", string(fileName)))
	}

	from, to := acme.ReadAddr(win)
	args = append(args, fmt.Sprintf("-offset=%v", from))
	args = append(args, fmt.Sprintf("-length=%v", to-from))

	var in = acme.Read(win, "body")
	var out bytes.Buffer
	cmd := exec.Command("clang-format", args...)
	cmd.Stdin = bytes.NewReader(in)
	cmd.Stdout = &out
	acme.CheckError(cmd.Run(), "Can't run clang-format")

	to = out.Len() - (len(in) - to)
	acme.Write(win, "data", out.Bytes()[from:to])
	acme.WriteAddr(win, "#%v,#%v", from, to)
	acme.Ctl(win, "dot=addr")
}
