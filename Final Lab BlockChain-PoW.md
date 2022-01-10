# 基于 PoW 的区块链仿真程序

## 1 Basic Principles

工作量证明 PoW （Proof of Work）是一个用来确认你做过一定量的工作的证明。关注工作的整个过程是较为低效的，而通过对工作结果进行认证来证明完成了相应的工作量是高效的。例如，现实生活中的驾驶证、毕业证、上司更关心你的工作成果（除非他喜欢一直观察你）等，也都是通过检测结果的方式所取得的证明。

工作量证明系统并不是区块链提出的新概念，而是由 Cynthia Dwork & Moni Naor 1993 年在学术论文中首次提出。它要求发起者进行一定量的运算，即消耗一定的计算时间。而工作量证明（PoW）这个名词是 1999 年 Markus Jakobsson & Ari Juels 仔文章中提出的。

哈希现金是一种工作量证明机制，由 Adam Back 在 1997 年发明，用于抵抗邮件的拒绝服务及垃圾邮件网关滥用。在比特币之前，哈希现金被用于用于垃圾邮件的过滤，也被微软用于 hotmail / exchange / outlook 等产品中（微软使用一种与哈希现金不兼容的格式并将之命名为电子邮戳），还被哈尔·芬尼以可重复使用的工作量证明（RPOW）的形式用于一种比特币之前的加密货币实验中。此外，戴伟的 B-money 、尼克·萨博的 Bit-Gold 这些比特币的先行者，都是在哈希现金的框架下进行挖矿的。

上述机制离不开哈希函数。哈希函数（Hash Function），也称为散列函数，给定一个输入 $x$ ，它会算出相应的输出 $H(x)$ ，即构成一种映射关系 $x \mapsto H(x)$ 。哈希函数输入输出的主要特征是：

1.   输入 $x$ 可以是任意长度的字符串
2.   输出结果即 $H(x)$ 的长度是固定的

哈希函数的主要性质是：

1.   确定性：给定相同的输入，产生相同的输出。
2.   快速：计算 $H(x)$ 的过程是高效的（对于长度为 $n$ 的字符串 $x$ ，计算出 $H(x)$ 的时间复杂度应为 $O(n)$ ）
3.   抗碰撞性，分为强/弱抗碰撞性：即对于固定/任意输入 $x$ ，找不到一个 $y$ ，使得 $H(x)=H(y)$ 。注意这种不可行性是指计算上不可行，因为从一个高维空间映射向低维空间，如果是满射，那么一定不是单射。
4.   抗原像性。哈希函数是数学上的单向陷门函数。对于一个给定输出结果 $H(x)$ ，反推出其输入 $x$ 在计算上不可行。注意这里还是计算上不可行，因为只要我们遍历所有输入总是可以找到的，只是时间不可容忍。
5.   雪崩效应：不存在比穷举更好的方法，可以使哈希结果 $H(x)$ 落在特定的范围内。

满足这些性质并不容易，感兴趣可以了解王小云院士 2004 年提出的 MD5 的碰撞攻击。

在讲述构造区块链（挖矿）之前，我们指出成功取决于三个因素：计算资源、网络时延、运气。比特币网络中的任何一个节点，想生成一个新的区块并写入区块链，必须解出比特币网络出的工作量证明题。这道题关键的三个要素是工作量证明函数、区块即难度值。工作量证明函数是这道题的计算方法，区块决定了这道题的输入数据，难度值决定了这道题的所需要的计算量。

对于工作量证明函数，比特币系统中使用的工作量证明函数是 SHA256 ，目前为止，还没有出现对SHA256 的有效攻击，这也是这次试验中使用的哈希函数。对于区块，是存储数据的地方。比特币的区块由区块头及该区块所包含的交易列表组成，其中区块头是工作量证明的输入字符串，而交易列表不是。因此，我们需要让区块头能体现交易列表，使得篡改交易列表是不可能的。为此，使用 Merkle Tree 算法生成 Merkle Root Hash ，并为此作为交易列表的摘要，加入到区块头中。Merkle Tree 的算法示意图如下：

