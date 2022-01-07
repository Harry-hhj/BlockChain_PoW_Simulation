package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"time"
)

// Define your message's struct here
type Block struct {
	LastHash  string // 上一个区块的 Hash
	Hash      string // 本区块 Hash
	MinerId   uint64 // 矿工 ID
	Data      string // 区块存储的数据（比如比特币UTXO模型 则此处可用于存储交易）
	UnixMilli int64  // 时间戳
	Height    uint   // 区块高度
	Target    uint   // 难度值
	Nonce     int64  // 随机数
}

func (b *Block) serialize() []byte {
	bytes, err := json.Marshal(b)
	if err != nil {
		log.Panic(err)
	}
	return bytes
}

func (block *Block) Deserialize(str string) {
	err := json.Unmarshal([]byte(str), block)
	if err != nil {
		fmt.Printf("json.Marshal,err:%s", err)
	}
}

func (b *Block) getHash() {
	result := sha256.Sum256(b.serialize())
	b.Hash = hex.EncodeToString(result[:])
}

// 判断工作量证明是否有效 TODO:Difficulty
func (b *Block) IsValid() bool {
	target := big.NewInt(1)
	target.Lsh(target, 256-b.Target)
	hashBytes, _ := hex.DecodeString(b.Hash)
	hashInt := big.Int{}
	hashInt.SetBytes(hashBytes)
	if hashInt.Cmp(target) == -1 {
		return true
	} else {
		return false
	}
}

// Block Chain
type BlockChainNode struct {
	block  *Block
	childs []*BlockChainNode
	parent *BlockChainNode
}
type BlockChain struct {
	root      *BlockChainNode
	workspace *BlockChainNode
	maxHeight uint
	index     map[uint][]*BlockChainNode
}

func NewBlockChain(genesisBlock *Block) *BlockChain {
	root := new(BlockChainNode)
	root.block = genesisBlock
	root.childs = nil
	root.parent = nil
	newBlockChain := new(BlockChain)
	newBlockChain.root = root
	newBlockChain.workspace = root
	newBlockChain.maxHeight = 1
	newBlockChain.index = make(map[uint][]*BlockChainNode)
	newBlockChain.index[1] = append(newBlockChain.index[1], newBlockChain.root)
	return newBlockChain
}

func (bc *BlockChain) append(b *Block) bool {
	if !b.IsValid() {
		fmt.Print("APPEND NOT VALID!\n")
		return false
	}
	if bc.index[b.Height-1] == nil {
		fmt.Printf("b.Height = %d, NIL!\n", b.Height)
		return false
	} else {
		for i := 0; i < len(bc.index[b.Height]); i++ {
			if b.Hash == bc.index[b.Height][i].block.Hash {
				fmt.Print("SAME\n")
				return false
			}
		}
		for i := 0; i < len(bc.index[b.Height-1]); i++ {
			if b.LastHash == bc.index[b.Height-1][i].block.Hash {
				newBCN := new(BlockChainNode)
				newBCN.block = b
				newBCN.childs = nil
				newBCN.parent = bc.index[b.Height-1][i]
				bc.index[b.Height-1][i].childs = append(bc.index[b.Height-1][i].childs, newBCN)
				bc.index[b.Height] = append(bc.index[b.Height], newBCN)
				if b.Height > bc.maxHeight {
					// switch to the longest chain
					bc.maxHeight = b.Height
					bc.workspace = newBCN
					return true
				}
				return false
			}
		}
		return false
	}
}

func (bc *BlockChain) statistics(n *BlockChainNode) int64 {
	if n.block.Height < IntervalNum+1 {
		return IntervalNum*Interval
	}
	cur := n.block.UnixMilli
	lstPtr := n
	for i:=0; i<IntervalNum; i++ {
		lstPtr = lstPtr.parent
	}
	lst := lstPtr.block.UnixMilli
	return cur - lst
}

func (bc *BlockChain) search(b *Block) *BlockChainNode {
	if b == nil {
		return nil
	}
	for _, bcn := range bc.index[b.Height] {
		if bcn.block.Hash == b.Hash {
			return bcn
		} 
	}
	return nil
}

func (bc *BlockChain) print() {
	tmp := []*BlockChainNode{}
	tmp = append(tmp, bc.workspace)
	for tmp[len(tmp)-1].parent != nil {
		tmp = append(tmp, tmp[len(tmp)-1].parent)
	}
	fmt.Print("===============================\n")
	for i := len(tmp) - 1; i >= 0; i-- {
		fmt.Println(tmp[i].block)
	}
	fmt.Print("===============================\n")
}

func (bc *BlockChain) formatPrintAll() {
	fmt.Print("===============================\n")
	stack := NewStack()
	stack.Push(bc.root)
	for !stack.IsEmpty() {
		cbcn := stack.Pop()
		for i := uint(1); i < cbcn.block.Height; i++ {
			fmt.Print("\t")
		}
		fmt.Println("└──", cbcn.block.Data)
		for _, bcn := range cbcn.childs {
			stack.Push(bcn)
		}
	}
	fmt.Print("===============================\n")
}

// 制造一个创世区块
func CreateGenesisBlock(data string) *Block {
	genesisBlock := new(Block)
	genesisBlock.UnixMilli = time.Now().UnixMilli()
	genesisBlock.Data = data
	genesisBlock.MinerId = math.MaxUint64
	genesisBlock.LastHash = "0000000000000000000000000000000000000000000000000000000000000000"
	genesisBlock.Height = 1
	genesisBlock.Nonce = rand.Int63()
	genesisBlock.Target = 19
	newBigInt := big.NewInt(1)
	newBigInt.Lsh(newBigInt, 256-genesisBlock.Target)
	for {
		genesisBlock.getHash()
		hashInt := big.Int{}
		hashBytes, _ := hex.DecodeString(genesisBlock.Hash)
		hashInt.SetBytes(hashBytes)
		if hashInt.Cmp(newBigInt) == -1 {
			break
		} else {
			genesisBlock.Nonce++
		}
	}
	return genesisBlock
}