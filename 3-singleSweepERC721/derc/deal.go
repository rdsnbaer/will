package derc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/my/repo/4-other/3-ERC721/loginfo"
)

func IsAddress(res string) (string, error) {
	lens := len(res)
	switch lens {
	case 66:
		return ("0x" + res[26:]), nil
	case 64:
		return ("0x" + res[24:]), nil
	case 40:
		return ("0x" + res), nil
	case 42:
		return res, nil
	default:
		err := errors.New("UnKnow addr format")
		return res, err
	}
}

func IsBool(res string) (bool, error) {
	lens := len(res)
	fmt.Println("len", lens)

	switch lens {
	case 66:
		str := res[65:]
		return isBool(str)
	case 64:
		str := res[63:]
		return isBool(str)
	case 1:
		str := res
		return isBool(str)
	default:
		err := errors.New("UnKnow bool format")
		return false, err
	}
}

func isBool(str string) (bool, error) {
	if str == "1" {
		return true, nil
	}
	if str == "0" {
		return false, nil
	}
	err := errors.New("UnKnow bool format")
	return false, err
}

func IsNumber(res string) *big.Int {
	n := RemoveZero(res[2:])
	num := new(big.Int)

	num, _ = num.SetString(n, 16)

	return num
}

func IsArray(res string) []int64 {
	var ownerTokenID []int64
	arr := res[66:]
	r := arr[:64]
	ln := RemoveZero(r)

	len, _ := strconv.ParseInt(ln, 16, 64)
	s := r[64:]
	var i int64
	for i = 0; i < len; i++ {
		k := (i + 1) * 64
		n := i * 64
		id := s[n:k]

		_tokenId := RemoveZero(id)
		tokenId, _ := strconv.ParseInt(_tokenId, 16, 64)
		ownerTokenID = append(ownerTokenID, tokenId)
	}

	return ownerTokenID
}

func IsString(res string) (string, error) {
	loginfo.Info("5.1-result:", res)
	fmt.Println("5.1-result:", res)

	r := res[66:]
	fmt.Println("1ok")
	s := r[:64]
	k := RemoveZero(s)

	len, err := strconv.ParseUint(k, 16, 64)
	if err != nil {
		return "", err
	}

	index := len*2 + 64
	str := r[64:index]
	src := []byte(str)

	return Decode(src)

}

func Decode(src []byte) (string, error) {
	dstDecode := make([]byte, hex.DecodedLen(len(src)))

	_, err := hex.Decode(dstDecode, src)
	if err != nil {
		return "", err
	}

	decodeStr, err := hex.DecodeString(string(src))
	if err != nil {
		return "", err
	}

	str := string(decodeStr)
	return str, nil
}

func RemoveZero(res string) string {
	var index int
	zero := rune(48)

	for n, v := range res {
		if v != zero {
			index = n
			break
		}
	}
	if index == len(res) {
		return "0"
	}
	return res[index:]
}

func ComplementNumber(num string) string {
	n, _ := strconv.Atoi(num)
	r := int64(n)

	str := ToHex(r, false)
	len := len(str)

	zeroNum := 64 - len
	for i := 0; i < zeroNum; i++ {
		str = "0" + str
	}
	return str
}

func ToHex(num int64, prefix bool) string {
	if prefix {
		return "0x" + strconv.FormatInt(num, 16)
	} else {
		return strconv.FormatInt(num, 16)
	}
}

func PostMsg(recv []byte, ip string) ([]byte, error) {
	cli := &http.Client{Timeout: 20 * time.Second}

	body := bytes.NewBuffer(recv)
	resq, err := http.NewRequest("POST", ip, body)
	if err != nil {
		return nil, err
	}
	resq.Close = true

	rep, err := cli.Do(resq)
	if err != nil {
		return nil, err
	}
	defer rep.Body.Close()

	res, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return res, err
	}

	return res, nil
}

func GetUnmarshal(recv []byte) (map[string]interface{}, error) {
	var jsonMap map[string]interface{}

	err := json.Unmarshal(recv, &jsonMap)
	if err != nil {
		return nil, err
	}

	return jsonMap, nil
}
