package derc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/my/repo/4-other/3-ERC721/loginfo"
	"math/big"
)

const (
	FixedSt          = "00000000000000000000000000000000000000000000000000000000"
	FixedStr         = "000000000000000000000000"
	ContractOwn      = "0x5e51ae6f"
	IsPaused         = "0xb187bd26"
	Name             = "0x06fdde03"
	Symbol           = "0x95d89b41"
	TotalSupply      = "0x18160ddd"
	TheNewestTokenId = "0x4d156d0f"
)

var (
	BalanceOf           = "0x70a08231" + FixedStr // + addr
	GetApproved         = "0x081812fc"            // + tokenId
	IsApprovedForAll    = "0xe985e9c5" + FixedStr // + addr + fixedStr + addr
	OwnerOf             = "0x6352211e"            // + tokenId
	SupportsInterface   = "0x01ffc9a7"            // + interfaceId(向右补齐)
	TokenByIndex        = "0x4f6ccce7"            // + index
	TokenOfOwnerByIndex = "0x2f745c59" + FixedStr // + addr + index
	TokenURI            = "0xc87b56dd"            // + TokenId
	TokensOfOwner       = "0x8462151c" + FixedStr // + addr
)

type ContractInterface interface {
	GetContractOwner(from, to, blockNum, ip interface{}) (string, error)
	GetContractIsPaused(from, to, blockNum, ip interface{}) (bool, error)
	GetNFTName(from, to, blockNum, ip interface{}) (string, error)
	GetNFTSymbol(from, to, blockNum, ip interface{}) (string, error)
	GetTotalSupply(from, to, blockNum, ip interface{}) (*big.Int, error)
	GetNewestTokenId(from, to, blockNum, ip interface{}) (*big.Int, error)

	GetAddrBalance(from, to, blockNum, ip interface{}) (*big.Int, error)
	GetTokenIdApproved(from, to, tokenID, blockNum, ip interface{}) (string, error)
	GetIsApprovedForAll(from, to, owner, operator, blockNum, ip interface{}) (bool, error)
	GetTokenIdOwnerOf(from, to, tokenID, blockNum, ip interface{}) (string, error)
	GetSupportsInterface(from, to, interfaceID, blockNum, ip interface{}) (bool, error)
	GetTokenIndex(from, to, index, blockNum, ip interface{}) (*big.Int, error)
	GetTokenOfOwnerByIndex(from, to, owner, index, blockNum, ip interface{}) (*big.Int, error)
	GetTokenURI(from, to, tokenID, blockNum, ip interface{}) (string, error)
	GetTokensOfOwner(from, to, owner, blockNum, ip interface{}) ([]int64, error)
}

type Query struct{}

type Contract struct {
	Name_C, Symbol_C, ContractOwner_C string
}

func BrcCall(from, to, data string) ([]byte, error) {
	return BrcCaller(from, to, data, "latest")
}

func BrcCaller(from, to, data, blockNum string) ([]byte, error) {
	jsonMap := make(map[string]interface{})
	jsonParams := make(map[string]interface{})

	jsonParams["from"] = from
	jsonParams["to"] = to
	jsonParams["data"] = data

	jsonMap["jsonrpc"] = "2.0"
	jsonMap["method"] = "brc_call"
	jsonMap["params"] = []interface{}{jsonParams, blockNum}
	jsonMap["id"] = 1

	recv, err := json.Marshal(jsonMap)
	if err != nil {
		return recv, err
	}

	return recv, nil
}

func (q *Query) GetContractOwner(from, to, blockNum, ip string) (string, error) {
	result, err := getBrcCallInfo(from, to, ContractOwn, blockNum, ip)
	if err != nil {
		return "", err
	}

	return IsAddress(result)
}

func (q *Query) GetContractIsPaused(from, to, blockNum, ip string) (bool, error) {
	result, err := getBrcCallInfo(from, to, IsPaused, blockNum, ip)
	if err != nil {
		return false, err
	}

	return IsBool(result)
}

func (q *Query) GetNFTName(from, to, blockNum, ip string) (string, error) {
	result, err := getBrcCallInfo(from, to, Name, blockNum, ip)
	if err != nil {
		return result, err
	}

	return IsString(result)
}

func (q *Query) GetNFTSymbol(from, to, blockNum, ip string) (string, error) {
	result, err := getBrcCallInfo(from, to, Symbol, blockNum, ip)
	if err != nil {
		return result, err
	}

	return IsString(result)
}

func (q *Query) GetTotalSupply(from, to, blockNum, ip string) (*big.Int, error) {
	result, err := getBrcCallInfo(from, to, TotalSupply, blockNum, ip)
	if err != nil {
		return big.NewInt(0), err
	}

	return IsNumber(result), nil
}

