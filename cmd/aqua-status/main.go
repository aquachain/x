package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aerth/tgun"
	"gitlab.com/aquachain/aquachain/opt/aquaclient"
	"gitlab.com/aquachain/aquachain/rpc/rpcclient"
)

const TimeoutDuration = 10 * time.Second

func NewHandler(aquaIPC string) *Handler {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	rpcc, err := rpc.DialIPC(ctx, aquaIPC)
	if err != nil {
		log.Fatal(err)
	}
	client := aquaclient.NewClient(rpcc)
	t := &tgun.Client{Timeout: TimeoutDuration}

	return &Handler{rpcc, client, map[string]*Cache{}, t}
}

type Handler struct {
	Rpc        *rpc.Client
	Client     *aquaclient.Client
	cache      map[string]*Cache
	httpclient *tgun.Client
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
	if r.Method != http.MethodGet {
		return
	}
	switch strings.TrimSuffix(r.URL.Path, "/") {
	case "", "/", "/status":
		h.HomePage(w, r)
	case "/status/miners":
		h.StatusMiners(w, r)
	case "/status/versions", "/status/peers":
		h.StatusVersions(w, r)
	case "/richlist.txt":
		h.Richlist(w, r)
	case "/supply":
		h.Supply(w, r)
	case "/supply.wei":
		h.SupplyWei(w, r)
	default:
		http.NotFound(w, r)
	}
}

func main() {
	var (
		addrFlag = flag.String("addr", "127.0.0.1:8881", "address to listen")
		rpcFlag  = flag.String("rpc", os.ExpandEnv("$HOME/.aquachain/aquachain.ipc"), "path to aqua")
	)
	flag.Parse()
	handler := NewHandler(*rpcFlag)
	err := http.ListenAndServe(*addrFlag, handler)
	if err != nil {
		log.Fatal(err)
	}
}
