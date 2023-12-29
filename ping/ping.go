package main

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

func factory_ping() gen.ProcessBehavior {
	return &ping{}
}

type ping struct {
	act.Actor
}

func (p *ping) Init(args ...any) error {
	p.Send(p.PID(), "")
	p.SetKeepNetworkOrder(true)
	return nil
}

func (p *ping) HandleMessage(from gen.PID, message any) error {
	if _, err := p.MonitorEvent(sendEvent); err != nil {
		return err
	}
	wg.Add(1)
	return nil
}

func (p *ping) HandleEvent(message gen.MessageEvent) error {
	switch m := message.Message.(type) {
	case sendCase11:
		for i := 0; i < m.n; i++ {
			wg.Add(1)
			p.Send(m.to, "hi")
		}
		wg.Done()

	case sendCase1N:
		// If we send messages sequentially over the process list m.to it makes the receiving
		// process switch state back and forth from sleep to running. So use a bit different approach...
		l := len(m.to)
		x := m.n/l + m.n%l
		for i := 0; i < m.n; i++ {
			n := i / x
			wg.Add(1)
			p.Send(m.to[n], "hi")
		}
		wg.Done()

	case sendCaseNN:
		p.SetKeepNetworkOrder(false)
		l := len(m.to)
		n := int(p.PID().ID) % l
		for i := 0; i < m.n; i++ {
			wg.Add(1)
			p.Send(m.to[n], "hi")
		}
		wg.Done()

	default:
		p.Log().Warning("unknown event: %#v", message)
	}

	return nil
}
