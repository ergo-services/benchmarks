package main

import (
	"ergo.services/ergo"
	"ergo.services/ergo/gen"
	"ergo.services/logger/colored"
	"time"
)

func runTestLocal11() {
	N := 1_000_000
	// prepare node
	options := gen.NodeOptions{}
	options.Network.Cookie = "cookie"
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

	nodeping, err := ergo.StartNode("node_local_11@localhost", options)
	if err != nil {
		panic(err)
	}

	token, err := nodeping.RegisterEvent(EVENT.Name, gen.EventOptions{})
	if err != nil {
		panic(err)
	}
	nodeping.Log().Info("-------------------------- LOCAL 1-1 (start) ----------------------------------")
	nodeping.Log().Info("N CPU: %d", NCPU)
	// starting 1 ping process
	WGready.Add(1)
	if _, err := nodeping.Spawn(factory_ping_local, gen.ProcessOptions{}); err != nil {
		panic(err)
	}
	nodeping.Log().Info("BENCHMARK: 1 process sends %d messages to 1 process", N)
	WGready.Wait() // created monitor on the event and spawned a pong process

	WGready.Add(1)
	if err := nodeping.SendEvent(EVENT.Name, token, gen.MessageOptions{}, startSend{n: N}); err != nil {
		panic(err)
	}
	WGready.Wait() // received event and started sending

	start := time.Now()
	WG.Wait()
	elapsed := time.Since(start)

	nodeping.Log().Info("received %d messages. %f msg/sec", N, float64(N)/elapsed.Seconds())

	nodeping.Log().Info("-------------------------- LOCAL 1-1 (end) ----------------------------------")
}
