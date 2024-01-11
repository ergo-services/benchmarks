package main

import (
	"runtime"
	"sync"
	"time"

	"ergo.services/ergo"
	"ergo.services/ergo/gen"
	"ergo.services/ergo/net/edf"
	"ergo.services/ergo/net/handshake"
	"ergo.services/logger/colored"
)

const (
	TestMessage string = "hi"
)

type sendCase11 struct {
	to gen.PID
	n  int
}

type sendCase1N struct {
	to []gen.PID
	n  int
}

type sendCaseNN struct {
	to []gen.PID
	n  int
}

var (
	wg                 sync.WaitGroup
	n                  uint32
	nodeping, nodepong gen.Node
	sendEvent          gen.Event = gen.Event{Name: "send", Node: "nodeping@localhost"}
	token              gen.Ref
	NCPU               int = runtime.NumCPU()
	POOLSIZE               = NCPU * 2
)

func init() {
	var err error
	options := gen.NodeOptions{}
	options.Network.Cookie = "cookie"
	l := gen.Listener{
		Handshake: handshake.Create(handshake.Options{PoolSize: POOLSIZE}),
	}
	options.Network.Listeners = append(options.Network.Listeners, l)

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

	nodeping, err = ergo.StartNode("nodeping@localhost", options)
	if err != nil {
		panic(err)
	}
	token, err = nodeping.RegisterEvent(sendEvent.Name, gen.EventOptions{})
	if err != nil {
		panic(err)
	}
	nodepong, err = ergo.StartNode("nodepong@localhost", options)
	if err != nil {
		panic(err)
	}

	if err := registerTypes(); err != nil {
		panic(err)
	}

	if _, err := nodeping.Network().GetNode(nodepong.Name()); err != nil {
		panic(err)
	}

}

func startProcesses(a gen.Node, numa int, b gen.Node, numb int) ([]gen.PID, []gen.PID) {
	apids := []gen.PID{}
	for i := 0; i < numa; i++ {
		pid, err := a.Spawn(factory_ping, gen.ProcessOptions{})
		if err != nil {
			panic(err)
		}
		apids = append(apids, pid)
	}
	bpids := []gen.PID{}
	for i := 0; i < numb; i++ {
		pid, err := b.Spawn(factory_pong, gen.ProcessOptions{})
		if err != nil {
			panic(err)
		}
		bpids = append(bpids, pid)
	}

	return apids, bpids
}

func killProcesses(a gen.Node, apids []gen.PID, b gen.Node, bpids []gen.PID) {
	for _, pid := range apids {
		a.Kill(pid)
	}
	for _, pid := range bpids {
		b.Kill(pid)
	}
}

func runTestLocal11() {
	a := nodeping
	apids, bpids := startProcesses(a, 1, a, 1)

	sc := sendCase11{
		to: bpids[0],
		n:  1_000_000,
	}
	time.Sleep(time.Second)
	nodepong.Log().Info("--------------------------------------------------------------------------")
	nodeping.Log().Info("BENCHMARK 1-1: 1 process sends %d messages to 1 process", sc.n)
	if err := nodeping.SendEvent(sendEvent.Name, token, gen.MessageOptions{}, sc); err != nil {
		panic(err)
	}

	start := time.Now()
	wg.Wait()
	elapsed := time.Since(start)

	nodeping.Log().Info("received %d messages. %f msg/sec", sc.n, float64(sc.n)/elapsed.Seconds())
	killProcesses(a, apids, a, bpids)
}

func runTestLocal1N() {
	a := nodeping
	NPROC := NCPU * 2
	apids, bpids := startProcesses(a, 1, a, NPROC)

	sc := sendCase1N{
		to: bpids,
		n:  1_000_000,
	}
	time.Sleep(time.Second)
	nodepong.Log().Info("--------------------------------------------------------------------------")
	nodeping.Log().Info("BENCHMARK 1-N: 1 process sends %d messages to %d processes ",
		sc.n, NPROC)
	if err := nodeping.SendEvent(sendEvent.Name, token, gen.MessageOptions{}, sc); err != nil {
		panic(err)
	}

	start := time.Now()
	wg.Wait()
	elapsed := time.Since(start)

	nodeping.Log().Info("received %d messages. %f msg/sec", sc.n, float64(sc.n)/elapsed.Seconds())
	killProcesses(a, apids, a, bpids)
}

