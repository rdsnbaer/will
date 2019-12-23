package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/my/repo/2-saokuai/deal"
	"github.com/my/repo/2-saokuai/logging"
)

const (
	ipJava = "http://127.0.0.1:8081"
	ipNet  = "http://127.0.0.1:8082"
	hash1  = "0x0a649348dec7300f180786ba1b1e5b9097b8f3ed4848d390ab5243f49f90f3fa"

	coinListStr     = "/coinList?"
	transferListStr = "/transferList?"
	coinInfoStr     = "/coinInfo?"

	pageNoStr   = "pageNo="
	pageSizeStr = "pageSize="
	symbolStr   = "symbol="
	fixStr      = "&"
)

var (
	pageNo, pageSize int64
	SIZE             int64 = 20
	NUM              int64 = 1
	IsOver           bool  = false
)

func main() {
	fmt.Println("start>>>")

OuterLoop:
	for {
		checkDeployContract(NUM)
		NUM++

		if IsOver {
			fmt.Println("16-checkDeployContract OuterLoop")
			logging.Info("16-checkDeployContract OuterLoop")
			break OuterLoop
		}
	}

	// symbol := "cba"
	// checkTransfer(symbol)
}

func checkTransfer(symbol string) {
	var size int64 = 20
	var num int64 = 1
	var isOver bool = false

OuterLoops:
	for {
		if isOver {
			fmt.Println("2.8-transfers OuterLoops")
			logging.Info("2.8-transfers OuterLoops")
			break OuterLoops
		}

		transferList := transferListHttp(symbol, num, size)
		recv, err := postMsg(transferList)
		if err != nil {
			fmt.Println("2.2-postMsg err:", err, ", [*transferList*]:", transferList)
			logging.Error("2.2-postMsg err:", err, ", [*transferList*]:", transferList)
			fmt.Println()
			return
		}

		res, err := deal.JsonTransfer(recv)
		if err == nil {
			lens := len(res)
			if lens == 0 {
				isOver = true
				fmt.Println("2.4-transfers is over")
				logging.Info("2.4-transfers is over")
				fmt.Println()

				return
			} else {
				deal.DealTransferInfo(res, ipNet)
				num++
			}
		} else {
			fmt.Println("2.3-JsonTransfer err:", err)
			logging.Info("2.3-JsonTransfer err:", err)
			fmt.Println()
		}
	}
}

func checkDeployContract(num int64) {
	_coinList := coinListHttp(num, SIZE)
	res, err := postMsg(_coinList)
	if err != nil {
		fmt.Println("2-postMsg err:", err, ", [*_coinList*]:", _coinList)
		logging.Error("2-postMsg err:", err, ", [*_coinList*]:", _coinList)

		return
	}

	_coins, err := deal.JsonCoinsInfo(res)
	if err == nil {
		getCoinCreateInfo(_coins)
	} else {
		fmt.Println("3-jsonCoin err:", err)
		logging.Error("3-jsonCoin err:", err, ", [*res*]:", string(res))
		fmt.Println()
	}
}

func getCoinCreateInfo(_coins []interface{}) {
	lens := len(_coins)
	if lens == 0 {
		IsOver = true
		fmt.Println("4-coinList is over")
		logging.Info("4-coinList is over")
		fmt.Println()

		return
	}

	for i := 0; i < lens; i++ {
		coinInfo, err := deal.AnalysisJsonCoin(_coins[i].(map[string]interface{}))
		if err == nil {
			fmt.Println("4-deploySuccessIsOne:", coinInfo)
			getContractCreateStatus(coinInfo)
		} else {
			fmt.Println("5-deploySuccessIsZero:", err, ",[*coinInfo*]", coinInfo)
			getContractCreateStatus(coinInfo)
		}
	}
}

