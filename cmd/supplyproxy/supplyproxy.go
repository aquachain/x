package main

import (
	"encoding/json"
	"log"
  "strings"
  "context"

	"net/http"
	"time"
  "gopkg.in/urfave/cli.v1"
	"github.com/gorilla/mux"
	"gitlab.com/aquachain/aquachain/rpc"
	"gitlab.com/aquachain/aquachain/common/hexutil"
)

type Coin struct {
	Algo            string  `json:"algo"`
	Port            string  `json:"port"`
	Name            string  `json:"name"`
	Height          int     `json:"height"`
	Workers         int     `json:"workers"`
	Shares          int     `json:"shares"`
	Hashrate        int64   `json:"hashrate"`
	Estimate        string  `json:"estimate"`
	Two4HBlocks     int     `json:"24h_blocks"`
	PercentBlocks   string  `json:"percent_blocks"`
	Two4HBtc        float64 `json:"24h_btc"`
	Two4HCoins      string  `json:"24h_coins"`
	Lastblock       int     `json:"lastblock"`
	Timesincelast   int     `json:"timesincelast"`
	Difficulty      string  `json:"difficulty"`
	NetworkHashrate float64 `json:"network_hashrate"`
	BlockReward     string  `json:"block_reward"`
}

type System struct {
	rpc *rpc.Client
}

func (s *System) HandleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	json.NewEncoder(w).Encode(time.Now())
}

func (s *System) HandleSupply(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")

  var headblock = ""
	err := s.rpc.Call(&headblock, "aqua_blockNumber")
	if err != nil {
		log.Println("error:", err)
    return
	}

	height, err := hexutil.DecodeUint64(headblock)
	if err != nil {
		log.Println("error:", err)
		return
	}

	log.Println("current height", height)
	json.NewEncoder(w).Encode(height)
}

func (s *System) listen(ctx *cli.Context) error {
	r := mux.NewRouter()
	r.HandleFunc("/api/status", s.HandleStatus)
	r.HandleFunc("/api/supply", s.HandleSupply)
	r.NotFoundHandler = http.HandlerFunc(notFound)
	return http.ListenAndServe(ctx.String("addr"), r)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusNotFound)
}

func getclient(ctx *cli.Context) (*rpc.Client, error) {
	if strings.HasPrefix(ctx.String("rpc"), "http") {
		return rpc.DialHTTP(ctx.String("rpc"))
	} else {
		return rpc.DialIPC(context.Background(), ctx.String("rpc"))
	}
}

func run(ctx *cli.Context) error {
  client, err := getclient(ctx)
  if err != nil {
    log.Fatalf("Epic fail: %v", err)
  }
	s := &System{rpc: client}
  return s.listen(ctx)
}
