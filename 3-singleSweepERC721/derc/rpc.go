package derc

import (
	"encoding/json"
	"fmt"

	"github.com/my/repo/4-other/3-ERC721/loginfo"
)

func BrcGetBlockDetialByNumber(blockNum, ip string) ([]byte, error) {
	jsonMap := make(map[string]interface{})

	jsonMap["jsonrpc"] = "2.0"
	jsonMap["method"] = "brc_getBlockDetialByNumber"
	jsonMap["params"] = []interface{}{blockNum, true}
	jsonMap["id"] = 1

	recv, err := json.Marshal(jsonMap)
	if err != nil {
		return recv, err
	}
	fmt.Println(string(recv))
	return recv, nil
}

func BrcGetTransactionReceipt(hash, ip string) ([]byte, error) {
	jsonMap := make(map[string]interface{})

	jsonMap["id"] = 1
	jsonMap["jsonrpc"] = "2.0"
	jsonMap["method"] = "brc_getTransactionReceipt"
	jsonMap["params"] = []interface{}{hash}

	recv, err := json.Marshal(jsonMap)
	if err != nil {
		return recv, err
	}

	fmt.Println("2.0-hash:", hash)
	loginfo.Info("2.0-hash:", hash)

	return recv, nil
}

func GetTransactions(jsonMap map[string]interface{}) []interface{} {
	jsonResult := jsonMap["result"].(map[string]interface{})
	if jsonResult == nil {
		return nil
	}

	txs := jsonResult["transactions"].([]interface{})
	if len(txs) == 0 {
		return nil
	}

	return txs
}

func DealTransactions(txs []interface{}, ip, contractAddr string) {
	num := len(txs)
	for i := 0; i < num; i++ {
		tx := txs[i].(map[string]interface{})
		hash := tx["hash"].(string)

		to := tx["to"]
		if to == nil {
			fmt.Println("2.5.1-deploy contract")
			jsonMap, err := DeployOrExecContract(hash, ip)
			if err == nil {
				DealDeployReceipt(jsonMap, contractAddr, ip)
			} else {
				fmt.Println()
			}
		} else {
			if to.(string) == contractAddr {
				fmt.Println("2.5.2-exec contract")
				jsonMap, err := DeployOrExecContract(hash, ip)
				if err == nil {
					DealExecContract(jsonMap, contractAddr, ip)
				} else {
					fmt.Println()
				}
			} else {
				fmt.Println("2.5.3-not deploy or exec contract")
			}
		}
	}
}

func DeployOrExecContract(hash, ip string) (map[string]interface{}, error) {
	recv, err := BrcGetTransactionReceipt(hash, ip)
	if err != nil {
		fmt.Println("2.1-BrcGetTransactionReceipt err:", err)
		loginfo.Error("2.1-BrcGetTransactionReceipt err:", err)
		return nil, err
	}

	res, err := PostMsg(recv, ip)
	if err != nil {
		fmt.Println("2.2-PostMsg err:", err)
		loginfo.Info("2.2-PostMsg err:", err)
		return nil, err
	}

	jsonMap, err := GetUnmarshal(res)
	if err != nil {
		fmt.Println("2.3-GetUnmarshal err:", err)
		loginfo.Error("2.3-GetUnmarshal err:", err)
		return nil, err
	}

	return jsonMap, nil
}

func DealDeployReceipt(jsonMap map[string]interface{}, contractAddr, ip string) {
	if jsonMap["result"] != nil {
		jsonResult := jsonMap["result"].(map[string]interface{})
		GetDeployContractInfo(jsonResult, contractAddr, ip)
	} else {
		fmt.Println("2.4-jsonResult is nil")
		loginfo.Error("2.4-jsonResult is nil")
	}
}

func GetDeployContractInfo(res map[string]interface{}, contractAddr, ip string) {
	num := res["blockNumber"].(float64)
	blockNum := int64(num)
	blockNumber := fmt.Sprintf("%d", blockNum)
	from := res["from"].(string)

	var q Query
	var contract Contract

	name, err := q.GetNFTName(from, contractAddr, blockNumber, ip)
	if err == nil {
		contract.Name_C = name
	} else {
		fmt.Println("4.1-GetNFTName err:", err)
		loginfo.Error("4.1-GetNFTName err:", err)
	}

	symbol, err := q.GetNFTSymbol(from, contractAddr, blockNumber, ip)
	if err == nil {
		contract.Symbol_C = symbol
	} else {
		fmt.Println("4.2-GetNFTSymbol err:", err)
		loginfo.Error("4.2-GetNFTSymbol err:", err)
	}

	contractOwner, err := q.GetContractOwner(from, contractAddr, blockNumber, ip)
	if err == nil {
		contract.ContractOwner_C = contractOwner
	} else {
		fmt.Println("4.3-GetContractOwner err:", err)
		loginfo.Error("4.3-GetContractOwner err:", err)
	}

	fmt.Println("4.3-create contract:", contract)
	loginfo.Info("4.3-create contract:", contract)
	loginfo.Info("-------------------------------------------------------------------")
	fmt.Println()
}

func DealExecContract(res map[string]interface{}, contractAddr, ip string) {
	if res["result"] != nil {
		jsonResult := res["result"].(map[string]interface{})
		if jsonResult["logs"] != nil {
			logs := jsonResult["logs"].([]interface{})
			len := len(logs)
			for i := 0; i < len; i++ {
				log := logs[i].(map[string]interface{})
				GetExecContractInfo(jsonResult, log, contractAddr, ip)
			}
		} else {
			fmt.Println("5-2:logs is nil")
		}
	} else {
		fmt.Println("5.1-jsonResult is nil")
		loginfo.Error("5.1-jsonResult is nil")
	}
}

func GetExecContractInfo(res, log map[string]interface{}, contractAddr, ip string) {
	if log["topics"] != nil {
		topic := log["topics"].([]interface{})
		DealResultAndTopic(res, log, topic, contractAddr, ip)
	} else {
		fmt.Println("5.3-topics is nil")
	}
}
