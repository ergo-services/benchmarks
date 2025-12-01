package main

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

func factory_producer() gen.ProcessBehavior {
	return &producer{}
}

type producer struct {
	act.Actor

	token     gen.Ref
	eventName gen.Atom
}

type doRegister struct{}

func (p *producer) Init(args ...any) error {
	p.eventName = args[0].(gen.Atom)
	p.Send(p.PID(), doRegister{})
	return nil
}

func (p *producer) HandleMessage(from gen.PID, message any) error {
	switch message.(type) {
	case doRegister:
		token, err := p.RegisterEvent(p.eventName, gen.EventOptions{})
		if err != nil {
			return err
		}
		p.token = token
		p.Log().Info("Producer registered event '%s'", p.eventName)
		WGready.Done()

	case startPublish:
		p.Log().Info("Producer publishing event...")
		if err := p.SendEvent(p.eventName, p.token, eventMessage{Payload: "test"}); err != nil {
			p.Log().Error("Failed to publish event: %v", err)
			return err
		}
		WGpublish.Done()
	}
	return nil
}
