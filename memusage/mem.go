package main

import (
	"runtime"
	"time"

	"ergo.services/ergo"
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
	"ergo.services/logger/colored"
	. "github.com/klauspost/cpuid/v2"
)

var (
	NCPU int = runtime.NumCPU()
)

func factory_simple() gen.ProcessBehavior {
	return &simple{}
}

type simple struct {
	act.Actor
}

func main() {
	var options gen.NodeOptions

	options.Network.Cookie = "123"
	loggercolored, err := colored.CreateLogger(colored.Options{
		TimeFormat: time.DateTime,
	})
	if err != nil {
		panic(err)
	}
	options.Log.DefaultLogger.Disable = true
	options.Log.Loggers = append(
		options.Log.Loggers,
		gen.Logger{Name: "colored", Logger: loggercolored},
	)
	node, err := ergo.StartNode("mem@localhost", options)
	if err != nil {
		panic(err)
	}

	mem := func(proc bool) {
		info, _ := node.Info()
		node.Log().Info("Memory allocated (runtime): %.2f Kb", float64(info.MemoryAlloc)/1024.0)
		node.Log().Info("Memory used (OS): %.2f Kb", float64(info.MemoryUsed)/1024.0)
		if proc {
			node.Log().Info("Total processes: %d (memory per process ~%.2f Kb)",
				info.ProcessesTotal,
				(float64(info.MemoryAlloc)/float64(info.ProcessesTotal))/1024.0)
		}
		runtime.GC()
	}
	node.Log().Info("-------------------------- Memory usage (start) ----------------------------------")
	node.Log().Info("Go Version : %s", runtime.Version())
	node.Log().Info("CPU: %s (Physical Cores: %d)", CPU.BrandName, CPU.PhysicalCores)
	node.Log().Info("Runtime CPUs: %d", NCPU)

	mem(false)
	node.Log().Info("Starting 1M processes...")
	start := time.Now()
	for i := 0; i < 1000000; i++ {
		if _, err := node.Spawn(factory_simple, gen.ProcessOptions{}); err != nil {
			panic(err)
		}
	}
	node.Log().Info("1M processes is started. Elapsed: %s", time.Since(start))

	for i := 0; i < 3; i++ {
		mem(true)
		time.Sleep(time.Second)
	}
	node.Log().Info("-------------------------- Memory usage (end) ----------------------------------")

}
