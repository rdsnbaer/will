package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/my/repo/2-saokuai/logging"
)

const (
	ip = "http://47.74.235.122:33333"

	//contractAddress = "0x8c78e9686442cd2fb8b3217d34d9e55c1b2171da"
	contractAddress = "0x4d4db69111deed3782056f067bfe799a4efa8b4b"

	blockNum = 7732621

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
	precision = "0x313ce567"
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
	Precision
	Balance
	Name
	Symbol
	Owner
)

type tokenIce interface {
	transferToken(txNum interface{}, topic []interface{}, blockNum interface{})
	increaseToken(txNum interface{}, topic []interface{}, blockNum interface{})
	approveToken(txNum interface{}, topic []interface{}, blockNum interface{})
	burnToken(txNum interface{}, topic []interface{}, blockNum interface{})
	changeContractOwner(topic []interface{}, blockNum interface{})
}

type token struct {
}

func main() {
	// // Obtain the information returned by the corresponding block height
	recv := brcGetBlockDetialByNumber(blockNum)
	//fmt.Println(string(recv))
	// query contract transfer, increaseSupply token, burn token info
	getInfo(recv, blockNum)

	// call RPC (brc_call) interface
	brccall := brcCall(totalSupply)
	// query contract token totalSupply numbers
	analysisBrcCallResult(brccall, TotalSupply)

	brccall = brcCall(precision)
	// query contract token decimals
	analysisBrcCallResult(brccall, Precision)

	brccall = brcCall(balance)
	// query address balance
	analysisBrcCallResult(brccall, Balance)

	brccall = brcCall(name)
	// query contract token name
	analysisBrcCallResult(brccall, Name)

	brccall = brcCall(symbol)
	// query contract token symbol
	analysisBrcCallResult(brccall, Symbol)

	brccall = brcCall(owner)
	// query contract token owner
	analysisBrcCallResult(brccall, Owner)
}

// Analysis query return results
func analysisBrcCallResult(recv []byte, call Call) {
	switch call {
	case Zero:
		return
	case Name, Symbol:
		isString(recv)
	case TotalSupply, Precision, Balance:
		fmt.Println(isNumber(recv))
	case Owner:
		fmt.Println(isAddress(recv))
	default:
		fmt.Println("Unknow!")
	}
}

func isAddress(recv []byte) string {
	var jsonMap map[string]interface{}

	err := json.Unmarshal(recv, &jsonMap)
	if err != nil {
		logging.Error(err)
		log.Fatal("err:", err)
	}

	result := jsonMap["result"].(interface{})
	addr := result.(string)[26:]

	return addr
}

func isString(recv []byte) {
	var jsonMap map[string]interface{}
	err := json.Unmarshal(recv, &jsonMap)
	if err != nil {
		return
	}

	result := jsonMap["result"].(interface{})
	rString := result.(string)
	res := rString[66:]

	ln := removeZero(res[:64])
	len, _ := strconv.ParseUint(ln, 16, 64)
	index := len*2 + 64
	str := res[64:index]

	src := []byte(str)
	decode(src)
}

func decode(src []byte) {
	dstDeCode := make([]byte, hex.DecodedLen(len(src)))

	_, err := hex.Decode(dstDeCode, src)
	if err != nil {
		log.Fatal(err)
	}

	decodedStr, _ := hex.DecodeString(string(src))
	fmt.Println(string(decodedStr))
	fmt.Println()
}

func isNumber(recv []byte) *big.Int {
	var jsonMap map[string]interface{}
	err := json.Unmarshal(recv, &jsonMap)
	if err != nil {
		log.Fatal(err)
	}

	result := jsonMap["result"].(interface{})
	res := result.(string)
	supply := getTxNum(res)

	fmt.Println(string(recv), "txNum is:", supply)
	//logging.Info("txNum:", supply)
	fmt.Println()

	return supply
}

func brcCall(data string) []byte {
	return _brcCall(data, "latest")
}

func _brcCall(data, blockNum string) []byte {
	msg := make(map[string]interface{})
	param := make(map[string]interface{})

	msg["jsonrpc"] = "2.0"
	msg["method"] = "brc_call"
	param["from"] = from
	param["to"] = contractAddress
	param["data"] = data
	msg["params"] = []interface{}{param, blockNum}
	msg["id"] = 1

	mjson, _ := json.Marshal(msg)
	mString := string(mjson)
	recv := postMsg(mString, ip)

	return recv
}

