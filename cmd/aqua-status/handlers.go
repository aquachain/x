package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"gitlab.com/aquachain/aquachain/params"
)

// type Handler struct {
//      Rpc    *rpc.Client
//      Client *aquaclient.Client
//		cache map[string]*Cache
// }
// func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request){

func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	t1 := time.Now()
	defer func() {
		log.Printf("%s finished in %s", r.URL.Path, time.Since(t1))

	}()
	if h.cached(w, r) {
		return
	}
	cache := h.cache[r.URL.Path]
	cache.b.WriteString("")
	var headblock string
	if err := h.Rpc.Call(&headblock, "aqua_blockNumber"); err != nil {
		log.Println(err)
		return
	}
	n, err := strconv.ParseInt(headblock[2:], 16, 64)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Fprintf(&cache.b, "Block number: %v", n)
	w.Write(cache.b.Bytes())
}

type PoolStats struct {
	Hashrate    int64 `json:'hashrate'`
	MinersTotal int64 `json:'minersTotal'`
	Nodes       struct {
		Height uint64 `json:'height'`
	} `json:'nodes'`
}

func (h *Handler) StatusMiners(w http.ResponseWriter, r *http.Request) {
	t1 := time.Now()
	defer func() {
		log.Printf("%s finished in %s", r.URL.Path, time.Since(t1))
	}()
	if h.cached(w, r) {
		return
	}
	cache := h.cache[r.URL.Path]
	pools, err := h.httpclient.GetBytes("https://aquachain.github.io/pools.json")
	if err != nil || pools == nil {
		log.Println("couldn't fetch pools, using default. err:", err)
		pools = []byte{0x5b, 0xa, 0x20, 0x20, 0x22, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x61, 0x71, 0x75, 0x61, 0x63, 0x68, 0x61, 0x2e, 0x69, 0x6e, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x22, 0x2c, 0xa, 0x20, 0x20, 0x22, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x61, 0x71, 0x75, 0x61, 0x2e, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x32, 0x6e, 0x6f, 0x69, 0x2e, 0x73, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x22, 0x2c, 0xa, 0x20, 0x20, 0x22, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x70, 0x6f, 0x6f, 0x6c, 0x2e, 0x61, 0x71, 0x75, 0x61, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2d, 0x66, 0x6f, 0x75, 0x6e, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x6f, 0x72, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x22, 0x2c, 0xa, 0x20, 0x20, 0x22, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x61, 0x71, 0x75, 0x61, 0x70, 0x6f, 0x6f, 0x6c, 0x2e, 0x72, 0x70, 0x6c, 0x61, 0x6e, 0x74, 0x2e, 0x78, 0x79, 0x7a, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x22, 0x2c, 0xa, 0x20, 0x20, 0x22, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x72, 0x75, 0x2e, 0x61, 0x71, 0x75, 0x61, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2e, 0x70, 0x6f, 0x6f, 0x6c, 0x2e, 0x63, 0x61, 0x6c, 0x63, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x74, 0x2e, 0x72, 0x75, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x22, 0xa, 0x5d, 0xa}
	}
	list := []string{}
	err = json.Unmarshal(pools, &list)
	if err != nil {
		log.Println("unm pool", err)
		return
	}
	var wg sync.WaitGroup
	for _, endpoint := range list {
		wg.Add(1)
		go func(w, buf io.Writer, poolapi string) {
			defer wg.Done()
			resp, err := h.httpclient.GetBytes(poolapi + "stats")
			if resp == nil || err != nil {
				log.Println(poolapi, resp, err)
				return
			}
			var poolStatus PoolStats
			if err := json.Unmarshal(resp, &poolStatus); err != nil {
				log.Println(poolapi, err)
				return
			}
			line := fmt.Sprintf("%3v miners (%8v h/s) (height %v) %s\n", poolStatus.MinersTotal, poolStatus.Hashrate, poolStatus.Nodes.Height, strings.TrimSuffix(poolapi, "/api/"))
			buf.Write([]byte(line))
			w.Write([]byte(line))
			return
		}(w, &cache.b, endpoint)
	}
	wg.Wait()
	// special case for async
	// w.Write(cache.b.Bytes())
}
func (h *Handler) StatusVersions(w http.ResponseWriter, r *http.Request) {
	t1 := time.Now()
	defer func() {
		log.Printf("%s finished in %s", r.URL.Path, time.Since(t1))
	}()
	if h.cached(w, r) {
		return
	}
	type peer struct {
		Name string `json:'name'`
	}
	type versionCount struct {
		version string
		score   int
	}

	var peers []peer
	if err := h.Rpc.Call(&peers, "admin_peers"); err != nil {
		log.Println(err)
		return
	}
	var (
		scoreboard = map[string]int{}
		versions   = []versionCount{}
	)

	for _, p := range peers {
		scoreboard[p.Name]++
	}
	for version, score := range scoreboard {
		versions = append(versions, versionCount{version, score})
	}

	sort.Slice(versions, func(i, j int) bool {
		return versions[i].score > versions[j].score
	})
	cache := h.cache[r.URL.Path]
	fmt.Fprintf(&cache.b, "%s \t%s\n", "PEERS", "VERSION")
	fmt.Fprintf(&cache.b, "%s \t%s\n", "_____", "_______")
	for _, v := range versions {
		fmt.Fprintf(&cache.b, "%00v %s\n", v.score, v.version)
	}
	w.Write(cache.b.Bytes())
}
func (h *Handler) Richlist(w http.ResponseWriter, r *http.Request) {
	t1 := time.Now()
	defer func() {
		log.Printf("%s finished in %s", r.URL.Path, time.Since(t1))
	}()
	if h.cached(w, r) {
		return
	}
	richlist := []string{}
	if err := h.Rpc.Call(&richlist, "admin_getRichlist", 500); err != nil {
		log.Println(err)
		return
	}
	for u := range richlist {
		fmt.Fprintf(&h.cache[r.URL.Path].b, "0x%s\n", richlist[u])
	}
	w.Write(h.cache[r.URL.Path].b.Bytes())
	return
}

func (h *Handler) Supply(w http.ResponseWriter, r *http.Request) {
	t1 := time.Now()
	defer func() {
		log.Printf("%s finished in %s", r.URL.Path, time.Since(t1))
	}()
	if h.cached(w, r) {
		return
	}
	cache := h.cache[r.URL.Path]
	var supply string
	if err := h.Rpc.Call(&supply, "admin_supply"); err != nil {
		log.Println(err)
		return
	}
	var BigAqua = new(big.Float).SetFloat64(params.Aqua)
	n, _ := new(big.Float).SetString(supply)
	fmt.Fprintf(&cache.b, "%v", new(big.Float).Quo(n, BigAqua))
	w.Write(cache.b.Bytes())
}
func (h *Handler) SupplyWei(w http.ResponseWriter, r *http.Request) {
	t1 := time.Now()
	defer func() {
		log.Printf("%s finished in %s", r.URL.Path, time.Since(t1))
	}()
	if h.cached(w, r) {
		return
	}
	cache := h.cache[r.URL.Path]
	var supply string
	if err := h.Rpc.Call(&supply, "admin_supply"); err != nil {
		log.Println(err)
		return
	}
	n, _ := new(big.Int).SetString(supply[2:], 16)
	fmt.Fprintf(&cache.b, "%v", n)
	w.Write(cache.b.Bytes())
}
