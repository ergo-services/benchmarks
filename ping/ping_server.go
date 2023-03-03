package main

import (
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
)

type ping struct {
	gen.Server
}

func (p *ping) Init(process *gen.ServerProcess, args ...etf.Term) error {
	return nil
}
