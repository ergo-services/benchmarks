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
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/ergo-services/ergo"
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
	"github.com/ergo-services/ergo/lib"
	"github.com/ergo-services/ergo/node"
)

type MapStringInt map[string]int
type SliceInt []int
type StructFlat struct {
	String string
	Int8   int8
	Int    int
	Int64  int64
	Float  float64
	PID    etf.Pid
	Alias  etf.Alias
	Ref    etf.Ref
}

type StructNest struct {
	StructFlat
	Map   MapStringInt
	Slice SliceInt
}

type testCase struct {
	name  string
	value any
}

var (
	nodePing, nodePong node.Node
	err                error

	table []testCase

	n uint32
)

func init() {
	nodePing, err = ergo.StartNode("nodePing@localhost", "cookies", node.Options{})
	if err != nil {
		panic(err)
	}
	nodePong, err = ergo.StartNode("nodePong@localhost", "cookies", node.Options{})
	if err != nil {
		panic(err)
	}

	if err := registerTypes(); err != nil {
		panic(err)
	}

	if err := nodePing.Connect(nodePong.Name()); err != nil {
		panic(err)
	}

	mapStringInt10 := make(MapStringInt)
	mapStringInt100 := make(MapStringInt)
	mapStringInt1000 := make(MapStringInt)

	sliceInt10 := make(SliceInt, 10)
	sliceInt100 := make(SliceInt, 100)
	sliceInt1000 := make(SliceInt, 1000)

	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("item %d", i)
		if i < 10 {
			mapStringInt10[key] = i
			sliceInt10[i] = i
		}
		if i < 100 {
			mapStringInt100[key] = i
			sliceInt100[i] = i
		}

		mapStringInt1000[key] = i
		sliceInt1000[i] = i
	}

	//structFlat := StructFlat{
	//	String: "Hello World!",
	//	Int8:   127,
	//	Int:    2147483647,
	//	Int64:  4294967296,
	//	Float:  3.14,
	//	PID:    etf.Pid{Node: etf.Atom("bench@127.0.0.1"), ID: 32767, Creation: 2},
	//	Ref:    etf.Ref{Node: etf.Atom("bench@127.0.0.1"), Creation: 8, ID: [5]uint32{73444, 3082813441, 2373634851, 0, 0}},
	//	Alias:  etf.Alias{Node: etf.Atom("bench@127.0.0.1"), Creation: 8, ID: [5]uint32{73444, 3082813441, 2373634851, 0, 0}},
	//}

	table = []testCase{
		{"atom", etf.Atom("Hello World!")},
		//{"string", "Hello World!"},
		//{"int8", 127},
		//{"int64", int64(2147483647)},
		//{"float", 3.14},
		//{"PID", etf.Pid{Node: etf.Atom("bench@127.0.0.1"), ID: 32767, Creation: 2}},
		//{"Ref", etf.Ref{Node: etf.Atom("bench@127.0.0.1"), Creation: 8, ID: [5]uint32{73444, 3082813441, 2373634851, 0, 0}}},
		//{"MapStringInt10", mapStringInt10},
		//{"MapStringInt100", mapStringInt100},
		//{"MapStringInt1000", mapStringInt1000},
		//{"SliceInt10", sliceInt10},
		//{"SliceInt100", sliceInt100},
		//{"SliceInt1000", sliceInt1000},
		//{"StructFlat", structFlat},
		//{"StructNest", StructNest{StructFlat: structFlat, Map: mapStringInt10, Slice: sliceInt10}},
	}
}

func TestNode(t *testing.T) {
	var wg sync.WaitGroup
	opts := gen.ProcessOptions{
		//	MailboxSize: 10000,
	}
	pingProcess, err := nodePing.Spawn("ping", opts, &ping{}, nil)
	if err != nil {
		panic(err)
	}
	//defer pingProcess.Kill()

	pongProcess, err := nodePong.Spawn("pong", opts, &pongTest{}, &wg)
	if err != nil {
		panic(err)
	}
	//defer pongProcess.Kill()

	for i, tc := range table {
		wg.Add(1)
		if err := pingProcess.Send(pongProcess.Self(), PongTestMessage{I: i, Value: tc.value}); err != nil {
			panic("can't send")
		}
		wg.Wait()
	}
}

func BenchmarkNodeNetworkMessaging(b *testing.B) {

	for _, tc := range table {
		b.Run(tc.name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				var wg sync.WaitGroup

				nn := atomic.AddUint32(&n, 1)
				pingProcessName := fmt.Sprintf("ping%d", nn)
				opts := gen.ProcessOptions{
					MailboxSize: 10000,
				}
				pingProcess, err := nodePing.Spawn(pingProcessName, opts, &ping{}, nil)
				if err != nil {
					panic(err)
				}
				//defer pingProcess.Kill()

				nn = atomic.AddUint32(&n, 1)
				pongProcessName := fmt.Sprintf("pong%d", nn)
				pongProcess, err := nodePong.Spawn(pongProcessName, opts, &pong{}, &wg)
				if err != nil {
					panic(err)
				}
				//defer pongProcess.Kill()

				b.ResetTimer()
				for pb.Next() {
					wg.Add(1)
					if err := pingProcess.Send(pongProcess.Self(), tc.value); err != nil {
						panic(err)
					}
				}
				wg.Wait()
			})
		})
	}

}
func BenchmarkNodeLocalMessaging(b *testing.B) {

	tc := table[2] // use the first case only
	cpus := []int{4, 16, 64}

	for _, cpu := range cpus {
		runtime.GOMAXPROCS(cpu)
		//fmt.Println(n)
		b.Run(tc.name+"CPU"+strconv.Itoa(cpu), func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				var wg sync.WaitGroup

				nn := atomic.AddUint32(&n, 1)
				pingProcessName := fmt.Sprintf("ping%d", nn)
				opts := gen.ProcessOptions{
					//				MailboxSize: 20000,
				}
				pingProcess, err := nodePing.Spawn(pingProcessName, opts, &ping{}, nil)
				if err != nil {
					panic(err)
				}
				//defer pingProcess.Kill()

				nn = atomic.AddUint32(&n, 1)
				pongProcessName := fmt.Sprintf("pong%d", nn)
				pongProcess, err := nodePing.Spawn(pongProcessName, opts, &pong{}, &wg)
				if err != nil {
					panic(err)
				}
				//defer pongProcess.Kill()

				b.ResetTimer()
				for pb.Next() {
					wg.Add(1)
					if err := pingProcess.Send(pongProcess.Self(), tc.value); err != nil {
						panic(err)
					}
				}
				wg.Wait()
			})
		})
	}

}

func registerTypes() error {
	types := []interface{}{
		MapStringInt{},
		SliceInt{},
		StructFlat{},
		StructNest{},
		PongTestMessage{},
	}

	opts := etf.RegisterTypeOptions{Strict: true}

	for _, t := range types {
		if _, err := etf.RegisterType(t, opts); err != nil && err != lib.ErrTaken {
			return err
		}
	}
	return nil
}
