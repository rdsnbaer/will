package derc

import (
	"fmt"

	"github.com/my/repo/4-other/3-ERC721/loginfo"
)

const (
	TransferMintBurnTopic = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	ApproveTopic          = "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925"
	SetApprovalTopic      = "0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31"
	ChangeOwnerTopic      = "0x07ff91ee303faa39fb867eb01411e233f560ef81a7eef1aff036a194236d1161"
	FixedTopic            = "0x0000000000000000000000000000000000000000000000000000000000000000"
)

type MintToken struct {
	From, To, TokenId, ContractAddr, Hash string
}

type comContract struct {
	objContract, actualContract, hash string
}
type Transfers struct {
	from, to, tokenId, contractAddr, hash string
}

type ExecInterface interface {
	GetMintInfo(res, topic, contractAddr, ip interface{})
	GetApproveInfo(res, topic, contractAddr, ip interface{})
	GetTransferFromInfo(res, topic, contractAddr, ip interface{})
	GetBurnInfo(topic []interface{}, contractAddr, ip interface{})
	GetSetApprovalForAllInfo(res, log, topic, contractAddr, ip interface{})
	GetChangeContractOwnerInfo(res, topic, contractAddr, ip interface{})
}

type Exec struct{}

func DealResultAndTopic(res, log map[string]interface{}, topic []interface{}, contractAddr, ip string) {
	len := len(topic)
	hash := res["transactionHash"].(string)
	switch len {
	case 3:
		DealThreeFiled(res, log, topic, contractAddr, ip)
	case 4:
		DealTransferInfo(res, topic, contractAddr, ip)
	default:
		fmt.Println("6.1-topic is not ERC721, hash:", hash)
		loginfo.Error("6.1-topic is not ERC721, hash:", hash)
	}
}

func DealThreeFiled(res, log map[string]interface{}, topic []interface{}, contractAddr, ip string) {
	topic_0 := topic[0].(string)
	var e Exec
	hash := res["transactionHash"].(string)

	switch topic_0 {
	case SetApprovalTopic:
		e.GetSetApprovalForAllInfo(res, log, topic, contractAddr, ip)
	case ChangeOwnerTopic:
		e.GetChangeContractOwnerInfo(res, topic, contractAddr, ip)
	default:
		fmt.Println("8.1-not ERC721 topic, hash:", hash)
		loginfo.Error("8.1-not ERC721 topic, hash:", hash)
	}
}

func DealTransferInfo(res map[string]interface{}, topic []interface{}, contractAddr, ip string) {
	topic_0 := topic[0].(string)
	hash := res["transactionHash"].(string)

	if topic_0 == TransferMintBurnTopic {
		var e Exec
		topic_1 := topic[1].(string)
		topic_2 := topic[2].(string)

		if (topic_1 == FixedTopic) || (topic_2 == FixedTopic) {
			if (topic_1 == FixedTopic) && (topic_2 == FixedTopic) {
				fmt.Println("6.4-topic[1] or topic[2] is the same, hash:", hash)
				loginfo.Info("6.4-topic[1] or topic[2] is the same, hash:", hash)
			} else {
				if (topic_1 == FixedTopic) && (topic_2 != FixedTopic) {
					e.GetMintInfo(res, topic, contractAddr, ip)
				}
				if (topic_1 != FixedTopic) && (topic_2 == FixedTopic) {
					e.GetBurnInfo(res, topic, contractAddr, ip)
				}
			}
		} else {
			e.GetTransferFromInfo(res, topic, contractAddr, ip)
		}
	} else {
		fmt.Println("6.2-topic[0] is not TransferMintBurnTopic, hash:", hash)
		loginfo.Error("6.2-topic[0] is not TransferMintBurnTopic, hash:", hash)
	}

}

func (e *Exec) GetMintInfo(res map[string]interface{}, topic []interface{}, contractAddr, ip string) {
	var mint MintToken

	toStr := topic[2].(string)
	id := topic[3].(string)

	to, _ := IsAddress(toStr)
	_id := IsNumber(id)
	tokenID := fmt.Sprintf("%v", _id)
	from := res["from"].(string)

	mint.From = from
	mint.To = to
	mint.TokenId = tokenID
	mint.Hash = res["transactionHash"].(string)
	mint.ContractAddr = res["to"].(string)

	loginfo.Info("7.1-[Mint] from:", from, ", to:", to, ",TokenId:",
		tokenID, ", contractAddr:", contractAddr, ",hash:", mint.Hash)
	loginfo.Info("-------------------------------------------------------------------")

	fmt.Println("7.1-[Mint] from:", from, ", to:", to, ", TokenId:",
		tokenID, ", contractAddr:", contractAddr, ", hash:", mint.Hash)
	fmt.Println("7.1-mint:", mint)
	fmt.Println()
}

