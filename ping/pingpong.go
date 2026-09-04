package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"
	"time"

	"ergo.services/ergo/gen"
)

type startSend struct {
	n int
}

var (
	WGready sync.WaitGroup
	WG      sync.WaitGroup
	EVENT   gen.Event = gen.Event{Name: "send"}
	NCPU    int       = runtime.NumCPU()
)

func main() {
	// f, err := os.Create("profile.prof")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()
	//
	// // Start CPU profiling
	// if err := pprof.StartCPUProfile(f); err != nil {
	// 	panic(err)
	// }
	// defer pprof.StopCPUProfile()
	//
	// // Start tracing
	// traceFile, err := os.Create("trace.out")
	// if err != nil {
	// 	panic(err)
	// }
	// defer traceFile.Close()
	//
	// if err := trace.Start(traceFile); err != nil {
	// 	panic(err)
	// }
	// defer trace.Stop()

	local11 := flag.Bool("local11", false, "run local 1-to-1 benchmark")
	localNN := flag.Bool("localNN", false, "run local N-to-N benchmark")
	network11 := flag.Bool("network11", false, "run network 1-to-1 benchmark")
	networkNN := flag.Bool("networkNN", false, "run network N-to-N benchmark")
	flag.Parse()

	// no flags selected -> run all
	all := (*local11 || *localNN || *network11 || *networkNN) == false

	var tests []func()
	if *local11 || all {
		tests = append(tests, runTestLocal11)
	}
	if *localNN || all {
		tests = append(tests, runTestLocalNN)
	}
	if *network11 || all {
		tests = append(tests, runTestNetwork11)
	}
	if *networkNN || all {
		tests = append(tests, runTestNetworkNN)
	}

	for i, t := range tests {
		if i > 0 {
			fmt.Println("")
			time.Sleep(time.Second * 10)
		}
		t()
	}
}
