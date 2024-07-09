package main

import (
	"runtime"
	"time"

	"ergo.services/ergo"
	"ergo.services/ergo/gen"
	"ergo.services/ergo/net/handshake"
	"ergo.services/logger/colored"
	. "github.com/klauspost/cpuid/v2"
)

func runTestNetworkNN() {
	N := 1_000_000
	// prepare nodes
	options := gen.NodeOptions{}
	options.Network.Cookie = "cookie"
	a := gen.AcceptorOptions{
		Handshake: handshake.Create(handshake.Options{PoolSize: NCPU / 2}),
	}
	options.Network.Acceptors = append(options.Network.Acceptors, a)
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

	nodeping, err := ergo.StartNode("node_network_NN_n1@localhost", options)
	if err != nil {
		panic(err)
	}
	nodepong, err := ergo.StartNode("node_network_NN_n2@localhost", options)
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
	nodeping.Log().Info("-------------------------- NETWORK N-N (start) ----------------------------------")
	nodeping.Log().Info("Go Version : %s", runtime.Version())
	nodeping.Log().Info("CPU: %s (Physical Cores: %d)", CPU.BrandName, CPU.PhysicalCores)
	nodeping.Log().Info("Runtime CPUs: %d", NCPU)
	// starting N ping processes
	np := NCPU
	WGready.Add(np)
	for i := 0; i < np; i++ {
		if _, err := nodeping.Spawn(factory_ping_network, gen.ProcessOptions{}, nodepong.Name(), pong); err != nil {
			panic(err)
		}
	}
	nodeping.Log().Info("BENCHMARK: %d processes send %d messages to %d processes", np, np*N, np)
	WGready.Wait() // created monitor on the event and spawned a pong process

	WGready.Add(np)
	if err := nodeping.SendEvent(EVENT.Name, token, gen.MessageOptions{}, startSend{n: N}); err != nil {
		panic(err)
	}
	WGready.Wait() // received event and started sending

	start := time.Now()
	WG.Wait()
	elapsed := time.Since(start)

	nodeping.Log().Info("received %d messages. %f msg/sec", N*np, float64(N*np)/elapsed.Seconds())

	nodeping.Log().Info("-------------------------- NETWORK N-N (end) ----------------------------------")
}
