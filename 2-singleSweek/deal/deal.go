package deal

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"

	"github.com/my/repo/2-saokuai/logging"
)

func IsAddress(recv []byte) string {
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

func IsString(recv []byte) string {
	var jsonMap map[string]interface{}
	err := json.Unmarshal(recv, &jsonMap)
	if err != nil {
		logging.Error(err)
	}

	result := jsonMap["result"].(interface{})
	rString := result.(string)
	res := rString[66:]

	ln := RemoveZero(res[:64])
	len, _ := strconv.ParseUint(ln, 16, 64)
	index := len*2 + 64
	str := res[64:index]

	src := []byte(str)
	return Decode(src)
}

func IsNumber(recv []byte) (string, *big.Int) {
	var jsonMap map[string]interface{}
	err := json.Unmarshal(recv, &jsonMap)
	if err != nil {
		log.Fatal(err)
	}

	result := jsonMap["result"].(interface{})
	res := result.(string)
	s := GetTxNum(res)
	supply := fmt.Sprintf("%d", s)

	return supply, s
}

func GetTxNum(txStr string) *big.Int {
	str := txStr[2:]
	num := RemoveZero(str)

	n := new(big.Int)
	n, _ = n.SetString(num, 16)

	return n
}

func RemoveZero(str string) string {
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

func Decode(src []byte) string {
	dstDeCode := make([]byte, hex.DecodedLen(len(src)))

	_, err := hex.Decode(dstDeCode, src)
	if err != nil {
		log.Fatal(err)
	}

	decodedStr, _ := hex.DecodeString(string(src))
	str := string(decodedStr)

	return str
}

func ToHex(num int64, prefix bool) string {
	if prefix {
		return "0x" + strconv.FormatInt(num, 16)
	} else {
		return strconv.FormatInt(num, 16)
	}
}
