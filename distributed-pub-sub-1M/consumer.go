package main

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

func factory_consumer() gen.ProcessBehavior {
	return &consumer{}
}

type consumer struct {
	act.Actor

	event gen.Event
}

type doSubscribe struct{}

func (c *consumer) Init(args ...any) error {
	c.event = args[0].(gen.Event)
	c.Send(c.PID(), doSubscribe{})
	return nil
}

func (c *consumer) HandleMessage(from gen.PID, message any) error {
	switch message.(type) {
	case doSubscribe:
		if _, err := c.MonitorEvent(c.event); err != nil {
			if err == gen.ErrTimeout {
				// Retry on timeout
				c.Send(c.PID(), doSubscribe{})
				return nil
			}
			return err
		}
		WGready.Done()
	}
	return nil
}

func (c *consumer) HandleEvent(message gen.MessageEvent) error {
	switch message.Message.(type) {
	case eventMessage:
		WGreceive.Done()
	}
	return nil
}
