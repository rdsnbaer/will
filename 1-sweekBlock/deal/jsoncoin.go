package deal

import (
	"encoding/json"
	"errors"
)

const (
	ZERO = "0"
)

type CoinInfo struct {
	Precision           int64
	Total, Name, Symbol string
}

func AnalysisJsonCoin(recv map[string]interface{}) (CoinInfo, error) {
	var coinInfo CoinInfo

	pre := recv["precision"].(float64)
	prec := int64(pre)
	tol := recv["total"].(string)
	var i int64
	for i = 0; i < prec; i++ {
		tol = tol + "0"
	}

	coinInfo.Precision = prec
	coinInfo.Total = tol
	coinInfo.Precision = prec
	coinInfo.Total = tol
	coinInfo.Name = recv["name"].(string)
	coinInfo.Symbol = recv["symbol"].(string)
	deploySuccess := recv["deploySuccess"].(string)

	if deploySuccess == ZERO {
		return coinInfo, errors.New("deploySuccess is ZERO:" + coinInfo.Symbol)
	}

	return coinInfo, nil
}

func JsonCoinsInfo(recv []byte) ([]interface{}, error) {
	var jsonMap map[string]interface{}

	err := json.Unmarshal(recv, &jsonMap)
	if err != nil {
		return nil, err
	}

	jsonCoins := jsonMap["coins"].([]interface{})

	return jsonCoins, nil
}

func JsonCoinInfo(recv []byte) (map[string]interface{}, error) {
	var jsonMap map[string]interface{}

	err := json.Unmarshal(recv, &jsonMap)
	if err != nil {
		return nil, err
	}

	if jsonMap["coin"] == nil {
		return nil, errors.New("coin is nil")
	}

	jsonCoin := jsonMap["coin"].(map[string]interface{})

	return jsonCoin, nil
}

func GetSymbolInfo(symbolInfo map[string]interface{}) interface{} {
	hash := symbolInfo["txHash"].(interface{})

	return hash
}

func JsonMaps(recv []byte) (map[string]interface{}, error) {
	var jsonMap map[string]interface{}

	err := json.Unmarshal(recv, &jsonMap)
	if err != nil {
		return nil, err
	}

	jsonResult := jsonMap["result"].(map[string]interface{})

	return jsonResult, nil
}
