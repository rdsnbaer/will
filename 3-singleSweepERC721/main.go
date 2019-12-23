package main

import (
	"fmt"
	"strconv"

	"github.com/my/repo/4-other/3-ERC721/derc"
	"github.com/my/repo/4-other/3-ERC721/loginfo"
)

const (
	//blockNum = "8572873" // 部署合约
	//blockNum        = "8573401" // mint
	//blockNum = "8579572" //setApprovalForAll
	//blockNum = "8573401" // changeContractOwner
	//blockNum = "8575091" // GetApproveInfo
	blockNum = "1"


	ip              = "http://10.150.18.231:8031"


	//contractAddress = "0x6067e36eb361f8f9a355c744fe98aae3d2bdb92e" //测试网
	contractAddress = "0x8ad0955c90382d6f9948f331841be47abd938af3" //
)

func main() {
	//doERC721()
	loopERC721(blockNum, contractAddress, ip)
}

func loopERC721(bNum, addr, ip string) {
	num, _ := strconv.Atoi(bNum)

	for i := num; ; i++ {
		n := fmt.Sprintf("%d", i)
		recv, err := derc.BrcGetBlockDetialByNumber(n, ip)
		if err == nil {
			res, err := derc.PostMsg(recv, ip)
			if err == nil {
				jsonMap, err := derc.GetUnmarshal(res)
				if err == nil {
					txs := derc.GetTransactions(jsonMap)
					if txs != nil {
						derc.DealTransactions(txs, ip, contractAddress)
					} else {
						fmt.Println("1.4-GetTransactions result or transactions is nil, blockNumber:", n)
					}
				} else {
					fmt.Println("1.3-GetUnmarshal err", err)
					loginfo.Error("1.3-GetUnmarshal err", err)
				}
			} else {
				fmt.Println("1.2-PostMsg err:", err)
				loginfo.Error("1.2-PostMsg err:", err)
			}
		} else {
			fmt.Println("1.1-BrcGetBlockDetialByNumber err:", err)
			loginfo.Error("1.1-BrcGetBlockDetialByNumber err:", err)
		}
	}
}

func doERC721() {
	recv, err := derc.BrcGetBlockDetialByNumber(blockNum, ip)
	if err == nil {
		res, err := derc.PostMsg(recv, ip)
		if err == nil {
			jsonMap, err := derc.GetUnmarshal(res)
			if err == nil {
				txs := derc.GetTransactions(jsonMap)
				if txs != nil {
					derc.DealTransactions(txs, ip, contractAddress)
				} else {
					fmt.Println("1.4-GetTransactions result or transactions is nil")
				}
			} else {
				fmt.Println("1.3-GetUnmarshal err", err)
				loginfo.Error("1.3-GetUnmarshal err", err)
			}
		} else {
			fmt.Println("1.2-PostMsg err:", err)
			loginfo.Error("1.2-PostMsg err:", err)
		}
	} else {
		fmt.Println("1.1-BrcGetBlockDetialByNumber err:", err)
		loginfo.Error("1.1-BrcGetBlockDetialByNumber err:", err)
	}
}
