package deal

import (
	"encoding/json"

	"github.com/my/repo/2-saokuai/http"
)

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
