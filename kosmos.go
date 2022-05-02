package main

import (
	dbg "runtime/debug"
	rt "runtime"
	"fmt"
	"os"

	el "github.com/kosmosJS/engine-node/eventloop"
	en "github.com/kosmosJS/engine"
	s "github.com/kosmosJS/std"
)

var (
	h = []string {
		"usage:",
		"kosmos <argument or file>",
		"",
		"arguments:",
		"help - print this help screen",
		"version - print kosmos version",
	}

	gc = "N/A"
	ver = "0.0.1"
)

func help() {
    for _, l := range h {
        fmt.Println(l)
    }
}

func run(p, d string) error {
	l := el.NewEventLoop(func() {
		s.RegisterAll()
	})

	c, e := en.Compile(p, d, false)

	if e != nil {
		return e
	}

	l.Run(func(v *en.Runtime) {
		_, e = v.RunProgram(c)
	})

	return e
}

func main() {
	a := os.Args[1]

	if len(a) == 0 {
		fmt.Println("must specify a command or file to run\n")
		help()
		os.Exit(1)
	}

	switch(a) {
	case "help":
		help()
		os.Exit(0)
	case "version":
		fmt.Printf("kosmos version %s %s/%s commit %s\n", ver, rt.GOOS, rt.GOARCH, gc)
		os.Exit(0)
	}

	if a[len(a)-3:] != ".js" && a[len(a)-4:] != ".cjs" {
		fmt.Println("file must be a valid CommonJS JavaScript file.")
		os.Exit(1)
	} else {
		_, se := os.Stat(a)
		if se != nil {
			fmt.Println("file cannot be found.")
			os.Exit(1)
		}
	}

	d, fe := os.ReadFile(a)

	if fe != nil {
		fmt.Println("file cannot be accessed.")
		os.Exit(1)
	}

	defer func() {
		if x := recover(); x != nil {
			dbg.Stack()
			panic(x)
		}
	}()

	if e := run(a, string(d)); e != nil {
		switch e := e.(type) {
		case *en.Exception:
			fmt.Println(e.String())
		case *en.InterruptedError:
			fmt.Println(e.String())
		default:
			fmt.Println(e)
		}
		os.Exit(64)
	}
}