func brcGetBlockDetialByNumber(num int64) []byte {
	msg := make(map[string]interface{})

	msg["jsonrpc"] = "2.0"
	msg["method"] = "brc_getBlockDetialByNumber"
	msg["params"] = []interface{}{toHex(num, true), true}
	msg["id"] = 1

	mjson, _ := json.Marshal(msg)
	mString := string(mjson)
	recv := postMsg(mString, ip)

	return recv
}

func toHex(num int64, prefix bool) string {
	if prefix {
		return "0x" + strconv.FormatInt(num, 16)
	} else {
		return strconv.FormatInt(num, 16)
	}
}

func getInfo(res []byte, bNum int64) {
	var jsonMap map[string]interface{}
	err := json.Unmarshal(res, &jsonMap)
	if err != nil {
		fmt.Println("json unmarshal faild!")
		return
	}

	jsonResult := jsonMap["result"].(map[string]interface{})
	jsonTx := jsonResult["transactions"].([]interface{})
	if len(jsonTx) != 0 {
		analysisTx(jsonTx, bNum)
	} else {
		fmt.Println("the block has no transaction")
	}
}

func analysisTx(jsonTX []interface{}, bNum int64) {
	txCount := len(jsonTX)

	for i := 0; i < txCount; i++ {
		tx := jsonTX[i].(map[string]interface{})
		//fmt.Println(tx)
		if tx["to"] != nil {
			if tx["to"].(string) == contractAddress {
				if tx["hash"] != nil {
					hash := tx["hash"]
					recv := brcGetTransactionReceipt(hash)
					jsonResult := jsonMaps(recv)
					transferInfo(jsonResult)
				}
			} else {
				fmt.Printf("The %d deal is not a contract deal\n", i)
			}
		} else {
			fmt.Println("deploy contract")
			if tx["hash"] != nil {
				hash := tx["hash"]
				recv := brcGetTransactionReceipt(hash)
				jsonResult := jsonMaps(recv)
				addr := jsonResult["contractAddress"].(string)

				if addr == contractAddress {
					contractAddr = jsonResult["from"].(string)

					input := tx["input"].(string)
					analysisInputAndNum(input, bNum)
				} else {
					logging.Debug("NeedContractAddr:", contractAddr, "TheContractAddr:", addr,
						"blockNumber:", bNum)
				}
			}
		}
	}
}

func analysisInputAndNum(input string, bNum int64) {
	_d := input[16208:16272]
	_t := input[16272:16336]

	_decimals := getTxNum(_d)
	_total := getTxNum(_t)

	_v := new(big.Int).Exp(big.NewInt(10), _decimals, nil)
	_total = _total.Mul(_total, _v)

	str := fmt.Sprintf("%d", bNum)
	_own := _brcCall(owner, str)
	_owner := isAddress(_own)

	recv := _brcCall(totalSupply, str)
	n := isNumber(recv)
	com := _total.Cmp(n)

	if com == 0 {
		logging.Info("blockNumber:", bNum, ",contractTotalSupply:", n,
			",contractAddress:", contractAddress, ",contractOwner:", _owner)
	} else {
		logging.Error("blockNumber", bNum, "contractAddr:", contractAddress, ",returnTotalSupply:", n,
			",constructorToken:", _total)
	}

	log.Println("contractAddr:", contractAddress, ",totalsupply:", n, ",blockNumer:", bNum)

}

func brcGetTransactionReceipt(h interface{}) []byte {
	msg := make(map[string]interface{})

	msg["jsonrpc"] = "2.0"
	msg["method"] = "brc_getTransactionReceipt"
	msg["params"] = []interface{}{h.(string)}
	msg["id"] = 1

	mjson, _ := json.Marshal(msg)
	mString := string(mjson)
	recv := postMsg(mString, ip)

	return recv
}

func jsonMaps(recv []byte) map[string]interface{} {
	var jsonMap map[string]interface{}

	err := json.Unmarshal(recv, &jsonMap)
	if err != nil {
		return nil
	}

	jsonResult := jsonMap["result"].(map[string]interface{})

	return jsonResult
}

func transferInfo(jsonResult map[string]interface{}) {
	jsonLogs := jsonResult["logs"].([]interface{})

	num := len(jsonLogs)
	if num != 0 {
		for i := 0; i < num; i++ {
			tx := jsonLogs[i].(map[string]interface{})
			if tx["topics"] != nil {
				topic := tx["topics"].([]interface{})
				tp := topic[0].(string)

				getTransferInfo(tp, tx, topic)
			}
		}
	}
}

