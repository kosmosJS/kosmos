package main

import (
    crand "crypto/rand"
    "encoding/binary"
	"runtime/debug"
    "math/rand"
	"bufio"
	"fmt"
    "os"

	"github.com/kosmosJS/engine-node/eventloop"
	"github.com/kosmosJS/kosmosJS/src/utils"
	"github.com/kosmosJS/engine"
)

func newRandSource() engine.RandSource {
	var seed int64
	if err := binary.Read(crand.Reader, binary.LittleEndian, &seed); err != nil {
		panic(fmt.Errorf("Could not read random bytes: %v", err))
	}
	return rand.New(rand.NewSource(seed)).Float64
}

func run(p string, d []string) error {
	loop := eventloop.NewEventLoop()

	prg, err := engine.Compile(p, utils.SliceToMLString(d), false)

	if err != nil {
		return err
	}

	loop.Run(func(vm *engine.Runtime) {
		_, err = vm.RunProgram(prg)
	})

	return err
}

func main() {
	a := os.Args[1]
	
	if len(a) == 0 {
		fmt.Println("must specify a command or file to run\n")
		utils.PrintHelp()
		os.Exit(1)
	}

	if utils.Contains(utils.Arguments, a) {
		switch(a) {
			case "help":
				utils.PrintHelp()
			case "version":
				utils.PrintVersion()
		}
		os.Exit(0)
	}

	if !utils.Contains([]string{".js", "cjs"}, a[len(a)-3:]) {
		fmt.Println("file must be a valid CommonJS JavaScript file.")
		os.Exit(1)
	} else {
		_, se := os.Stat(a)
		if se != nil {
			fmt.Println("file cannot be found.")
			os.Exit(1)
		}
	}
	
	var d []string

	f, fe := os.Open(a)

	if fe == nil {
		sc := bufio.NewScanner(f)
		sc.Split(bufio.ScanLines)
		for sc.Scan() {
			d = append(d, sc.Text())
		}
	} else {
		fmt.Println("file cannot be accessed.")
		os.Exit(1)
	}

	f.Close()

	defer func() {
		if x := recover(); x != nil {
			debug.Stack()
			panic(x)
		}
	}()

	if err := run(a, d); err != nil {
		switch err := err.(type) {
		case *engine.Exception:
			fmt.Println(err.String())
		case *engine.InterruptedError:
			fmt.Println(err.String())
		default:
			fmt.Println(err)
		}
		os.Exit(64)
	}
}
