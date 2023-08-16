package main

import (
	"runtime"
	"time"

	"ergo.services/ergo"
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
	"ergo.services/logger/colored"
)

func factory_simple() gen.ProcessBehavior {
	return &simple{}
}

type simple struct {
	act.Actor
}

func main() {
	var mstat runtime.MemStats
	var options gen.NodeOptions

	options.Log.DefaultLogger.Disable = true
	loggercolored := gen.Logger{
		Name:   "colored",
		Logger: colored.CreateLogger(colored.Options{ShortTimestamp: true, ShortLevelNames: true}),
	}
	options.Log.Loggers = append(options.Log.Loggers, loggercolored)

	node, err := ergo.StartNode("demo@localhost", options)
	if err != nil {
		panic(err)
	}

	mem := func() {
		runtime.ReadMemStats(&mstat)
		node.Log().Info("memory usage: %d", mstat.Alloc)
		runtime.GC()
	}

	mem()
	node.Log().Info("starting 1M processes...")
	for i := 0; i < 1000000; i++ {
		if _, err := node.Spawn(factory_simple, gen.ProcessOptions{}); err != nil {
			panic(err)
		}
	}
	node.Log().Info("1M processes is started")

	for {
		mem()
		time.Sleep(time.Second)
		info, err := node.Info()
		if err != nil {
			panic(err)
		}
		node.Log().Info("total processes: %d", info.ProcessesTotal)
	}

	node.Wait()

}
