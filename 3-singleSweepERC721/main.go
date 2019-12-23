package main

import (
	"fmt"

	"github.com/my/repo/4-other/3-ERC721/derc"
	"github.com/my/repo/4-other/3-ERC721/loginfo"
)

const (
	blockNum = "8590137" // changeContractOwner
	ip              = "http://10.150.18.231:8031"
	contractAddress = "0x6067e36eb361f8f9a355c744fe98aae3d2bdb92e" 
)

func main() {
	doERC721()
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
