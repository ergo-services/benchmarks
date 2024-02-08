package main

import (
	"runtime"
	"time"

	"ergo.services/ergo"
	"ergo.services/ergo/gen"
	"ergo.services/logger/colored"
	. "github.com/klauspost/cpuid/v2"
)

func runTestNetwork11() {
	N := 3_000_000
	// prepare nodes
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

	nodeping, err := ergo.StartNode("node_network_11_n1@localhost", options)
	if err != nil {
		panic(err)
	}
	nodepong, err := ergo.StartNode("node_network_11_n2@localhost", options)
	if err != nil {
		panic(err)
	}

	if _, err := nodeping.Network().GetNode(nodepong.Name()); err != nil {
		panic(err)
	}

	pong := gen.Atom("pong")
	nodepong.Network().EnableSpawn(pong, factory_pong)

	token, err := nodeping.RegisterEvent(EVENT.Name, gen.EventOptions{})
	if err != nil {
		panic(err)
	}
	nodeping.Log().Info("-------------------------- NETWORK 1-1 (start) ----------------------------------")
	nodeping.Log().Info("Go Version : %s", runtime.Version())
	nodeping.Log().Info("CPU: %s (Physical Cores: %d)", CPU.BrandName, CPU.PhysicalCores)
	nodeping.Log().Info("Runtime CPUs: %d", NCPU)
	// starting 1 ping process
	WGready.Add(1)
	if _, err := nodeping.Spawn(factory_ping_network, gen.ProcessOptions{}, nodepong.Name(), pong); err != nil {
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

	nodeping.Log().Info("-------------------------- NETWORK 1-1 (end) ----------------------------------")
}