<img src="Final%20Lab%20BlockChain-PoW.assets/555.png" alt="img" style="zoom:80%;" />

但是，我们在本次实验中不考虑如此复杂的结构，仅使用 `data` 作为区块头中的数据字段。对于难度值，它决定了矿工大约需要经过多少次哈希运算才能产生一个合法的区块。难度值必须根据全网算力的变化进行调整，比特币规定新区块产生的速率都保持在 10 分钟一个。注意，当区块产生速度增大，即交易吞吐量增加时，新区块的产生速率仍然是有上限的，即网络的时延。试想一下，当一个区块的产生速率大于这个区块被广播到最远节点的网络延时，那么必然会存在不同的视图（即同一时刻存在两个最长的链），这会导致算力的分散（即两组矿工在不同的最长链上挖矿），在特殊情况下，之后提到的 51% 攻击将会变成 26% 攻击，这使得区块链的安全性降低。难度值直观理解就是规定了生成的哈希开头有多少个 0 ，由于不存在比遍历更好的方法，因此我们可以估算出需要进行的计算次数，这就是为什么哈希可以提现工作量。注意此题是一定有解的，不同矿工可以从不同的随机值开始遍历，这就是之前提到的运气因素，一个好的起点可能和结果很接近，所花的时间也更少。

接下来我们讨论存在的攻击方法： 51% 攻击。假定一个恶意节点试图双花之前的已花费的交易，攻击者需要重做包含这个交易的区块，以及这个区块之后的所有的区块，创建一个比目前诚实区块链更长的区块链。只有网络中的大多数节点都转向攻击者创建的区块链，攻击者的攻击才算成功了。由于每一个区块都包含了之前的所有区块的交易信息，所以随着块高的增加，之前的区块都会被再次确认一次，确认超过6次，可以理解为无法被修改。如果攻击者掌握了全网 51% 的计算资源，那么他就一定能构造出一个比现有最长链更长的区块，因为攻击链的增长速率更快。当节点数不多时，这个下限会变得不稳定，很少的节点实行分叉攻击也可能会成功，因为之前提到过，挖矿和运气有关。需要指出的是，即使不存在攻击者，区块链也可能会在某一时刻分叉。所谓的链分叉，主要是由于在计算 hash 时，每个人拿到的区块内容是不同的，导致算出的区块结果也不同。但它们都是正确结果，于是区块链在这个时刻，出现了两个都满足要求的不同区块。由于距离远近、网络等原因，不同旷工看到这两个区块的先后顺序是不一样的。通常情况下，旷工会把自己先看到的区块链复制过来，然后接着在这个区块上开始新的挖矿工作，于是就出现了链分叉。但由于所有矿工都遵从同样的机制，且这两个良性矿工打包的数据是相同的，因此不会影响区块链的安全性，仅仅是一个矿工白白付出了算力。

基于这些知识，现在我们可以搭建我们的简易区块链系统了。



## 2 Introduction

本项目是为了熟悉区块链以及 PoW 的工作原理，模拟真实的区块链，以及可能出现的分叉攻击的仿真。在一开始，我通过所有矿工在一个全局的区块链进行仿真，基于的原因是在如此小规模的网络和多线程仿真条件下，几乎不会存在节点因不网络延迟存在的视图差异。虽然这么做有一定道理，但为了更好地展示区块链，还是修改成为每一个节点保存一个独自视角的区块，节点之间通过通道（channel）进行信息传递。即使这么做，但我还是离真实的区块链有一定的差异。我指出：

1.   实验的节点是一开始确定的，不能模拟新产生的节点获得区块链视图的过程
2.   我假设网络是不会丢包的，这在真实场景中不一定
3.   受限于实验要求，需要 1s 左右出块，因此受系统调度、运气等因素影响更大

我实现了：

1.   基于多线程仿真的 PoW 区块链系统
2.   PoW 的解题和验证
3.   矿工的挖矿、接收和传播
4.   攻击（分叉攻击）者节点、攻击者秘密信道同步



