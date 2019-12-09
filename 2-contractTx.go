package main

import (
	"github.com/my/repo/2-saokuai/cmdWallet"
)

var (
	from     = "e79fead329b69540142fabd881099c04424cc49f"
	to       = "0x9c39af839ab7bf4c3bc37acbf7086a720ce38df8"
	contract = "0x3ead67b50000000000000000000000000f63870e5062f01b87bfc804f413e251982a4764"
	key      = "6vsDoEvshzWJUhY5Aa8xd9WtHwNs9tjmeFq8NnJBST6t"
)

func main() {
	recv, err := cmdWallet.ExContract(from, to, contract, key)
	if err != nil {
		return
	}

	cmdWallet.WriteJson(recv)

}

// // 交易发起方keystore文件地址
// var fromKeyStoreFile = "";

// // keystore文件对应的密码
// var password = "";

// // 交易接收方地址
// var toAddress = ""

// // http服务地址, 例:http://localhost:8545
// var httpUrl = "http://ip:port"

// /*
// 	以太坊交易发送
//  */
// func TestSendTx(t *testing.T){
// 	// 交易发送方
// 	// 获取私钥方式一，通过keystore文件
// 	fromKeystore,err := ioutil.ReadFile(fromKeyStoreFile)
// 	require.NoError(t,err)
// 	fromKey,err := keystore.DecryptKey(fromKeystore,password)
// 	fromPrivkey := fromKey.PrivateKey
// 	fromPubkey := fromPrivkey.PublicKey
// 	fromAddr := crypto.PubkeyToAddress(fromPubkey)

// 	// 获取私钥方式二，通过私钥字符串
// 	//privateKey, err := crypto.HexToECDSA("私钥字符串")

// 	// 交易接收方
// 	toAddr := common.StringToAddress(toAddress)

// 	// 数量
// 	amount := big.NewInt(14)

// 	// gasLimit
// 	var gasLimit uint64 = 300000

// 	// gasPrice
// 	var gasPrice *big.Int = big.NewInt(200)

// 	// 创建客户端
// 	client, err:= ethclient.Dial(httpUrl)
// 	require.NoError(t, err)

// 	// nonce获取
// 	nonce, err := client.PendingNonceAt(context.Background(), fromAddr)

// 	// 认证信息组装
// 	auth := bind.NewKeyedTransactor(fromPrivkey)
// 	//auth,err := bind.NewTransactor(strings.NewReader(mykey),"111")
// 	auth.Nonce = big.NewInt(int64(nonce))
// 	auth.Value = amount     // in wei
// 	//auth.Value = big.NewInt(100000)     // in wei
// 	auth.GasLimit = gasLimit // in units
// 	//auth.GasLimit = uint64(0) // in units
// 	auth.GasPrice = gasPrice
// 	auth.From = fromAddr

// 	// 交易创建
// 	tx := types.NewTransaction(nonce,toAddr,amount,gasLimit,gasPrice,[]byte{})

// 	// 交易签名
// 	signedTx ,err:= auth.Signer(types.HomesteadSigner{}, auth.From, tx)
// 	//signedTx ,err := types.SignTx(tx,types.HomesteadSigner{},fromPrivkey)
// 	require.NoError(t, err)

// 	// 交易发送
// 	serr := client.SendTransaction(context.Background(),signedTx)
// 	if serr != nil {
// 		fmt.Println(serr)
// 	}

// 	// 等待挖矿完成
// 	bind.WaitMined(context.Background(),client,signedTx)

// }
