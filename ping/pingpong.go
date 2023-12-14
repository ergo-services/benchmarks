// to run bench:
//   $ go test -bench . | tee ./output.version1.txt
// repeat it for the another version
//
// if 'benchstat' is not installed:
//   $ go install golang.org/x/perf/cmd/benchstat@latest
//
// to see the diff:
//   $ benchstat ./output.version1.txt ./output.version2.txt

package main

import (
	"sync"
	"time"

	"ergo.services/ergo"
	"ergo.services/ergo/gen"
	"ergo.services/ergo/lib"
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

var (
	wg                 sync.WaitGroup
	n                  uint32
	nodeping, nodepong gen.Node
	sendEvent          gen.Event = gen.Event{Name: "send", Node: "nodeping@localhost"}
	token              gen.Ref
)

func init() {
	var err error
	options := gen.NodeOptions{}
	// options.Log.Level = gen.LogLevelTrace
	options.Network.Cookie = "cookie"
	l := gen.Listener{
		Handshake: handshake.Create(handshake.Options{PoolSize: 15}),
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

func runTestNetwork11() {
	a := nodeping
	b := nodepong
	apids, bpids := startProcesses(a, 1, b, 1)

	sc := sendCase11{
		to: bpids[0],
		n:  10_000_000,
	}
	time.Sleep(time.Second)
	nodepong.Log().Info("-------------------------------------------------------------------")
	nodeping.Log().Info("BENCHMARK 1-1: 1 process (%s) sends %d messages to 1 process (%s) ",
		nodeping.Name(), sc.n, nodepong.Name())
	if err := nodeping.SendEvent(sendEvent.Name, token, gen.MessageOptions{}, sc); err != nil {
		panic(err)
	}

	start := time.Now()
	wg.Wait()
	elapsed := time.Since(start)

	nodepong.Log().Info("received %d messages. %v msg/sec", sc.n, float64(sc.n)/elapsed.Seconds())
	nodepong.Log().Info("-------------------------------------------------------------------")
	killProcesses(a, apids, b, bpids)
}

func runTestNetwork1N() {
	a := nodeping
	b := nodepong
	apids, bpids := startProcesses(a, 1, b, 100)

	sc := sendCase1N{
		to: bpids,
		n:  10_000_000,
	}
	time.Sleep(time.Second)
	nodepong.Log().Info("-------------------------------------------------------------------")
	nodeping.Log().Info("BENCHMARK 1-N: 1 process (%s) sends %d messages to 100 processes (%s) ",
		nodeping.Name(), sc.n, nodepong.Name())
	if err := nodeping.SendEvent(sendEvent.Name, token, gen.MessageOptions{}, sc); err != nil {
		panic(err)
	}

	start := time.Now()
	wg.Wait()
	elapsed := time.Since(start)

	nodepong.Log().Info("received %d messages. %v msg/sec", sc.n, float64(sc.n)/elapsed.Seconds())
	nodepong.Log().Info("-------------------------------------------------------------------")
	killProcesses(a, apids, b, bpids)
}

func runTestNetworkNN() {
	a := nodeping
	b := nodepong
	apids, bpids := startProcesses(a, 10, b, 10)

	sc := sendCase1N{
		to: bpids,
		n:  200_000,
	}
	time.Sleep(time.Second)
	nodepong.Log().Info("-------------------------------------------------------------------")
	nodeping.Log().Info("BENCHMARK N-N: 10 processes (%s) sends %d messages to 10 processes (%s) ",
		nodeping.Name(), sc.n, nodepong.Name())
	if err := nodeping.SendEvent(sendEvent.Name, token, gen.MessageOptions{}, sc); err != nil {
		panic(err)
	}

	start := time.Now()
	nodeping.Log().Info("Started at %d...", start.UnixNano())
	wg.Wait()
	elapsed := time.Since(start)

	nodepong.Log().Info("received %d messages. %v msg/sec", sc.n*10, float64(sc.n*10)/elapsed.Seconds())
	nodepong.Log().Info("-------------------------------------------------------------------")
	killProcesses(a, apids, b, bpids)
}
func main() {

	lib.StatBuffers()
	nodeping.Log().Info("-------------------------- OVER NETWORK ---------------------------")
	time.Sleep(3 * time.Second)
	// runTestNetwork11()
	// runTestNetwork1N()
	runTestNetworkNN()

	lib.StatBuffers()
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
