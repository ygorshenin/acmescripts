package main

import (
	"bytes"
	"sort"

	acme "github.com/ygorshenin/acmescripts"
)

type ByteSlice [][]byte

func (b ByteSlice) Len() int {
	return len(b)
}

func (b ByteSlice) Less(i, j int) bool {
	lhs, rhs := b[i], b[j]
	for k := 0; k < len(lhs) && k < len(rhs); k++ {
		if lhs[k] < rhs[k] {
			return true
		}
		if lhs[k] > rhs[k] {
			return false
		}
	}
	return false
}

func (b ByteSlice) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func main() {
	win := acme.GetCurrentWindow()
	defer win.CloseFiles()

	lines := bytes.Split(acme.ReadDot(win), []byte("\n"))
	start, _ := acme.ReadAddr(win)

	sort.Sort(ByteSlice(lines))
	var result []byte
	for _, line := range lines {
		result = append(result, line...)
		result = append(result, '\n')
	}

	result = result[0 : len(result)-1]

	acme.Write(win, "data", result)
	acme.WriteAddr(win, "#%v,#%v", start, start+len(result))
	acme.Ctl(win, "dot=addr")
}
