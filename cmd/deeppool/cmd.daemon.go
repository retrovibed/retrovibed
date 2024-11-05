package main

import "github.com/james-lawrence/deeppool/cmd/cmdopts"

type cmdDaemon struct{}

func (t cmdDaemon) Run(ctx *cmdopts.Global) (err error) {

	<-ctx.Context.Done()
	return nil
}
