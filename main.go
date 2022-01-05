package main

import (
	"math"
	"math/rand"
	"time"
)

const (
	MinersNumber   = 4
	AttackerNumber = 1
	MonitorNumber  = 1
	MaxChannelSize = 1000
	Interval       = 1000 // ms
)


// 制造一个创世区块
func CreateGenesisBlock(data string) *Block {
	genesisBlock := new(Block)
	genesisBlock.UnixMilli = time.Now().UnixMilli()
	genesisBlock.Data = data
	genesisBlock.MinerId = math.MaxUint64
	genesisBlock.LastHash = "0000000000000000000000000000000000000000000000000000000000000000"
	genesisBlock.Height = 1
	genesisBlock.Nonce = 0
	genesisBlock.Target = 19
	genesisBlock.getHash()
	return genesisBlock
}

func main() {
	// random seed
	rand.Seed(time.Now().Unix())

	nodes := make([]*Node, MinersNumber)
	attackers := make([]*Attacker, AttackerNumber)
	peers := make(map[uint64]chan Block)
	for i := 0; i < MinersNumber+AttackerNumber+MonitorNumber; i++ {
		peers[uint64(i)] = make(chan Block, MaxChannelSize)
	}
	genesisBlock := CreateGenesisBlock("IS416 HHJ")
	for i := uint64(0); i < MinersNumber; i++ {
		nodes[i] = NewNode(i, *genesisBlock, peers, peers[i])
		peers[i] <- *genesisBlock
	}
	cahoots := make(map[uint64]chan ConspiratorialTarget)
	for i := uint64(0); i < AttackerNumber; i++ {
		cahoots[uint64(i)] = make(chan ConspiratorialTarget, MaxChannelSize)
	}
	for i := uint64(0); i < AttackerNumber; i++ {
		attackers[i] = NewAttacker(i+MinersNumber, i, *genesisBlock, peers, peers[i+MinersNumber], cahoots, cahoots[i])
		peers[i+MinersNumber] <- *genesisBlock
	}
	monitor := NewMonitor(0, *genesisBlock, peers, peers[MinersNumber+AttackerNumber])
	go monitor.Run()

	// start all nodes
	for i := 0; i < MinersNumber; i++ {
		go nodes[i].Run()
	}
	for i := 0; i < AttackerNumber; i++ {
		go attackers[i].Run()
	}

	// block to wait for all nodes' threads
	<-make(chan int)
}
