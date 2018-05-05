// Copyright 2015 The aquachain Authors
// This file is part of the aquachain library.
//
// The aquachain library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The aquachain library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the aquachain library. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aquachain/x/utils"
	"github.com/aquanetwork/aquachain/core/state"
	"github.com/aquanetwork/aquachain/core/types"
	"github.com/aquanetwork/aquachain/params"
	"github.com/aquanetwork/aquachain/rpc"
	"github.com/aquanetwork/ok/aquachain/aquaclient"
	chart "github.com/wcharczuk/go-chart"

	"gopkg.in/urfave/cli.v1"
)

var gitCommit = ""
var (
	app = utils.NewApp(gitCommit, "usage")
)

func init() {
	app.Name = "aqua-explorer"
	app.Action = switcher
	app.Flags = []cli.Flag{cli.StringFlag{
		Value: filepath.Join(utils.DataDirFlag.Value.String(), "aquachain.ipc"),
		Name:  "rpc",
		Usage: "path or url to rpc",
	}}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func switcher(ctx *cli.Context) error {
	if argc := ctx.NArg(); argc == 0 {
		return fmt.Errorf("command not specified: timing diff rich")
	}
	arg1 := ctx.Args()[0]
	switch arg1 {
	case "timing":
		return drawchart_timing(ctx)
	case "diff":
		return drawchart_diff(ctx)
	case "rich", "richlist":
		return richlist(ctx)
	default:
		return fmt.Errorf("command not found")
	}

}

func dummy(ctx *cli.Context) error { return nil }

func getclient(ctx *cli.Context) (*rpc.Client, error) {
	if strings.HasPrefix(ctx.String("rpc"), "http") {
		return rpc.DialHTTP(ctx.String("rpc"))
	} else {
		return rpc.DialIPC(context.Background(), ctx.String("rpc"))
	}
}

func richlist(ctx *cli.Context) error {
	client, err := getclient(ctx)
	if err != nil {
		return err
	}
	f, err := os.Create("richlist.png")
	if err != nil {
		return err
	}
	defer f.Close()

	start := time.Now()
	var result = map[string]state.DumpAccount{} // address -> balance
	if err := client.Call(&result, "admin_getDistribution"); err != nil {
		return err
	}
	bigaqua := new(big.Float).SetFloat64(params.Aqua)
	results := []struct {
		s  string
		ss state.DumpAccount
	}{}
	for i, v := range result {
		ss := struct {
			s  string
			ss state.DumpAccount
		}{i, v}
		results = append(results, ss)
	}
	sort.Slice(results, func(i int, j int) bool {
		ii, _ := new(big.Int).SetString(results[i].ss.Balance, 10)
		jj, _ := new(big.Int).SetString(results[j].ss.Balance, 10)
		return ii.Cmp(jj) > 0
	})
	var balances, indexes []float64
	for i, v := range results {
		if v.ss.Balance != "0" {
			f, _ := new(big.Float).SetString(v.ss.Balance)
			f = f.Quo(f, bigaqua)
			fmt.Printf("%v: %00.0008f\n", v.s, f)
			indexes = append(indexes, float64(i+1))
			f6, _ := f.Float64()

			balances = append(balances, f6)
		}
	}

	graph := newChart(indexes, balances, "Distribution", "balance")
	err = graph.Render(chart.PNG, f)
	if err != nil {
		utils.Fatalf("Render error: %v\n", err)
	}
	fmt.Printf("Image render done in %v\n", time.Since(start))
	return nil
}

func drawchart_timing(ctx *cli.Context) error {
	return fmt.Errorf("not implemented")
}

func drawchart_diff(ctx *cli.Context) error {
	rc, err := getclient(ctx)
	if err != nil {
		return err
	}
	client := aquaclient.NewClient(rc)
	f, err := os.Create("latest.png")

	if err != nil {
		return err
	}
	defer f.Close()

	count := uint64(360)

	rawstr := ctx.Args()[0]
	if rawstr != "" {
		num, err := strconv.Atoi(rawstr)
		if err == nil {
			if num > 9 && num < 1000 {
				count = uint64(num)
			}
		}
	}

	head, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}
	var (
		start      = time.Now()
		headNumber = head.Number.Uint64()
		numbers    = []float64{}
		diffs      = []float64{}
		times      = []float64{}
	)
	//fp := ctx.Args().First()

	var latest *types.Header

	for i := uint64(headNumber - count); i < headNumber; i++ {
		h, err := client.HeaderByNumber(context.Background(), new(big.Int).SetUint64(i))
		if err != nil {
			return err
		}
		timediff := uint64(240)
		if latest != nil {
			timediff = h.Time.Uint64() - latest.Time.Uint64()
		}
		log.Println("Processing header:", h.Number)
		log.Println(float64(h.Number.Uint64()), float64(timediff), float64(h.Difficulty.Uint64()))
		numbers = append(numbers, float64(i))
		times = append(times, float64(timediff*1000))
		diffs = append(diffs, float64(h.Difficulty.Uint64()))
		latest = h
	}

	graph := newChart(numbers, times, "Difficulty, most recent "+strconv.Itoa(int(count))+" blocks", "mining difficulty")
	err = graph.Render(chart.PNG, f)
	if err != nil {
		utils.Fatalf("Render error: %v\n", err)
	}
	fmt.Printf("Image render done in %v\n", time.Since(start))
	return nil
}

func float2Tick(f []float64) (ticks []chart.Tick) {
	for i, v := range f {
		if i%4 == 0 {
			ticks = append(ticks, chart.Tick{Label: fmt.Sprintf("%v", int(v)), Value: v})
		}
	}
	return ticks
}
