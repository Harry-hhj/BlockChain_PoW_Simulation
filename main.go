package main

import (
	"math/rand"
	"time"
)

const (
	MinersNumber   = 5
	AttackerNumber = 2
	MonitorNumber  = 1
	MaxChannelSize = 1000
	Interval       = 1000 // ms
	IntervalNum    = 5
)

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
	}
	cahoots := make(map[uint64]chan ConspiratorialTarget)
	for i := uint64(0); i < AttackerNumber; i++ {
		cahoots[uint64(i)] = make(chan ConspiratorialTarget, MaxChannelSize)
	}
	for i := uint64(0); i < AttackerNumber; i++ {
		attackers[i] = NewAttacker(i+MinersNumber, i, *genesisBlock, peers, peers[i+MinersNumber], cahoots, cahoots[i])
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
