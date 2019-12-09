package rpc

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/my/repo/2-saokuai/deal"
	"github.com/my/repo/2-saokuai/logging"
)

func Aaaa(v *big.Int) {
	AA(v)
	Value(v)
	AA(v)
}

func Value(v *big.Int) {
	fmt.Println("1-v:", v)
	v = v.Add(v, v)
	fmt.Println("2-v:", v)
}

func AA(v *big.Int) {
	fmt.Println("3-v:", v)
}

func GetInfo(res []byte, bNum int64, _ip, _conAddr string, _totals, _sum *big.Int) {
	var jsonMap map[string]interface{}
	err := json.Unmarshal(res, &jsonMap)
	if err != nil {
		fmt.Println("json unmarshal faild!")
		return
	}

	jsonResult := jsonMap["result"].(map[string]interface{})
	jsonTx := jsonResult["transactions"].([]interface{})
	if len(jsonTx) != 0 {
		analysisTx(jsonTx, bNum, _conAddr, _ip, _totals, _sum)
	} else {
		fmt.Println("the block has no transaction:", bNum)
	}
}

func analysisTx(jsonTX []interface{}, bNum int64, conAddr, _ip string, _totals, _sum *big.Int) {
	txCount := len(jsonTX)

	for i := 0; i < txCount; i++ {
		tx := jsonTX[i].(map[string]interface{})
		//fmt.Println(tx)
		if tx["to"] != nil {
			if tx["to"].(string) == conAddr {
				if tx["hash"] != nil {
					hash := tx["hash"]
					_from := tx["from"].(string)
					recv := BrcGetTransactionReceipt(hash, _ip)
					jsonResult := JsonMaps(recv)
					transferInfo(jsonResult, _from, conAddr, _ip, _totals, _sum)
				}
			} else {
				fmt.Printf("The %d deal is not a contract deal\n", i)
			}
		} else {
			fmt.Println("deploy contract")
			if tx["hash"] != nil {
				hash := tx["hash"]
				_from := tx["from"].(string)
				recv := BrcGetTransactionReceipt(hash, _ip)
				jsonResult := JsonMaps(recv)
				addr := jsonResult["contractAddress"].(string)

				if addr == conAddr {
					input := tx["input"].(string)
					analysisInputAndNum(input, _from, conAddr, _ip, bNum, _totals, _sum)
				} else {
					//contractAddr = jsonResult["to"].(string)
					logging.Debug("NeedContractAddr:", contractAddr, "TheContractAddr:", addr,
						"blockNumber:", bNum)
				}
			}
		}
	}
}

func JsonMaps(recv []byte) map[string]interface{} {
	var jsonMap map[string]interface{}

	err := json.Unmarshal(recv, &jsonMap)
	if err != nil {
		return nil
	}

	jsonResult := jsonMap["result"].(map[string]interface{})

	return jsonResult
}

func transferInfo(jsonResult map[string]interface{}, _from, _conAddr, _ip string,
	_totals, _sum *big.Int) {

	jsonLogs := jsonResult["logs"].([]interface{})

	num := len(jsonLogs)
	if num != 0 {
		for i := 0; i < num; i++ {
			tx := jsonLogs[i].(map[string]interface{})
			if tx["topics"] != nil {
				topic := tx["topics"].([]interface{})
				tp := topic[0].(string)

				getTransferInfo(tp, _from, _conAddr, _ip, tx, topic, _totals, _sum)
			}
		}
	}
}

func analysisInputAndNum(input, _from, _conAddr, _ip string, bNum int64,
	_totals, _sum *big.Int) {

	_d := input[15780:15844]
	_t := input[15844:15908]

	_decimals := deal.GetTxNum(_d)
	_total := deal.GetTxNum(_t)

	_v := new(big.Int).Exp(big.NewInt(10), _decimals, nil)
	_total = _total.Mul(_total, _v)

	str := fmt.Sprintf("%d", bNum)
	_own := BrcCaller(_from, _conAddr, owner, _ip, str)
	_owner := deal.IsAddress(_own)

	recv := BrcCaller(_from, _conAddr, totalSupply, _ip, str)
	_, n := deal.IsNumber(recv)

	_bo := _balance + _owner
	_bOwner := BrcCaller(_owner, _conAddr, _bo, _ip, str)
	_, bOwner := deal.IsNumber(_bOwner)

	com := _total.Cmp(n)
	coms := _total.Cmp(bOwner)

	if com == 0 && coms == 0 {
		_totals = _total
		_sum = bOwner

		logging.Info("blockNumber:", bNum, ",contractTotalSupply:", n,
			",contractAddress:", _conAddr, ",contractOwner:", _owner, ",balanceOwner:", bOwner)
		log.Println("contractAddr:", _conAddr, ",totalsupply:", n, ",blockNumer:", bNum)
	} else {
		logging.Error("blockNumber", bNum, "contractAddr:", _conAddr, ",returnTotalSupply:", n,
			",constructorToken:", _total)
		log.Fatal("contract deploy err")
	}

}

func getTransferInfo(tp, _from, _conAddr, _ip string, tx map[string]interface{},
	topic []interface{}, _totals, _sum *big.Int) {

	if tx["data"] != nil {
		txStr := tx["data"].(string)
		blockNum := uint64(tx["blockNumber"].(float64))
		//fmt.Println("1-blockNum", blockNum)
		t := new(token)

		switch tp {
		case transfer:
			txNum := deal.GetTxNum(txStr)
			t.transferToken(txNum, _totals, _sum, topic, blockNum, _conAddr, _ip)
		case increaseSupply:
			txNum := deal.GetTxNum(txStr)
			t.increaseToken(txNum, _totals, _sum, topic, blockNum, _conAddr, _ip)
		case burn:
			txNum := deal.GetTxNum(txStr)
			t.burnToken(txNum, _totals, _sum, topic, blockNum, _conAddr, _ip)
		case approve:
			txNum := deal.GetTxNum(txStr)
			t.approveToken(txNum, topic, blockNum, _conAddr, _ip)
		case changeContractOwn:
			t.changeContractOwner(topic, blockNum, _conAddr, _ip)
		default:
			fmt.Println("Unknow!")
		}
	}
}
