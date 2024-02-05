package main

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

func factory_ping_network() gen.ProcessBehavior {
	return &ping_network{}
}

type ping_network struct {
	act.Actor

	remote      gen.Atom
	remote_pong gen.Atom
	pair        gen.PID
}

func (p *ping_network) Init(args ...any) error {
	p.remote = args[0].(gen.Atom)
	p.remote_pong = args[1].(gen.Atom)
	p.Send(p.PID(), "")
	return nil
}

func (p *ping_network) HandleMessage(from gen.PID, message any) error {
	if _, err := p.MonitorEvent(EVENT); err != nil {
		return err
	}
	remote, err := p.Node().Network().Node(p.remote)
	if err != nil {
		return err
	}
	pid, err := remote.Spawn(p.remote_pong, gen.ProcessOptions{})
	if err != nil {
		return err
	}
	p.pair = pid
	WGready.Done()
	return nil
}

func (p *ping_network) HandleEvent(message gen.MessageEvent) error {
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
