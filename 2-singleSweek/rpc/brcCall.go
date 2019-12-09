package rpc

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/my/repo/2-saokuai/deal"
	"github.com/my/repo/2-saokuai/http"
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

var (
	// Operators calling or deploying contracts
	from = "0xe79fead329b69540142fabd881099c04424cc49f"
	// contract address
	//to             = "0x22c614bb0fc26bbdc90bdb88c2df87d377917d5c"
	deContractAddr = "0x0000000000000000000000000000000000000000"

	contractAddr string
	totalToken   *big.Int
)

type Call int

const (
	Zero Call = iota
	TotalSupply
	Decimals
	Balance
	Name
	Symbol
	Owner
)

func BrcCall(from, to, data, ip string) []byte {
	return BrcCaller(from, to, data, ip, "latest")
}

func BrcCaller(from, to, data, ip, blockNum string) []byte {
	msg := make(map[string]interface{})
	param := make(map[string]interface{})

	msg["jsonrpc"] = "2.0"
	msg["method"] = "brc_call"
	param["from"] = from
	param["to"] = to
	param["data"] = data
	msg["params"] = []interface{}{param, blockNum}
	msg["id"] = 1

	mjson, _ := json.Marshal(msg)
	mString := string(mjson)
	recv := http.PostMsg(mString, ip)

	return recv
}

func AnalysisBrcCallResult(recv []byte, call Call) {
	switch call {
	case Zero:
		return
	case Name:
		_name := deal.IsString(recv)
		fmt.Println("token name:", _name)
	case Symbol:
		_symbol := deal.IsString(recv)
		fmt.Println("token symbol:", _symbol)
	case TotalSupply:
		_totalSupply, _ := deal.IsNumber(recv)
		fmt.Println("token totalSupply:", _totalSupply)
	case Decimals:
		_decimals, _ := deal.IsNumber(recv)
		fmt.Println("token decimals:", _decimals)
	case Balance:
		_balance, _ := deal.IsNumber(recv)
		fmt.Println("token balance:", _balance)
	case Owner:
		_owner := deal.IsAddress(recv)
		fmt.Println("token owner:", _owner)
	default:
		fmt.Println("Unknow!")
	}
}
