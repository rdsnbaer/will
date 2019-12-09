package deal

import (
	"errors"
	"fmt"

	"github.com/my/repo/2-saokuai/logging"
)

const (
	DeAddr = "0x0000000000000000000000000000000000000000"
	//CreateError = "-contract create err"
)

type Contract struct {
	Precision, BlockNumber                         int64
	Total, Name, Symbol, From, Owner, ContractAddr string
}

func GetContractCreateInfos(res map[string]interface{}, coinInfo CoinInfo, _ip string) (CoinInfo, error) {
	var contract Contract

	_to := res["to"].(string)
	if DeAddr == _to {
		_contractAddr := res["contractAddress"].(string)
		_from := res["from"].(string)
		n := res["blockNumber"].(float64)
		num := int64(n)
		blockNum := fmt.Sprintf("%d", num)

		_coinInfo := DealBrcCallInfo(_from, _contractAddr, _ip, blockNum)

		isok := IsCreateContract(_coinInfo, coinInfo)
		if isok {
			contract.Precision = _coinInfo.Precision
			contract.BlockNumber = num
			contract.Total = _coinInfo.Total
			contract.Name = _coinInfo.Name
			contract.Symbol = _coinInfo.Symbol
			contract.From = _from
			contract.ContractAddr = _contractAddr

			_owner, err := GetOwner(_from, _contractAddr, _ip, blockNum)
			if err == nil {
				contract.Owner = _owner
			}

			fmt.Println("12-create contract success,", "name:", contract.Name, "synbol:", contract.Symbol,
				"decimials:", contract.Precision, "totalSupply:", contract.Total, "from:", contract.From,
				"owner:", contract.Owner, "contractAdd:", contract.ContractAddr,
				"blockNumber:", contract.BlockNumber)
			logging.Info("12-create contract success,", "name:", contract.Name, "synbol:", contract.Symbol,
				"decimials:", contract.Precision, "totalSupply:", contract.Total, "from:", contract.From,
				"owner:", contract.Owner, "contractAdd:", contract.ContractAddr,
				"blockNumber:", contract.BlockNumber)

			fmt.Println("13-create success:", _coinInfo)
			logging.Info("13-create success:", _coinInfo)
			fmt.Println()

			return _coinInfo, nil
		} else {
			fmt.Println("11-Plan create contract: ", coinInfo, ", [*actual*]: ", _coinInfo)
			logging.Error("11-Plan create contract: ", coinInfo, ", [*actual*]: ", _coinInfo)
			fmt.Println()
			return _coinInfo, errors.New("err")
		}

	} else {
		fmt.Println("14-to is not:", DeAddr, ", **not create contract:", coinInfo)
		logging.Error("14-to is not:", DeAddr, ", **not create contract:", coinInfo)
		fmt.Println()

		return coinInfo, errors.New("err")
	}
}

func IsCreateContract(_coinInfo, coinInfo CoinInfo) (isOk bool) {
	isOk = _coinInfo.Precision == coinInfo.Precision && _coinInfo.Total == coinInfo.Total && _coinInfo.Name == coinInfo.Name && _coinInfo.Symbol == coinInfo.Symbol

	return isOk
}

func DealReceipt(jsonResult map[string]interface{}) (string, string, int64) {
	_from := jsonResult["from"].(string)
	_contractAddress := jsonResult["contractAddress"].(string)
	_blockNum := jsonResult["blockNumber"].(float64)
	_blockNumber := int64(_blockNum)

	return _from, _contractAddress, _blockNumber
}
