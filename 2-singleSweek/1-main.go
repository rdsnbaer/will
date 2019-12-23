package main

import (
	"math/big"

	"github.com/my/repo/2-saokuai/rpc"
)

const (
	ip = "http://127.0.0.1:8081"

	from = "0xe79fead329b69540142fabd881099c04424cc49f"

	//contractAddress = "0x8c78e9686442cd2fb8b3217d34d9e55c1b2171da"
	//contractAddress = "0x22c614bb0fc26bbdc90bdb88c2df87d377917d5c"
	contractAddress = "0x4d4db69111deed3782056f067bfe799a4efa8b4b"

	blockNum = 7732621
)

const (
	// Transfer fixed field
	transfer = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	// Add token fixed field
	increaseSupply = "0x430080bcaaf4832a6fedf62b6ce7849b50c273246fdc69a13ce54e9772008f1c"
	// Destroy token fixed field
	burn = "0x8102980f65c5fee2d21c925f22d022abc395e0c409616c78835cdf6d6dff613b"

	changeContractOwn = "0x07ff91ee303faa39fb867eb01411e233f560ef81a7eef1aff036a194236d1161"

	approve = "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925"

	// Parameters for querying the number of issued tokens
	totalSupply = "0x18160ddd"
	// Parameters to query token precision
	decimals = "0x313ce567"
	// Parameters to query the address balance of the contract token
	balance  = "0x70a08231000000000000000000000000" + "e79fead329b69540142fabd881099c04424cc49f"
	_balance = "0x70a08231000000000000000000000000"
	// Parameters for querying the name of a contract token
	name = "0x06fdde03"
	// Parameters for querying symbols of contract tokens
	symbol = "0x95d89b41"
	// Query contract owner parameters
	owner = "0x8da5cb5b"
)

var Totals *big.Int
var SumBalances *big.Int

func main() {
	dealAllContractCall()
	// Totals = big.NewInt(0)
	// SumBalances = big.NewInt(0)

	// var b int64
	// b = 5707670
	// for {
	// 	recv := rpc.BrcGetBlockDetialByNumber(b, ip)
	// 	rpc.GetInfo(recv, b, ip, contractAddress, Totals, SumBalances)
	// 	b++
	// }
}

func dealAllContractCall() {
	//file, _ := os.OpenFile("token.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)

	brcCall := rpc.BrcCall(from, contractAddress, owner, ip)
	rpc.AnalysisBrcCallResult(brcCall, rpc.Owner)

	brcCall = rpc.BrcCall(from, contractAddress, totalSupply, ip)
	rpc.AnalysisBrcCallResult(brcCall, rpc.TotalSupply)

	brcCall = rpc.BrcCall(from, contractAddress, decimals, ip)
	rpc.AnalysisBrcCallResult(brcCall, rpc.Decimals)

	b := _balance + "9e19f91a85ca1bccbe783ad8190a55534615c544"
	brcCall = rpc.BrcCall(from, contractAddress, b, ip)
	rpc.AnalysisBrcCallResult(brcCall, rpc.Balance)

	brcCall = rpc.BrcCall(from, contractAddress, name, ip)
	rpc.AnalysisBrcCallResult(brcCall, rpc.Name)

	brcCall = rpc.BrcCall(from, contractAddress, symbol, ip)
	rpc.AnalysisBrcCallResult(brcCall, rpc.Symbol)

	brcCall = rpc.BrcCall(from, contractAddress, owner, ip)
	rpc.AnalysisBrcCallResult(brcCall, rpc.Owner)
}
