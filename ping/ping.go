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
	return nil
}

func (p *ping) HandleMessage(from gen.PID, message any) error {
	if _, err := p.MonitorEvent(sendEvent); err != nil {
		return err
	}
	p.Log().Info("subscribed to %s", sendEvent)
	wg.Add(1)
	return nil
}

func (p *ping) HandleEvent(message gen.MessageEvent) error {
	switch m := message.Message.(type) {
	case sendCase11:
		p.Log().Info("sending %d messages", m.n)
		for i := 0; i < m.n; i++ {
			wg.Add(1)
			p.Send(m.to, "hi")
		}
		p.Log().Info("sent %d messages", m.n)
		wg.Done()
	case sendCase1N:
		p.Log().Info("sending %d messages", m.n)
		l := len(m.to)
		for i := 0; i < m.n; i++ {
			wg.Add(1)
			n := i % l
			p.Send(m.to[n], "hi")
		}
		p.Log().Info("sent %d messages", m.n)
		wg.Done()

	default:
		p.Log().Warning("unknown event: %#v", message)
	}

	return nil
}
