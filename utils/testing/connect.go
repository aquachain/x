package main

import (
	"fmt"
	"math/big"
	"time"

	"gitlab.com/aquachain/aquachain/consensus/aquahash"
	"gitlab.com/aquachain/aquachain/core/types"
	"gitlab.com/aquachain/aquachain/params"
)

const (
	start     = 80000
	limit     = 30
	diffstart = 1024 * 9
)

var (
	config  = params.MainnetChainConfig
	big1    = big.NewInt(1)
	big240  = big.NewInt(350)
	timenow = time.Now().Unix()
)

func main() {
	config.HF = map[int]*big.Int{}
	var (
		genesis, parent, head *types.Header
	)

	genesis = &types.Header{
		Version:    2,
		Number:     big.NewInt(start),
		Time:       big.NewInt(timenow),
		Difficulty: big.NewInt(diffstart),
	}

	parent = new(types.Header)
	head = new(types.Header)

	*parent = *genesis
	*head = *genesis

	for i := 0; i < limit; i++ {
		head.Time = head.Time.Add(parent.Time, big240)
		head.Difficulty = aquahash.CalcDifficulty(config, head.Time.Uint64(), parent)
		head.Number = new(big.Int).Add(big1, parent.Number)
		head.ParentHash = parent.Hash()
		printblock(head, parent)
		parent = head
	}

}

func printblock(head, parent *types.Header) {
	if head == nil {
		fmt.Println(head)
		return
	}

	fmt.Printf("header: %x\n", head.Hash())
	fmt.Printf("number: %s\n", head.Number)
	fmt.Printf("time: %s\n", head.Time)
	fmt.Printf("difficulty: %s (%2.2f)\n", head.Difficulty, float64(head.Difficulty.Uint64())/float64(parent.Difficulty.Uint64()))
	fmt.Printf("parent: %x\n\n", head.ParentHash)
}
