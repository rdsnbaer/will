package rpc

import (
	"fmt"
	"log"
	"math/big"

	"github.com/my/repo/2-saokuai/deal"
	"github.com/my/repo/2-saokuai/logging"
)

type tokenIce interface {
	transferToken(txNum, _totals, _sum interface{}, topic []interface{}, blockNum interface{}, conAddr, _ip interface{})
	increaseToken(txNum, _totals, _sum interface{}, topic []interface{}, blockNum interface{}, conAddr, _ip interface{})
	approveToken(txNum, _totals, _sum interface{}, topic []interface{}, blockNum interface{}, conAddr, _ip interface{})
	burnToken(txNum interface{}, topic []interface{}, blockNum interface{}, conAddr, _ip interface{})
	changeContractOwner(topic []interface{}, blockNum interface{}, conAddr, _ip interface{})
}

type token struct {
}

func (t *token) approveToken(txNum *big.Int, topic []interface{}, blockNum uint64, conAddr, _ip string) {
	_from := "0x" + topic[1].(string)[26:]
	_to := "0x" + topic[2].(string)[26:]

	logging.Info("blockNumber:", blockNum, ",approveFrom:", _from, ",approveTo:", _to,
		",approveToken:", txNum, ",contractAddr:", conAddr, ",ip:", _ip)
}

func (t *token) transferToken(txNum, _totals, _sum *big.Int, topic []interface{},
	blockNum uint64, conAddr, _ip string) {

	_from := "0x" + topic[1].(string)[26:]
	_to := "0x" + topic[2].(string)[26:]

	_bFrom := _balance + _from[2:]
	_bTo := _balance + _to[2:]
	_blockNum := fmt.Sprintf("%d", blockNum)
	_Preb := fmt.Sprintf("%d", blockNum-1)

	bFrom := BrcCaller(_from, conAddr, _bFrom, _ip, _blockNum)
	bTo := BrcCaller(_to, conAddr, _bTo, _ip, _blockNum)
	_bf := BrcCaller(_from, conAddr, _bFrom, _ip, _Preb)
	_bt := BrcCaller(_to, conAddr, _bTo, _ip, _Preb)

	balanceFrom, _fromB := deal.IsNumber(bFrom)
	balanceTo, _toB := deal.IsNumber(bTo)
	_, preFrom := deal.IsNumber(_bf)
	_, preTo := deal.IsNumber(_bt)

	var subFrom, subTo *big.Int
	subFrom = subFrom.Sub(preFrom, _fromB)
	subTo = subTo.Sub(_toB, preTo)

	val := subTo.Cmp(subFrom)
	vas := txNum.Cmp(subFrom)
	vab := txNum.Cmp(subTo)
	if val == 0 && vas == 0 && vab == 0 {
		logging.Info("blockNumber:", blockNum, ",transferFrom:", _from,
			",to:", _to, ",TxNum:", txNum, ",balanceFrom:", balanceFrom,
			",balanceTo:", balanceTo, ",contractAddr:", conAddr)
	} else {

		_sum = _sum.Add(_sum, subTo)
		_sum = _sum.Sub(_sum, subFrom)
		logging.Error("blockNumber:", blockNum, ",transferFrom:", _from,
			",to:", _to, ",TxNum:", txNum, ",contractAddr:", conAddr, ",_sum:", _sum)
		log.Fatal("transfer error")
	}

}

func (t *token) increaseToken(txNum, _totals, _sum *big.Int, topic []interface{},
	blockNum uint64, conAddr, _ip string) {

	_operateAddr := "0x" + topic[1].(string)[26:]
	_bAddr := _balance + _operateAddr[2:]
	_blockNum := fmt.Sprintf("%d", blockNum)

	_total := BrcCaller(_operateAddr, conAddr, totalSupply, _ip, _blockNum)
	_baAddr := BrcCaller(_operateAddr, conAddr, _bAddr, _ip, _blockNum)

	_totalSupply, _t := deal.IsNumber(_total)
	_balanceAddr, _ := deal.IsNumber(_baAddr)

	_totals = _totals.Add(_totals, _t)
	_sum = _sum.Add(_sum, _t)

	logging.Info("blockNumber:", blockNum, ",increaseTokenAddr:", _operateAddr,
		",increaseToken:", txNum, ",nowTotalSupply:", _totalSupply,
		",increaseTokenAddrBalance:", _balanceAddr, ",contractAddr:", conAddr)
}

func (t *token) burnToken(txNum, _totals, _sum *big.Int, topic []interface{}, blockNum uint64,
	conAddr, _ip string) {

	_from := "0x" + topic[1].(string)[26:]
	_bFrom := _balance + _from[2:]
	_blockNum := fmt.Sprintf("%d", blockNum)

	_total := BrcCaller(_from, conAddr, totalSupply, _ip, _blockNum)
	_baFrom := BrcCaller(_from, conAddr, _bFrom, _ip, _blockNum)

	_totalSupply, _t := deal.IsNumber(_total)
	_balanceFrom, _ := deal.IsNumber(_baFrom)

	_totals = _totals.Sub(_totals, _t)
	_sum = _sum.Sub(_sum, _t)

	logging.Info("blockNumber:", blockNum, ",burnTokenAddr:", _from, ",burnToken:", txNum,
		",nowTotalSupply:", _totalSupply, ",burnTokenAddrBalance:", _balanceFrom,
		",contractAddr:", conAddr)
}

func (t *token) changeContractOwner(topic []interface{}, blockNum uint64, conAddr, _ip string) {
	oldOwner := "0x" + topic[1].(string)[26:]
	newOwner := "0x" + topic[2].(string)[26:]

	logging.Info("blockNumber:", blockNum, ",oldOwner:", oldOwner, "newOwner:", newOwner)
}
