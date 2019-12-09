package cmdWallet

import (
	"encoding/json"
	"fmt"
)

type Data struct {
	TypeName int    `json:"type"`
	Contract string `json:"contract"`
}

type Source struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Value    string `json:"value"`
	Data     []Data `json:"data"`
	ChainId  string `json:"chainId"`
	Nonce    string `json:"nonce"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
}

type Contract struct {
	Sources []Source `json:"source"`
	Key     []string `json:"keys"`
}

func WriteJson([]byte) {
	//	file, _ := os.OpenFile("./Json/1-exContract.json", os.O_TRUNC|os.O_CREATE|os.O_WRONLY)
}

func ExContract(from, to, data, ks string) ([]byte, error) {
	var contract Contract
	var _s Source
	var _d Data

	_k := []string{ks}
	contract.Key = _k

	_d.TypeName = 6
	_d.Contract = data

	_s.Data = []Data{_d}
	_s.From = from
	_s.To = to
	_s.Value = "0x0"
	_s.ChainId = "0x1"
	_s.Nonce = "0x0"
	_s.Gas = "0x1388000"
	_s.GasPrice = "0x5"

	contract.Sources = []Source{_s}

	mJson, err := json.Marshal(contract)
	if err != nil {
		fmt.Println("1-mJson is err:", err)
		return mJson, err
	}

	return mJson, nil
}
