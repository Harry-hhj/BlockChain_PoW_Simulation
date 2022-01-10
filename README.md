# Blockchain PoW Simulation

## 1 Introduction

This project implements a PoW-based blockchain and simulates the existence of a group of  fork attack attackers.

Final project for SJTU IS416.

## 2 Detail

See `Final Lab BlockChain-PoW.md` for details. (Chinese avalable only)

## 3 Usage

Open a terminal, enter:

```bash
go build
.pow
```

If you want to change the output of the program, modify the `handler` function of the class monitor under the `node.go` file:

```go
func (m *Monitor) handler(msg Block)
```

## 4 Remark

There are still some problems with the program, such as

- [ ] The verification function does not verify the difficulty of the blockchain
- [ ] The adjustment accuracy of the difficulty value is not enough

If you find any problem, welcome to contact me.

<br/>

**_If you think this project is helpful to you, please click star, it is very useful to me._**
**_如果你觉得本项目对你有所帮助，欢迎点star，这对我很有用。_**

Note: This project is just a final assignment of a school curriculum, and there are many things that are not good enough. It is only for novice learning and schoolmates' reference. It is not recommended to use it in any of your projects.