## 3 Architecture

```go
.
├── blockchain.go  // 区块链实现
├── go.mod  // mod
├── main.go  // main 函数
├── node.go  // 矿工、恶意攻击者、系统监听者实现
└── stack.go  // 栈实现

0 directories, 5 files
```



## 4 Implementation

本实验基于 Go 语言编写，用多线程模拟网络上不同的矿工节点，用 `channel` 模拟网络通信信道和攻击者间的秘密信道。

第一步我们先要实现区块 `Block` ：

```go
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
```

基于此，我们可以构造区块链了。一个需要思考的问题是如何存储区块链，因为我们的区块链也不会很大，不需要考虑物理介质的存储，不用考虑 B+ 树这种存储。一个区块有且只有一个父节点（初创世区块），可以有 0 个或多个子节点（考虑分叉），这种结构最适合树形存储：

```go
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
```

注意到，我们为 `BlockChain` 定义了一个 `index map[uint64][]*BlockChainNode` ，将 `Height` 映射到 `[]*BlockChainNode` ，这是为了实现快速的索引。一方面，它能够快速找到新收到的区块的父节点（还需要判断是否已添加），另一方面，它也更适合攻击者选择最容易分叉的节点（最长链的长度少一的不属于这条链的任意一个节点）。

![IS416.drawio](Final%20Lab%20BlockChain-PoW.assets/IS416.drawio.png)

为了实现区块链的常见操作，我们定义了以下函数：

```go
func (b *Block) serialize() []byte  // 为了进行哈希
func (block *Block) Deserialize(str string)  // 逆序列化
func (b *Block) getHash()  // 计算哈希
func (b *Block) IsValid() bool  // 验证哈希和难度值

func NewBlockChain(genesisBlock *Block) *BlockChain  // 构造初始区块链
func (bc *BlockChain) append(b *Block) bool  // 添加区块
func (bc *BlockChain) statistics(n *BlockChainNode) int64  // 获得 2 个区块构造速度
func (bc *BlockChain) search(b *Block) *BlockChainNode  // 根据 hash 查找区块 
func (bc *BlockChain) print()  // 打印最长链
func (bc *BlockChain) formatPrintAll()  // 打印树形结构（包含分叉）
```

其中为了实现打印树形结构，我们还需要用到栈 `stack` ，具体算法见代码。

第二步我们要实现良性矿工 `Node` ，矿工需要异步接受其他节点的新区块，让自己尽可能地在最长的上挖掘，还需要循环不断的挖矿，这两个任务通过两个线程实现。新的消息是根据前一个区块和节点属性构造的，当收到一个新产生的区块时，矿工的接收线程将自己的更新表示位置成 `true` ，挖矿线程每次重新计算 `nonce` 时会先检查 `update` 标志，如果为 `true` 则放弃本次挖矿，转而到最长的链上进行挖矿。

```go
// May need other fields
type Node struct {
	id              uint64
	blockChain      BlockChain
	blockChainMutex sync.RWMutex
	peers           map[uint64]chan Block
	receiveChan     chan Block
	update          bool
	flagMutex       sync.RWMutex
}

func NewNode(id uint64, genesisBlock Block, peers map[uint64]chan Block, recvChan chan Block) *Node  // 产生一个新的矿工
func (n *Node) Run()  // 启动节点
func (n *Node) Receive() // 启动接收线程
func (n *Node) handler(msg Block)  // 处理接受到的信息
func (n *Node) Broadcast(msg Block) // 广播新的节点
func (n *Node) Mine() *Block  // 挖矿
```

注意到 `blockChain` 和 `update` 会同时被两个线程读写，为了确保安全性，我们需要加入读写锁。矿工挖出一个新的区块或接收到一个新的区块加入自己的区块链时，需要获得 `blockChain` 写锁；判断目前不是在最长链上挖矿后需要获得 `update` 写锁更新状态；挖矿函数获得上一个区块或统计前两个区块耗时需要获得 `blockChain` 读锁，挖矿时判断 状态需要获得 `update` 读锁。读写锁的使用可以让一个节点有多个挖矿函数时获得最好的性能，这可以用于仿真每个节点有不同的算力的情况。