func runTestLocalNN() {
	a := nodeping
	NPROC := NCPU * 8
	apids, bpids := startProcesses(a, NPROC, a, NPROC)

	sc := sendCaseNN{
		to: bpids,
		n:  100_000,
	}
	time.Sleep(time.Second)
	nodepong.Log().Info("--------------------------------------------------------------------------")
	nodeping.Log().Info("BENCHMARK N-N: %d processes send %d messages to %d processes",
		NPROC, sc.n*NPROC, NPROC)
	if err := nodeping.SendEvent(sendEvent.Name, token, gen.MessageOptions{}, sc); err != nil {
		panic(err)
	}

	start := time.Now()
	// nodeping.Log().Info("Started at %d...", start.UnixNano())
	wg.Wait()
	elapsed := time.Since(start)

	nodeping.Log().Info("received %d messages. %f msg/sec", sc.n*NPROC, float64(sc.n*NPROC)/elapsed.Seconds())
	killProcesses(a, apids, a, bpids)
}

func runTestNetwork11() {
	a := nodeping
	b := nodepong
	apids, bpids := startProcesses(a, 1, b, 1)

	sc := sendCase11{
		to: bpids[0],
		n:  1_000_000,
	}
	time.Sleep(time.Second)
	nodepong.Log().Info("--------------------------------------------------------------------------")
	nodeping.Log().Info("BENCHMARK 1-1: 1 process (%s) sends %d messages to 1 process (%s) ",
		nodeping.Name(), sc.n, nodepong.Name())
	if err := nodeping.SendEvent(sendEvent.Name, token, gen.MessageOptions{}, sc); err != nil {
		panic(err)
	}

	start := time.Now()
	wg.Wait()
	elapsed := time.Since(start)

	nodepong.Log().Info("received %d messages. %f msg/sec", sc.n, float64(sc.n)/elapsed.Seconds())
	killProcesses(a, apids, b, bpids)
}

func runTestNetwork1N() {
	a := nodeping
	b := nodepong
	NPROC := NCPU * 2
	apids, bpids := startProcesses(a, 1, b, NPROC)

	sc := sendCase1N{
		to: bpids,
		n:  1_000_000,
	}
	time.Sleep(time.Second)
	nodepong.Log().Info("--------------------------------------------------------------------------")
	nodeping.Log().Info("BENCHMARK 1-N: 1 process (%s) sends %d messages to %d processes (%s) ",
		nodeping.Name(), sc.n, NPROC, nodepong.Name())
	if err := nodeping.SendEvent(sendEvent.Name, token, gen.MessageOptions{}, sc); err != nil {
		panic(err)
	}

	start := time.Now()
	wg.Wait()
	elapsed := time.Since(start)

	nodepong.Log().Info("received %d messages. %f msg/sec", sc.n, float64(sc.n)/elapsed.Seconds())
	killProcesses(a, apids, b, bpids)
}

func runTestNetworkNN() {
	a := nodeping
	b := nodepong
	NPROC := NCPU
	apids, bpids := startProcesses(a, NPROC, b, NPROC*2)

	sc := sendCaseNN{
		to: bpids,
		n:  1_000_000,
	}
	time.Sleep(time.Second)
	nodepong.Log().Info("--------------------------------------------------------------------------")
	nodeping.Log().Info("BENCHMARK N-N: %d processes (%s) send %d messages to %d processes (%s) ",
		NPROC, nodeping.Name(), sc.n*NPROC, NPROC*2, nodepong.Name())
	if err := nodeping.SendEvent(sendEvent.Name, token, gen.MessageOptions{}, sc); err != nil {
		panic(err)
	}

	start := time.Now()
	wg.Wait()
	elapsed := time.Since(start)

	nodepong.Log().Info("received %d messages. %f msg/sec", sc.n*NPROC, float64(sc.n*NPROC)/elapsed.Seconds())
	killProcesses(a, apids, b, bpids)
}

func main() {
	time.Sleep(3 * time.Second)
	nodeping.Log().Info("-------------------------- LOCAL (start) ----------------------------------")
	nodeping.Log().Info("N CPU: %d", NCPU)
	runTestLocal11()
	runTestLocal1N()
	runTestLocalNN()
	nodeping.Log().Info("-------------------------- LOCAL (end) ------------------------------------")

	nodeping.Log().Info("-------------------------- OVER NETWORK (start) ---------------------------")
	nodeping.Log().Info("Network pool: %d TCP-links", POOLSIZE)
	nodeping.Log().Info("N CPU: %d", NCPU)
	time.Sleep(3 * time.Second)
	runTestNetwork11()
	runTestNetwork1N()
	runTestNetworkNN()
	nodeping.Log().Info("-------------------------- OVER NETWORK (stop) ----------------------------")

	nodeping.Wait()
}

func registerTypes() error {
	types := []any{}

	for _, t := range types {
		err := edf.RegisterTypeOf(t)
		if err == nil || err == gen.ErrTaken {
			continue
		}
		panic(err)
	}
	return nil
}
