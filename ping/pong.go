package main

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

func factory_pong() gen.ProcessBehavior {
	return &pong{}
}

type pong struct {
	act.Actor
}

func (p *pong) HandleMessage(from gen.PID, message any) error {
	WG.Done()
	return nil
}
