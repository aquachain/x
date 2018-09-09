package main

import (
	"context"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"sync"

	"gitlab.com/aquachain/aquachain/common/log"
	"gitlab.com/aquachain/aquachain/core/state"
	"gitlab.com/aquachain/aquachain/params"
	"gitlab.com/aquachain/aquachain/rpc"
	cli "gopkg.in/urfave/cli.v1"
)

const (
	DISTRIB int = 0
	POOL
	LOL
)

var BigAqua = new(big.Float).SetFloat64(params.Aqua)

type Explorer struct {
	rpc *rpc.Client
	mu  sync.Mutex          // guards entire cache map
	c   map[int]interface{} // cache
	ctx *cli.Context        // cli args and flags

}

// NewExplorer returns new instance of Explorer. It is ready to use.
func NewExplorer(ctx *cli.Context) *Explorer {
	return &Explorer{
		c:   map[int]interface{}{},
		ctx: ctx,
	}
}

func (e *Explorer) Serve() (err error) {
	if err := e.refreshDistribution(); err != nil {
		return err
	}
	return nil
}
func (e *Explorer) refreshDistribution() (err error) {
	type distribResult struct {
		s  string
		ss state.DumpAccount
	}
	var (
		result    = map[string]state.DumpAccount{}
		results   = []distribResult{}
		headblock = ""
	)

	log.Info("Welcome. Connecting to Aquachain.", "node", e.ctx.String("rpc"))
	e.rpc, err = getclient(e.ctx)
	if err != nil {
		return err
	}

	err = e.rpc.Call(&headblock, "aqua_blockNumber")
	if err != nil {
		return err
	}

	nmb, _ := new(big.Int).SetString(headblock, 10)
	log.Info("Current Block:", "Number", nmb)

	err = e.rpc.Call(&result, "admin_getDistribution")
	if err != nil {
		return err
	}

	for address, bal := range result {
		results = append(results, distribResult{s: address, ss: bal})
	}

	sort.Slice(results, func(i, j int) bool {
		ii, _ := new(big.Int).SetString(results[i].ss.Balance, 10)
		jj, _ := new(big.Int).SetString(results[j].ss.Balance, 10)
		return ii.Cmp(jj) > 0
	})

	for i, v := range results {
		if v.ss.Balance != "0" {
			f, _ := new(big.Float).SetString(v.ss.Balance)
			f = f.Quo(f, BigAqua)
			log.Info("holder "+strconv.Itoa(i), "acct", v.s, "bal", f)
		}
	}
	e.mu.Lock()
	e.c[DISTRIB] = results
	e.mu.Unlock()
	return nil
}

func getclient(ctx *cli.Context) (*rpc.Client, error) {
	if strings.HasPrefix(ctx.String("rpc"), "http") {
		return rpc.DialHTTP(ctx.String("rpc"))
	} else {
		return rpc.DialIPC(context.Background(), ctx.String("rpc"))
	}
}