func (q *Query) GetNewestTokenId(from, to, blockNum, ip string) (*big.Int, error) {
	result, err := getBrcCallInfo(from, to, TheNewestTokenId, blockNum, ip)
	if err != nil {
		return big.NewInt(0), err
	}

	return IsNumber(result), nil
}

func (q *Query) GetAddrBalance(from, to, blockNum, ip string) (*big.Int, error) {
	var data string
	lens := len(from)
	if lens == 42 {
		data = BalanceOf + from[2:]
	}
	if lens == 40 {
		data = BalanceOf + from
	}

	result, err := getBrcCallInfo(from, to, data, blockNum, ip)
	if err != nil {
		return big.NewInt(0), err
	}

	return IsNumber(result), nil
}

func (q *Query) GetTokenIdApproved(from, to, tokenID, blockNum, ip string) (string, error) {
	num := ComplementNumber(tokenID)
	data := GetApproved + num

	result, err := getBrcCallInfo(from, to, data, blockNum, ip)
	if err != nil {
		return "", err
	}

	return IsString(result)
}

func (q *Query) GetIsApprovedForAll(from, to, owner, operator, blockNum, ip string) (bool, error) {
	if len(owner) == 42 {
		owner = owner[2:]
	}
	if len(operator) == 42 {
		operator = operator[2:]
	}

	data := IsApprovedForAll + owner + FixedStr + operator

	result, err := getBrcCallInfo(from, to, data, blockNum, ip)
	if err != nil {
		return false, err
	}

	return IsBool(result)
}

func (q *Query) GetTokenIdOwnerOf(from, to, tokenID, blockNum, ip string) (string, error) {
	num := ComplementNumber(tokenID)
	data := OwnerOf + num

	result, err := getBrcCallInfo(from, to, data, blockNum, ip)
	if err != nil {
		return "", err
	}

	return IsAddress(result)
}

func (q *Query) GetSupportsInterface(from, to, interfaceID, blockNum, ip string) (bool, error) {
	data := interfaceID + FixedSt

	result, err := getBrcCallInfo(from, to, data, blockNum, ip)
	if err != nil {
		return false, err
	}

	return IsBool(result)
}

func (q *Query) GetTokenIndex(from, to, index, blockNum, ip string) (*big.Int, error) {
	num := ComplementNumber(index)
	data := TokenByIndex + num

	result, err := getBrcCallInfo(from, to, data, blockNum, ip)
	if err != nil {
		return big.NewInt(0), err
	}

	return IsNumber(result), nil
}

func (q *Query) GetTokenOfOwnerByIndex(from, to, owner, index, blockNum, ip string) (*big.Int, error) {
	num := ComplementNumber(index)
	if len(owner) == 42 {
		owner = owner[2:]
	}
	data := TokenOfOwnerByIndex + owner + num

	result, err := getBrcCallInfo(from, to, data, blockNum, ip)
	if err != nil {
		return big.NewInt(0), err
	}

	return IsNumber(result), nil
}

func (q *Query) GetTokenURI(from, to, tokenID, blockNum, ip string) (string, error) {
	num := ComplementNumber(tokenID)
	data := TokenOfOwnerByIndex + num

	result, err := getBrcCallInfo(from, to, data, blockNum, ip)
	if err != nil {
		return "", err
	}

	return IsString(result)
}

func (q *Query) GetTokensOfOwner(from, to, owner, blockNum, ip string) ([]int64, error) {
	if len(owner) == 42 {
		owner = owner[2:]
	}
	data := TokensOfOwner + owner

	result, err := getBrcCallInfo(from, to, data, blockNum, ip)
	if err != nil {
		return nil, err
	}

	return IsArray(result), nil
}

func getBrcCallInfo(from, to, data, blockNum, ip string) (string, error) {
	call, err := BrcCaller(from, to, data, blockNum)
	if err != nil {
		fmt.Println("3.1-getBrcCallInfo->BrcCaller err:", err)
		loginfo.Error("3.1-getBrcCallInfo->BrcCaller err:", err)
		return "", err
	}

	recv, err := PostMsg(call, ip)
	if err != nil {
		fmt.Println("3.2-getBrcCallInfo->PostMsg err:", err)
		loginfo.Info("3.2-getBrcCallInfo->PostMsg err:", err)
		return "", err
	}

	jsonMap, err := GetUnmarshal(recv)
	if err != nil {
		fmt.Println("3.3-getBrcCallInfo->GetUnmarshal err:", err)
		loginfo.Info("3.3-getBrcCallInfo->GetUnmarshal err:", err)
		return "", err
	}

	if jsonMap["result"] == nil {
		err = errors.New("result is nil")
		fmt.Println("3.4-getBrcCallInfo err:", err)
		loginfo.Error("3.4-getBrcCallInfo err:", err)
		return "", err
	}

	return jsonMap["result"].(string), nil
}
