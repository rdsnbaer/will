package rpc

import (
	"encoding/json"

	"github.com/my/repo/2-saokuai/deal"
	"github.com/my/repo/2-saokuai/http"
)

func BrcGetBlockDetialByNumber(num int64, ip string) []byte {
	msg := make(map[string]interface{})

	msg["jsonrpc"] = "2.0"
	msg["method"] = "brc_getBlockDetialByNumber"
	msg["params"] = []interface{}{deal.ToHex(num, true), true}
	msg["id"] = 1

	mjson, _ := json.Marshal(msg)
	mString := string(mjson)
	recv := http.PostMsg(mString, ip)

	return recv
}

func BrcGetTransactionReceipt(h interface{}, _ip string) []byte {
	msg := make(map[string]interface{})

	msg["jsonrpc"] = "2.0"
	msg["method"] = "brc_getTransactionReceipt"
	msg["params"] = []interface{}{h.(string)}
	msg["id"] = 1

	mjson, _ := json.Marshal(msg)
	mString := string(mjson)
	recv := http.PostMsg(mString, _ip)

	return recv
}
