package main

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/aquanetwork/aquachain/common/log"
	"github.com/aquanetwork/aquachain/core/state"
	"github.com/aquanetwork/aquachain/core/types"
	"github.com/aquanetwork/aquachain/opt/aquaclient"
	"github.com/aquanetwork/aquachain/params"
	"github.com/aquanetwork/aquachain/rpc"

	chart "github.com/wcharczuk/go-chart"
	"gopkg.in/urfave/cli.v1"
)

func dummy(ctx *cli.Context) error { return nil }

func getclient(ctx *cli.Context) (*rpc.Client, error) {
	if strings.HasPrefix(ctx.String("rpc"), "http") {
		return rpc.DialHTTP(ctx.String("rpc"))
	} else {
		return rpc.DialIPC(context.Background(), ctx.String("rpc"))
	}
}

var BigAqua = new(big.Float).SetFloat64(params.Aqua)

func richlist(ctx *cli.Context) error {
	client, err := getclient(ctx)
	if err != nil {
		return err
	}
	var (
		start   = time.Now()
		result  = map[string]state.DumpAccount{} // address -> balance
		results = []struct {
			s  string
			ss state.DumpAccount
		}{}
	)

	err = client.Call(&result, "admin_getDistribution")
	if err != nil {
		return err
	}

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
			f = f.Quo(f, BigAqua)
			fmt.Printf("%v: %00.0008f\n", v.s, f)
			indexes = append(indexes, float64(i+1))
			f6, _ := f.Float64()
			balances = append(balances, f6)
		}
	}

	graph := newChart(indexes, balances, "Distribution", "balance")
	fmt.Printf("chart render done in %v\n", time.Since(start))
	var w io.Writer
	w = os.Stdout
	if ctx.String("o") != "" {
		log.Debug("writing", "file", ctx.String("o"))
		w, err = os.OpenFile(ctx.String("o"), os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
	}
	return drawChart(graph, w)
}

func drawchart_timing(ctx *cli.Context) error {
	return fmt.Errorf("not implemented")
}

//
// func analyzeDifficulty(ctx *cli.Context) ([]float64, []float64, error) {
// 	rc, err := getclient(ctx)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	client := aquaclient.NewClient(rc)
// 	f, err := os.Create(ctx.String("o"))
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	defer f.Close()
// 	begin := ctx.Uint64("from")
// 	if begin == 0 {
// 		return nil, nil, fmt.Errorf("cant start from zero yet")
// 	}
// 	head, err := client.HeaderByNumber(context.Background(), nil)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	var (
// 		headNumber = head.Number.Uint64()
// 		numbers    = []float64{}
// 		diffs      = []float64{}
// 		times      = []float64{}
// 	)
// 	var latest *types.Header
// 	for i := begin; i < headNumber; i++ {
// 		h, err := client.HeaderByNumber(context.Background(), new(big.Int).SetUint64(i))
// 		if err != nil {
// 			return nil, nil, err
// 		}
// 		timediff := uint64(240)
// 		if latest != nil {
// 			timediff = h.Time.Uint64() - latest.Time.Uint64()
// 		}
// 		log.Debug("processing", "header", h.Number, "miner", h.Coinbase.Hex(), "blocktiming", int64(time.Duration(timediff*uint64(time.Second)).Seconds()), "difficulty", h.Difficulty.Uint64())
// 		numbers = append(numbers, float64(i))
// 		times = append(times, float64(timediff*1000))
// 		diffs = append(diffs, float64(h.Difficulty.Uint64()))
// 		latest = h
// 	}
// 	return numbers, diffs, err
// }
func analyzeDifficulty(ctx *cli.Context) ([]float64, []float64, []float64, error) {
	rc, err := getclient(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	client := aquaclient.NewClient(rc)
	f, err := os.Create(ctx.String("o"))
	if err != nil {
		return nil, nil, nil, err
	}
	defer f.Close()
	begin := ctx.Uint64("from")
	head, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, nil, nil, err
	}
	var (
		headNumber = head.Number.Uint64()
		numbers    = []float64{}
		diffs      = []float64{}
		times      = []float64{}
	)
	var latest *types.Header
	for i := begin; i < headNumber; i++ {
		h, err := client.HeaderByNumber(context.Background(), new(big.Int).SetUint64(i))
		if err != nil {
			return nil, nil, nil, err
		}
		timediff := uint64(240)
		if latest != nil {
			timediff = h.Time.Uint64() - latest.Time.Uint64()
		}
		log.Debug("processing", "header", h.Number, "miner", h.Coinbase.Hex(), "blocktiming", int64(time.Duration(timediff*uint64(time.Second)).Seconds()), "difficulty", h.Difficulty.Uint64())
		numbers = append(numbers, float64(i))
		times = append(times, float64(timediff))
		diffs = append(diffs, float64(h.Difficulty.Uint64()))
		latest = h
	}
	return numbers, diffs, times, err
}

func analyzeRichlist(ctx *cli.Context) *Chart {
	client, err := getclient(ctx)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var (
		//start   = time.Now()
		result  = map[string]state.DumpAccount{} // address -> balance
		results = []struct {
			s  string
			ss state.DumpAccount
		}{}
	)

	err = client.Call(&result, "admin_getDistribution")
	if err != nil {
		fmt.Println(err)
		return nil
	}

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
			f = f.Quo(f, BigAqua)
			fmt.Printf("%v: %00.0008f\n", v.s, f)
			indexes = append(indexes, float64(i+1))
			f6, _ := f.Float64()
			balances = append(balances, f6)
		}
	}

	return Ch(newChart(indexes, balances, "Distribution", "balance"))
}

type Chart struct {
	chart *chart.Chart
}

func (c Chart) HandleChart(w http.ResponseWriter, r *http.Request) {
	drawChart(c.chart, w)
}

func Ch(chart *chart.Chart) *Chart {
	return &Chart{chart: chart}
}

func getcharts_difftiming(ctx *cli.Context) (*Chart, *Chart) {
	log.Info("Reading blockchain")
	t1 := time.Now()
	numbers, diffs, times, err := analyzeDifficulty(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	log.Info("Generating diff.png", "elapsed", time.Now().Sub(t1))
	t2 := time.Now()
	diffchart := newChart(numbers, diffs, "Difficulty", "mining difficulty")

	log.Info("Generating timing.png", "elapsed", time.Now().Sub(t2))
	timechart := newChart(numbers, times, "Timing", "blocktimes (seconds)")
	log.Info("Two charts done", "total", time.Now().Sub(t1))
	return Ch(diffchart), Ch(timechart)
}