func getTransferInfo(tp string, tx map[string]interface{}, topic []interface{}) {
	if tx["data"] != nil {
		txStr := tx["data"].(string)
		blockNum := uint64(tx["blockNumber"].(float64))
		//fmt.Println("1-blockNum", blockNum)
		t := new(token)

		switch tp {
		case transfer:
			txNum := getTxNum(txStr)
			t.transferToken(txNum, topic, blockNum)
		case increaseSupply:
			txNum := getTxNum(txStr)
			t.increaseToken(txNum, topic, blockNum)
		case burn:
			txNum := getTxNum(txStr)
			t.burnToken(txNum, topic, blockNum)
		case approve:
			txNum := getTxNum(txStr)
			t.approveToken(txNum, topic, blockNum)
		case changeContractOwn:
			t.changeContractOwner(topic, blockNum)
		default:
			fmt.Println("Unknow!")
		}
	}
}

func (t *token) approveToken(txNum *big.Int, topic []interface{}, blockNum uint64) {
	_from := "0x" + topic[1].(string)[26:]
	_to := "0x" + topic[2].(string)[26:]

	// _bFrom := _balance + _from[2:]
	// _bTo := _balance + _to[2:]

	// bFrom := brcCall(_bFrom)
	// bTo := brcCall(_bTo)

	// balanceFrom := isNumber(bFrom)
	// balanceTo := isNumber(bTo)

	logging.Info("blockNumber:", blockNum, ",approveFrom:", _from, ",approveTo:", _to,
		",approveToken:", txNum)
}

func (t *token) transferToken(txNum *big.Int, topic []interface{}, blockNum uint64) {
	_from := "0x" + topic[1].(string)[26:]
	_to := "0x" + topic[2].(string)[26:]

	_bFrom := _balance + _from[2:]
	_bTo := _balance + _to[2:]

	bFrom := brcCall(_bFrom)
	bTo := brcCall(_bTo)

	balanceFrom := isNumber(bFrom)
	balanceTo := isNumber(bTo)

	logging.Info("blockNumber:", blockNum, ",transferFrom:", _from,
		",to:", _to, ",TxNum:", txNum, ",balanceFrom:", balanceFrom, ",balanceTo:", balanceTo)
}

func (t *token) increaseToken(txNum *big.Int, topic []interface{}, blockNum uint64) {
	_operateAddr := "0x" + topic[1].(string)[26:]
	_bAddr := _balance + _operateAddr[2:]

	_total := brcCall(totalSupply)
	_baAddr := brcCall(_bAddr)

	_totalSupply := isNumber(_total)
	_balanceAddr := isNumber(_baAddr)

	logging.Info("blockNumber:", blockNum, ",increaseTokenAddr:", _operateAddr,
		",increaseToken:", txNum, ",nowTotalSupply:", _totalSupply,
		",increaseTokenAddrBalance:", _balanceAddr)
}

func (t *token) burnToken(txNum *big.Int, topic []interface{}, blockNum uint64) {
	_from := "0x" + topic[1].(string)[26:]
	_bFrom := _balance + _from[2:]

	_total := brcCall(totalSupply)
	_baFrom := brcCall(_bFrom)

	_totalSupply := isNumber(_total)
	_balanceFrom := isNumber(_baFrom)

	logging.Info("blockNumber:", blockNum, ",burnTokenAddr:", _from, ",burnToken:", txNum,
		",nowTotalSupply:", _totalSupply, ",burnTokenAddrBalance:", _balanceFrom)
}

func (t *token) changeContractOwner(topic []interface{}, blockNum uint64) {
	oldOwner := "0x" + topic[1].(string)[26:]
	newOwner := "0x" + topic[2].(string)[26:]

	logging.Info("blockNumber:", blockNum, ",oldOwner:", oldOwner, "newOwner:", newOwner)
}

func getTxNum(txStr string) *big.Int {
	str := txStr[2:]
	num := removeZero(str)

	n := new(big.Int)
	n, _ = n.SetString(num, 16)

	return n
}

func removeZero(str string) string {
	a := rune(48)
	var n int
	for _, v := range str {
		if v == a {
			n++
		} else {
			break
		}
	}
	str = str[n:]

	return str
}

func postMsg(msg, ip string) []byte {
	cli := &http.Client{
		Timeout: 200 * time.Second,
	}

	body := bytes.NewBuffer([]byte(msg))
	resq, err := http.NewRequest("POST", ip, body)
	if err != nil {
		fmt.Println("resq faild")
	}
	resq.Close = true

	rep, err := cli.Do(resq)
	if err != nil {
		fmt.Println("get response faild!")
	}

	defer rep.Body.Close()
	rcv, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		fmt.Println("get body faild")
	}

	return rcv
}
