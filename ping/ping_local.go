package main

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

func factory_ping_local() gen.ProcessBehavior {
	return &ping_local{}
}

type ping_local struct {
	act.Actor

	pair gen.PID
}

func (p *ping_local) Init(args ...any) error {
	p.Send(p.PID(), "")
	return nil
}

func (p *ping_local) HandleMessage(from gen.PID, message any) error {
	if _, err := p.MonitorEvent(EVENT); err != nil {
		return err
	}

	pid, err := p.Spawn(factory_pong, gen.ProcessOptions{})
	if err != nil {
		return err
	}
	p.pair = pid

	WGready.Done()
	return nil
}

func (p *ping_local) HandleEvent(message gen.MessageEvent) error {
	switch m := message.Message.(type) {
	case startSend:
		WG.Add(1 + m.n)
		WGready.Done()
		for i := 0; i < m.n; i++ {
			p.SendPID(p.pair, 1)
		}
		WG.Done()

	default:
		p.Log().Warning("unknown event: %#v", message)
	}

	return nil
}