第三步我们实现系统监听员，他的任务就是获取系统中完整的区块链，进行输出和统计信息。他的实现比矿工还要简单，不需要考虑锁，因为区块链仅被接收线程操控。

```go
type Monitor struct {
	id          uint64
	blockChain  BlockChain
	peers       map[uint64]chan Block
	receiveChan chan Block
}

func NewMonitor(id uint64, genesisBlock Block, peers map[uint64]chan Block, recvChan chan Block) *Monitor
func (m *Monitor) Run()
func (m *Monitor) Receive()
func (m *Monitor) handler(msg Block)
```

第四步我们实现攻击者，也是系统中最复杂的节点。我们先分析攻击者需要什么功能。我们这里假设攻击者是不理性的，其目的就是为了破坏系统的安全性，而不是希望实现特定区块的覆盖。为了实行分叉攻击，攻击者需要选择最有可能成功的区块，即当前最长链的倒数第二个区块（原因是攻击者一般是针对一个已经存在的区块而不是一个正在被解决的区块进行攻击）或在区块链最大高度不大于正在解决的区块的高度加 1 时仍将原目标区块作为目标区块。如果有多个通信者，这几个通信者之间应该达成某种共识，选择共同的目标作为攻击对象。这需要为他们建立秘密信道，并设置特殊的通信协议达成共识。

```go
type MsgType uint
// 枚举类型：攻击者之间通信类型
const (
	INIT MsgType = iota
	CONTINUE
	ACKNOWLEGGE  // 攻击者投票机制，未用
	RECLINE      // 攻击者投票机制，未用
)
// 攻击者间定义的秘密消息格式
type ConspiratorialTarget struct {
	target  *Block  // 目标父区块
	msgType MsgType  // 父区块的产生原因
}

type Attacker struct {
	id                uint64  // 攻击者对外的 uid
	secret_id         uint64  // 攻击者内的 uid
	blockChain        BlockChain
	blockChainMutex   sync.RWMutex
	peers             map[uint64]chan Block
	receiveChan       chan Block
	update            bool
	flagMutex         sync.RWMutex
	cahoot            map[uint64]chan ConspiratorialTarget  // 攻击者秘密信道
	secretReceiveChan chan ConspiratorialTarget  // 攻击者秘密接收信道
	request           bool  // 目标父区块更新标志位
	requestMutex      sync.RWMutex  // 标志位锁
	target            *BlockChainNode  // 目标父区块
	targetMutex       sync.RWMutex  // 目标父区块读写锁
	votes             uint64  // 攻击者投票机制，未用
}

func NewAttacker(id uint64, secret_id uint64, genesisBlock Block, peers map[uint64]chan Block, recvChan chan Block, cahoot map[uint64]chan ConspiratorialTarget, secretRecvChan chan ConspiratorialTarget) *Attacker
func (a *Attacker) Run()
func (a *Attacker) Receive()  // 接收两类消息：正常区块信息和攻击者共识消息
func (a *Attacker) Broadcast(msg Block)
func (a *Attacker) handler(msg Block)
func (a *Attacker) shandler(msg ConspiratorialTarget)  // 处理攻击者共识消息
func (a *Attacker) polt(msg ConspiratorialTarget)  // 攻击者内秘密广播
func (a *Attacker) Attack() *Block  // 产生分叉区块
```



## 5 Experiments and Results

在矿工 10 ，恶意节点 4 的情况下，系统监听员记录的区块链增长速度和难度变化值：

![fig1](Final%20Lab%20BlockChain-PoW.assets/fig1.jpg)

无论节点是否恶意，都会遵守工作量证明机制。可以看到，当系统区域稳定后，时间的变化幅度在 10% 以内，考虑到单机仿真的不确定性，这个结果是良好的。

