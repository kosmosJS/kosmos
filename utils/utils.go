package utils

import (
    "runtime"
    "fmt"
)

var (
	Arguments = []string {
		"help",
		"version",
	}

	Help = []string {
		"usage:",
		"kosmos <argument or file>",
		"",
		"arguments:",
		"help - print this help screen",
		"version - print void version",
	}

	Commit = "N/A"
)

const (
    Version = "0.0.1"
)

func PrintVersion() {
    fmt.Printf("kosmos version %s %s/%s commit %s\n", Version, runtime.GOOS, runtime.GOARCH, Commit)
    return
}

func PrintHelp() {
    for _, h := range Help {
        fmt.Println(h)
    }
    return
}

func Contains(s []string, v interface{}) bool {
    for _, cv := range s {
        if cv == v {
            return true
        }
    }
    return false
}

func SliceToMLString(s []string) string {
	var f string
	for _, l := range s {
		f += l+`
`
	}

	return f
}
