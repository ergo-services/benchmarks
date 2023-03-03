package main

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
)

type pong struct {
	gen.Server
}

type pongState struct {
	done *sync.WaitGroup
}

func (p *pong) Init(process *gen.ServerProcess, args ...etf.Term) error {
	process.State = &pongState{
		done: args[0].(*sync.WaitGroup),
	}
	return nil
}

func (p *pong) HandleInfo(process *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	state := process.State.(*pongState)
	state.done.Done()
	return gen.ServerStatusOK
}

// pongTest
type pongTest struct {
	gen.Server
}
type PongTestMessage struct {
	I     int
	Value interface{}
}

func (p *pongTest) Init(process *gen.ServerProcess, args ...etf.Term) error {
	process.State = &pongState{
		done: args[0].(*sync.WaitGroup),
	}
	return nil
}

func (p *pongTest) HandleInfo(process *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	state := process.State.(*pongState)
	m := message.(PongTestMessage)
	if !reflect.DeepEqual(table[m.I].value, m.Value) {
		fmt.Printf("MISMATCH      (got): %#v\n", m.Value)
		fmt.Printf("MISMATCH (expected): %#v\n", table[m.I].value)
		panic("mismatch")
	}
	state.done.Done()
	return gen.ServerStatusOK
}