从图 1 中可以看到， 100 个区块后区块的构造速度趋于稳定。为了让系统趋于稳定，从而得到更可靠的攻击可行性数据，我们规定攻击者在区块链生长到 100 的高度时才开始攻击。

我们让良性矿工和攻击者节点各减半，观察到图 2 。

![fig2](Final%20Lab%20BlockChain-PoW.assets/fig2.jpg)

同样系统也会趋于稳定

同样规定攻击者在 100 个区块后开始活动。不同恶意节点比例条件下，以 1000 个区块长度统计平均（实际有效长度 900 ），攻击者尝试分叉攻击中成功的比例：

| 攻击者节点比例 | 尝试分叉攻击的次数 | 成功分叉攻击的次数 | 攻击成功的概率 |
| -------------- | ------------------ | ------------------ | -------------- |
| 10%            | 116                | 13                 | 11.21%         |
| 20%            | 251                | 76                 | 30.28%         |
| 30%            | 304                | 130                | 42.76%         |
| 40%            | 401                | 190                | 47.38%         |

![fig3](Final%20Lab%20BlockChain-PoW.assets/fig3.jpg)

随着恶意节点数的增大，恶意节点尝试攻击的次数增大，成功的次数增大，成功的概率也增大。



## 6 Discussion

本项目用 Go 语言搭建了一个 PoW 的多线程反正程序，每个线程模拟一个节点生成区块链的状态。我构造了一个类似于 Bitcoin 的数据结构，但是省略了交易信息，而用一个字符串信息代替。实现过程中一个比较重要的点是一个节点会有至少两个线程（挖矿和接收），它们可能会操作同一个变量，因此需要加入读写锁。目前项目中无死锁存在。第二个比较重要的点是区块链的更新，所有节点共同、及时地切换到最新的节点，避免浪费算力。比较麻烦的是实现攻击者。其实一个简单的攻击者也很好实现，就是不在最长链上挖矿就行了。但是，当一群攻击者各自为营时，这种攻击是低效的，不可能成功的。因此，需要设计一种通信协议，让所有攻击者的行动统一起来。还有一种人是系统的监听者，顾名思义，他只接收信息而不产出信息，对区块链目前的状况做各种分析。

为了让区块的生产速度尽可能快地调整在 1s 左右，同时又避免受挖矿速度随机性的影响，我将区块调整难度的时机定为 5 个区块一次。可以看到图 1 中曲线和真实的比特币的挖矿曲线很相近，证明了实验的成功。其次，由于多线程本质还是在有限的 CPU 核（本实验中 4 核）上跑的，因此减少节点的数量并不一定代表着减少算力，对难度值并没有太大的变化，只是系统的不确定性减小了。

增加攻击者的比例，确实可以让成功的比例增大。同时，我观察到，系统对于一定数量的攻击者其实是有较好的抵御能力的，一旦突破了这个限制，会导致系统的安全性急剧下降（本实验中位于 10%-20% ）。

需要指出，本试验结果存在一定的随机性，但我已经尽可能选取最常见的试验结果，并使用大量统计的方法得到实验结论。



## 7 Lesson

其实在做这个项目之前我一点都不会 Go 语言，可以说在做这个项目之后我还是不能算会，但是我依然能够简单实现我想要的功能，体现了 Go 确实“简单、可靠且高效”。goroutine、mutex、channel 应有尽有，写分布式仿真代码确实很爽。遗憾是我也不知道为啥这个项目就跑起来了，有没有漏洞、死锁、不合理的假设？接口 Interface 怎么做模版？struct 继承？太多问题没来得及思考就把项目写完了。

本次实验给我最大的收获就是不再害怕 Go 语言，对于区块链有了一定的浅层理解，提升了一些能力。以后我也会认真对待每一个实验、每一次机会，不断提升自己，不断体会编程的快乐。



## 8 Acknowledgement

感谢老师在课上讲解的知识，以及提供的仿真模版，它们对我如有神助！

<br/>

本次实验所有代码均会在实验截止日期后开源 https://github.com/Harry-hhj/BlockChain_PoW_Simulation ，供需要学习、借鉴的人参考！