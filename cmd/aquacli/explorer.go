package main

import (
	"context"
	"math/big"
	"strings"

	"gitlab.com/aquachain/aquachain/common/log"
	"gitlab.com/aquachain/aquachain/params"
	"gitlab.com/aquachain/aquachain/rpc"
	cli "gopkg.in/urfave/cli.v1"
)

var BigAqua = new(big.Float).SetFloat64(params.Aqua)

type System struct{}

func NewSystem() *System {
	return &System{}
}
func (s *System) Run(ctx *cli.Context) error {
	log.Info("Welcome. Connecting to Aquachain.", "node", ctx.String("rpc"))
	rpc, err := getclient(ctx)
	if err != nil {
		return err
	}
	var headblock string
	err = rpc.Call(&headblock, "aqua_blockNumber")
	if err != nil {
		return err
	}
	nmb, _ := new(big.Int).SetString(headblock, 10)
	log.Info("Current Block:", "number", nmb)

	return nil
}

func getclient(ctx *cli.Context) (*rpc.Client, error) {
	if strings.HasPrefix(ctx.String("rpc"), "http") {
		return rpc.DialHTTP(ctx.String("rpc"))
	} else {
		return rpc.DialIPC(context.Background(), ctx.String("rpc"))
	}
}
