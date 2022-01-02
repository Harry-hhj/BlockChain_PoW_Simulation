package main

import (
	"fmt"
)

const (
	NodeNumber     = 4
	MaxChannelSize = 1000
)

// Define your message's struct here
type Message struct {
	sender uint64
}

// May need other fields
type Node struct {
	id          uint64
	peers       map[uint64]chan Message
	receiveChan chan Message
}

func NewNode(id uint64, peers map[uint64]chan Message, recvChan chan Message) *Node {
	return &Node{
		id:          id,
		peers:       peers,
		receiveChan: recvChan,
	}
}

// Run start one consensus node, remember to start one thread for receiving messages (n.Receive())
func (n *Node) Run() {
	fmt.Println("start node : ", n.id)
	go n.Receive()

	n.Broadcast(Message{sender : n.id})
}

func (n *Node) Receive() {
	for {
		select {
		case msg := <-n.receiveChan:
			n.handler(msg)
		}
	}
}

func (n *Node) handler(msg Message) {
	fmt.Println("Node", n.id, "received message from node", msg.sender)
}

func (n *Node) Broadcast(msg Message) {
	for id, ch := range n.peers {
		if id == n.id {
			continue
		}
		ch <- msg
	}
}

func main() {
	nodes := make([]*Node, NodeNumber)
	peers := make(map[uint64]chan Message)
	for i := 0; i < NodeNumber; i++ {
		peers[uint64(i)] = make(chan Message, MaxChannelSize)
	}
	for i := uint64(0); i < NodeNumber; i++ {
		nodes[i] = NewNode(i, peers, peers[i])
	}

	// start all nodes
	for i := 0; i < NodeNumber; i++ {
		go nodes[i].Run()
	}

	// block to wait for all nodes' threads
	<-make(chan int)
}
