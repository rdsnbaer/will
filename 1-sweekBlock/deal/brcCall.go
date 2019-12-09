package deal

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"

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
	from = "0xe79fead329b69540142fabd881099c04424cc49f"
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

// type Coin struct {
// 	CoinName, CoinSymbol, CoinDecimals, CoinTotalSupply string
// }

func DealBrcCallInfo(_from, _to, _ip, _blockNumber string) CoinInfo {
	var coin CoinInfo

	_name, err := GetName(_from, _to, _ip)
	if err == nil {
		coin.Name = _name
	}

	_symbol, err := GetSymbol(_from, _to, _ip)
	if err == nil {
		coin.Symbol = _symbol
	}

	_decimals, err := GetDecimals(_from, _to, _ip, _blockNumber)
	if err == nil {
		coin.Precision, _ = strconv.ParseInt(_decimals, 10, 64)
	}

	_totalSupply, err := GetTotalSupply(_from, _to, _ip, _blockNumber)
	if err == nil {
		coin.Total = _totalSupply
	}

	return coin
}

func GetName(_from, _to, _ip string) (string, error) {
	_n := BrcCall(_from, _to, name, _ip)
	_name, err := AnalysisBrcCallResult(_n, Name)
	if err != nil {
		return _name, err
	}

	return _name, nil
}

func GetSymbol(_from, _to, _ip string) (string, error) {
	_s := BrcCall(_from, _to, symbol, _ip)
	_symbol, err := AnalysisBrcCallResult(_s, Symbol)
	if err != nil {
		return _symbol, err
	}

	return _symbol, nil
}

func GetDecimals(_from, _to, _ip, _blockNumber string) (string, error) {
	_d := BrcCaller(_from, _to, decimals, _ip, _blockNumber)
	_decimals, err := AnalysisBrcCallResult(_d, Decimals)
	if err != nil {
		return _decimals, err
	}

	return _decimals, nil
}

func GetTotalSupply(_from, _to, _ip, _blockNumber string) (string, error) {
	_t := BrcCaller(_from, _to, totalSupply, _ip, _blockNumber)
	_totalSupply, err := AnalysisBrcCallResult(_t, TotalSupply)
	if err != nil {
		return _totalSupply, err
	}

	return _totalSupply, nil
}

func GetOwner(_from, _to, _ip, _blockNumber string) (string, error) {
	_w := BrcCaller(_from, _to, owner, _ip, _blockNumber)
	_owner, err := AnalysisBrcCallResult(_w, Owner)
	if err != nil {
		return _owner, err
	}

	return _owner, nil
}

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

func AnalysisBrcCallResult(recv []byte, call Call) (string, error) {
	switch call {
	case Zero:
		err := errors.New("Zero error")
		return IsErr, err
	case Name:
		_name := IsString(recv)
		if _name == IsErr {
			err := errors.New("name is nil")
			return IsErr, err
		}
		fmt.Println("token name:", _name)
		return _name, nil
	case Symbol:
		_symbol := IsString(recv)
		if _symbol == IsErr {
			err := errors.New("symbol is nil")
			return IsErr, err
		}
		fmt.Println("token symbol:", _symbol)
		return _symbol, nil
	case TotalSupply:
		_totalSupply, _ := IsNumber(recv)
		if _totalSupply == IsErr {
			err := errors.New("totalSupply is nil")
			return IsErr, err
		}
		_total := fmt.Sprintf("%v", _totalSupply)
		fmt.Println("token totalSupply:", _total)
		return _total, nil
	case Decimals:
		_decimals, _ := IsNumber(recv)
		if _decimals == IsErr {
			err := errors.New("decimals is nil")
			return IsErr, err
		}
		_precision := fmt.Sprintf("%v", _decimals)
		fmt.Println("token decimals:", _decimals)
		return _precision, nil
	case Balance:
		_balance, _ := IsNumber(recv)
		if _balance == IsErr {
			err := errors.New("name is nil")
			return IsErr, err
		}
		fmt.Println("token balance:", _balance)
		return _balance, nil
	case Owner:
		_owner := IsAddress(recv)
		if _owner == IsErr {
			err := errors.New("name is nil")
			return IsErr, err
		}
		fmt.Println("token owner:", _owner)
		return _owner, nil
	default:
		fmt.Println("UnKnow!")
		return IsErr, nil
	}
}
