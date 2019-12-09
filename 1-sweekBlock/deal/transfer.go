package deal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/my/repo/2-saokuai/logging"
)

const (
	TRANSFER = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	SYMBOL   = "0x95d89b41"
)

type TransferInfo struct {
	blockHeight                                       int64
	from, to, amount, contractAddress, txHash, symbol string
}

func DealTransferInfo(transfers []interface{}, ip string) {
	lens := len(transfers)
	for i := 0; i < lens; i++ {
		transferInfo := GetTransferInfo(transfers[i].(map[string]interface{}))
		val, err := GetReceiptResult(transferInfo, ip)
		if err == nil {
			_transferInfo, err := GetResult(val, ip)
			if err == nil {
				if transferInfo == _transferInfo {
					fmt.Println("2.7-transferInfo is correct:", transferInfo)
					logging.Info("2.7-transferInfo is correct:", transferInfo)
					fmt.Println()
				} else {
					fmt.Println("2.8-transferInfo is wrong, Plan transfer:", transferInfo, ", actual:", _transferInfo)
					logging.Error("2.8-transferInfo is wrong, Plan transfer:", transferInfo, ", actual:", _transferInfo)
					fmt.Println()
				}
			} else {
				fmt.Println("2.6-transferInfo not same GetReceiptResult:", transferInfo, ",[*err*]:", err)
				logging.Error("2.6-transferInfo not same GetReceiptResult:", transferInfo, ",[*err*]:", err)
				fmt.Println()
			}
		} else {
			fmt.Println("2.5-JsonMaps err, result is:", err, ",[*GetTransferInfo*]:", transferInfo)
			logging.Error("2.5-JsonMaps err, result is:", err, ",[*GetTransferInfo*]:", transferInfo)
			fmt.Println()
		}
	}
}

func GetResult(res map[string]interface{}, ip string) (TransferInfo, error) {
	var _transferInfo TransferInfo
	var tx map[string]interface{}

	logs := res["logs"].([]interface{})
	lens := len(logs)

	if lens == 0 {
		return _transferInfo, errors.New("err, not transfer")
	}
	if lens == 1 {
		tx = logs[0].(map[string]interface{})
	}
	if lens == 2 {
		tx = logs[1].(map[string]interface{})
	}

	isOk, err := checkLogsInfo(res, tx)
	if !isOk {
		return _transferInfo, err
	}

	return getLogsInfo(res, tx, ip)
}

func getLogsInfo(res, tx map[string]interface{}, ip string) (TransferInfo, error) {
	var _transferInfo TransferInfo

	n := tx["blockNumber"].(float64)
	_transferInfo.blockHeight = int64(n)

	topics := tx["topics"].([]interface{})
	_transferInfo.from = "0x" + topics[1].(string)[26:]
	_transferInfo.to = "0x" + topics[2].(string)[26:]

	data := tx["data"].(string)
	_transferInfo.amount = RemoveZero(data[2:])

	_transferInfo.contractAddress = res["contractAddress"].(string)
	_transferInfo.txHash = tx["transactionHash"].(string)

	recv, err := GetSymbol(_transferInfo.from, _transferInfo.contractAddress, ip)
	if err != nil {
		return _transferInfo, err
	}

	_transferInfo.symbol = recv

	return _transferInfo, nil

}

func checkLogsInfo(res, tx map[string]interface{}) (bool, error) {
	_resBlockNumber := res["blockNumber"].(float64)
	resBlockNumber := int64(_resBlockNumber)
	_txBlockNumber := tx["blockNumber"].(float64)
	txBlockNumber := int64(_txBlockNumber)
	if resBlockNumber != txBlockNumber {
		r := fmt.Sprintf("%d", resBlockNumber)
		t := fmt.Sprintf("%d", txBlockNumber)
		err := "blockNumber not the same, resBlockNumber:" + r + ", txBlockNumber" + t
		return false, errors.New(err)
	}

	resHash := res["transactionHash"].(string)
	txHash := tx["transactionHash"].(string)
	if resHash != txHash {
		err := "transactionHash not the same, resHash:" + resHash + ", xHash" + txHash
		return false, errors.New(err)
	}

	resContractAddr := res["to"].(string)
	txContractAddr := tx["address"].(string)
	if resContractAddr != txContractAddr {
		err := "contractAddr is not same, resContractAddr:" + resContractAddr + ", txContractAddr:" + txContractAddr
		return false, errors.New(err)
	}

	if tx["topics"] == nil {
		return false, errors.New("topics is nil")
	}

	topics := tx["topics"].([]interface{})
	if len(topics) != 3 {
		return false, errors.New("len topics is not 3")
	}

	topic1 := topics[0].(string)
	if topic1 != TRANSFER {
		err := "topics[0] is not same TRANSFER, topics[0]:" + topic1 + ", TRANSFER:" + TRANSFER
		return false, errors.New(err)
	}

	_txFrom := topics[1].(string)[26:]
	_resFrom := res["from"].(string)
	if _txFrom != _resFrom {
		err := "_txFrom is not _resFrom, _txFrom:" + _txFrom + ", _resFrom" + _resFrom
		return false, errors.New(err)
	}

	return true, nil
}

func GetReceiptResult(res TransferInfo, ip string) (map[string]interface{}, error) {
	recv := BrcGetTransactionReceipt(res.txHash, ip)
	val, err := JsonMaps(recv)

	return val, err
}

func GetTransferInfo(res map[string]interface{}) TransferInfo {
	var transferInfo TransferInfo

	num := res["blockHeight"].(float64)

	transferInfo.blockHeight = int64(num)
	transferInfo.from = res["from"].(string)
	transferInfo.to = res["to"].(string)
	transferInfo.amount = res["amount"].(string)
	transferInfo.contractAddress = res["contractAddress"].(string)
	transferInfo.txHash = res["txHash"].(string)
	transferInfo.symbol = res["symbol"].(string)

	return transferInfo
}

func JsonTransfer(recv []byte) ([]interface{}, error) {
	var jsonMap map[string]interface{}

	err := json.Unmarshal(recv, &jsonMap)
	if err == nil {
		return nil, err
	}

	res := jsonMap["transfers"].([]interface{})

	return res, nil
}
