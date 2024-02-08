package main

import (
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

	runTestLocal11()
	fmt.Println("")
	time.Sleep(time.Second * 10)
	runTestLocalNN()
	fmt.Println("")
	time.Sleep(time.Second * 10)
	runTestNetwork11()
	fmt.Println("")
	time.Sleep(time.Second * 10)
	runTestNetworkNN()
}