func (e *Exec) GetBurnInfo(res map[string]interface{}, topic []interface{}, contractAddr, ip string) {

	conAddr := res["to"].(string)
	hash := res["transactionHash"].(string)

	if conAddr == contractAddr {
		fromStr := topic[1].(string)
		from, _ := IsAddress(fromStr)

		id := topic[3].(string)
		_tokenID := IsNumber(id)
		tokenID := fmt.Sprintf("%v", _tokenID)

		operator := res["from"].(string)
		if from == operator {
			fmt.Println("7.3-Burn, from:", from, ", tokenID:", tokenID, ", operator:", operator)
			loginfo.Info("7.3-Burn, from:", from, ", tokenID:", tokenID, ", operator:", operator)
			loginfo.Info("-------------------------------------------------------------------")
		} else {
			fmt.Println("7.4-Burn err, operator:", operator, ",from: ", from, ", tokenID:", tokenID,
				", hash:", hash)
			loginfo.Error("7.4-Burn err, operator:", operator, ",from: ", from, ", tokenID:", tokenID,
				", hash:", hash)
		}
	} else {
		fmt.Println("7.2-contractAddr err, objContract:", contractAddr, ", actualContract:", conAddr,
			", hash:", hash)
		loginfo.Error("7.2-contractAddr err, objContract:", contractAddr, ", actualContract:", conAddr,
			", hash:", hash)
	}
}

func (e *Exec) GetTransferFromInfo(res map[string]interface{}, topic []interface{}, contractAddr, ip string) {

	addr := res["to"].(string)
	if addr == contractAddr {
		hash := res["transactionHash"].(string)
		fromStr := topic[1].(string)
		toStr := topic[2].(string)
		id := topic[3].(string)

		from, _ := IsAddress(fromStr)
		to, _ := IsAddress(toStr)
		iD := IsNumber(id)
		tokenID := fmt.Sprintf("%v", iD)

		fmt.Println("7.5-TransferFrom, from:", from, ", to:", to, ", tokenId:", tokenID, ", hash:", hash)
		loginfo.Info("7.5-TransferFrom, from:", from, ", to:", to, ", tokenId:", tokenID, ", hash:", hash)
		loginfo.Info("-------------------------------------------------------------------")
		fmt.Println()
	} else {
		fmt.Println("not the same contractAddr, objContract:", contractAddr, ", actualContract:", addr)
		loginfo.Error("not the same contractAddr, objContract:", contractAddr, ", actualContract:", addr)
	}

}

func (e *Exec) GetSetApprovalForAllInfo(res, log map[string]interface{}, topic []interface{}, contractAddr, ip interface{}) {
	fmt.Println("7.10-topic:", topic)
	from, _ := IsAddress(topic[1].(string))
	to, _ := IsAddress(topic[2].(string))
	operator := res["from"].(string)
	hash := res["transactionHash"].(string)

	if from == operator {
		a := log["data"].(string)
		approval, _ := IsBool(a)
		fmt.Println("7.11-from:", from, ",to:", to, ", approval:", approval, ", hash:", hash)
		loginfo.Info("7.11-from:", from, ",to:", to, ", approval:", approval, ", hash:", hash)
		fmt.Println()
		loginfo.Info("-------------------------------------------------------------------")
	} else {
		fmt.Println("7.10-operator and from not same, operator:", operator, ", from:", from, ", hash:", hash)
		loginfo.Info("7.10-operator and from not same, operator:", operator, ", from:", from, ", hash:", hash)
	}

}

func (e *Exec) GetChangeContractOwnerInfo(res map[string]interface{}, topic []interface{}, contractAddr, ip interface{}) {
	from, _ := IsAddress(topic[1].(string))
	to, _ := IsAddress(topic[2].(string))
	operator := res["from"].(string)
	hash := res["transactionHash"].(string)
	addr := res["to"].(string)

	if from == operator && addr == contractAddr {
		num := res["blockNumber"].(float64)
		blockNum := int64(num)
		fmt.Println("7.13-contract owner, from:", from, ", to:", to, ", contractAddr:", addr,
			",hash:", hash, ", blockNumber", blockNum)
		loginfo.Info("7.13-contract owner, from:", from, ", to:", to, ", contractAddr:", addr,
			",hash:", hash, ", blockNumber", blockNum)
		fmt.Println()
		loginfo.Info("-------------------------------------------------------------------")
	} else {
		fmt.Println("7.12-the operator or contractAddr err, operator:", operator, ",from:", from,
			", objContract:", contractAddr, ",actualContract:", addr, ",hash:", hash)
		loginfo.Info("7.12-the operator or contractAddr err, operator:", operator, ",from:", from,
			", objContract:", contractAddr, ", to:", to, ",actualContract:", addr, ",hash:", hash)
	}
}

func (e *Exec) GetApproveInfo(topic []interface{}, contractAddr, ip interface{}) {

}