func getContractCreateStatus(coinInfo deal.CoinInfo) {
	_coin := coinInfoHttp(coinInfo.Symbol)

	// Get(_coin)
	// time.Sleep(1000 * time.Second)
	recv, err := postMsg(_coin)
	//fmt.Println("3.1-:", string(recv))
	if err == nil {
		_hash, err := getCoinCreateHash(recv)
		if err == nil {
			if _hash != nil {
				res := deal.BrcGetTransactionReceipt(_hash, ipNet)
				fmt.Println("8-txHash:", _hash)
				logging.Info("8-txHash:", _hash)

				getContractCreateInfo(res, coinInfo, ipNet)
			} else {
				fmt.Println("9-create contract err, txHash:", _hash, ", [*coinInfo*]:", coinInfo)
				logging.Warn("9-create contract err, txHash:", _hash, ", [*coinInfo*]:", coinInfo)
				fmt.Println()
			}
		} else {
			fmt.Println("16-find hash err:", err, "_coin:", _coin)
			logging.Error("16-find hash err:", err, "_coin:", _coin)
			fmt.Println()
		}

	} else {
		fmt.Println("14-postMsg err:", err, "_coin:", _coin)
		logging.Error("14-postMsg err:", err, "_coin:", _coin)
		fmt.Println()
		return
	}
}

func getContractCreateInfo(recv []byte, coinInfo deal.CoinInfo, _ip string) {
	_res, err := deal.JsonMaps(recv)

	if err == nil {
		_coin, err := deal.GetContractCreateInfos(_res, coinInfo, _ip)
		if err != nil {
			fmt.Println(err)
		} else {
			checkTransfer(_coin.Symbol)
		}
	} else {
		fmt.Println("10-contract create err, result:", err, ", [*coinInfo*]:", coinInfo)
		logging.Error("10-contract create err, result:", err, ", [*coinInfo*]:", coinInfo)
		fmt.Println()
	}
}

func getCoinCreateHash(recv []byte) (interface{}, error) {
	_coin, err := deal.JsonCoinInfo(recv)
	if err != nil {
		return nil, err
	}
	_hash := _coin["txHash"]

	return _hash, nil
}

func Get(url string) ([]byte, error) {
	client := http.Client{Timeout: 5 * time.Second}
	response, err := client.Get(url)
	defer response.Body.Close()
	if err != nil {
		fmt.Println("111-err:", err)
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll err=", err)
		return nil, err
	}

	return body, err
}

func postMsg(msg string) ([]byte, error) {
	client := http.Client{Timeout: 5 * time.Second}
	response, err := client.Get(msg)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll err=", err)
		return nil, err
	}

	return body, err
}

func coinListHttp(_pageNo, _pageSize int64) string {
	_pageNoStr := pageNoStr + fmt.Sprintf("%d", _pageNo)
	_pageSizeStr := pageSizeStr + fmt.Sprintf("%d", _pageSize)
	coinList := ipJava + coinListStr + _pageNoStr + fixStr + _pageSizeStr

	fmt.Println("1-coinListHttp:", coinList)

	return coinList
}

func transferListHttp(_symbol string, _pageNo, _pageSize int64) string {
	_s := url.PathEscape(_symbol)

	_symbolStr := symbolStr + _s
	_pageNoStr := pageNoStr + fmt.Sprintf("%d", _pageNo)
	_pageSizeStr := pageSizeStr + fmt.Sprintf("%d", _pageSize)

	transferList := ipJava + transferListStr + _symbolStr + fixStr + _pageNoStr + fixStr + _pageSizeStr

	fmt.Println("2.1-transferListHttp:", transferList)
	logging.Info("2.1-transferListHttp:", transferList)

	return transferList
}

func coinInfoHttp(_symbol string) string {
	_s := url.PathEscape(_symbol)

	_symbolStr := symbolStr + _s
	coinInfo := ipJava + coinInfoStr + _symbolStr

	fmt.Println("6-coinInfoHttp:", coinInfo)

	return coinInfo
}
