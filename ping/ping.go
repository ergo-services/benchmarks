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
	p.SetKeepOrder(true)
	return nil
}

func (p *ping) HandleMessage(from gen.PID, message any) error {
	if _, err := p.MonitorEvent(sendEvent); err != nil {
		return err
	}
	// p.Log().Info("subscribed to %s", sendEvent)
	wg.Add(1)
	return nil
}

func (p *ping) HandleEvent(message gen.MessageEvent) error {
	switch m := message.Message.(type) {
	case sendCase11:
		// p.Log().Info("sending %d messages", m.n)
		for i := 0; i < m.n; i++ {
			wg.Add(1)
			p.Send(m.to, "hi")
		}
		// p.Log().Info("sent %d messages", m.n)
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
		// p.Log().Info("sent %d messages", m.n)
		wg.Done()

	case sendCaseNN:
		// p.Log().Info("sending %d messages", m.n)
		p.SetKeepOrder(false)
		l := len(m.to)
		for i := 0; i < m.n; i++ {
			wg.Add(1)
			n := i % l
			p.Send(m.to[n], "hi")
		}
		// p.Log().Info("sent %d messages", m.n)
		wg.Done()

	default:
		p.Log().Warning("unknown event: %#v", message)
	}

	return nil
}
