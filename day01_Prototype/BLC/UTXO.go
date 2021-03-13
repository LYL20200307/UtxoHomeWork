package BLC

import (
	"encoding/hex"
	"fmt"
	"os"
)

//UTXO结构
type Utxo struct {
	//交易哈希
	TxHash []byte
	//交易的输出索引
	index int
	//交易的输出(包含交易金额和用户名)
	TxOutput *TxOutPut
}

//通过utxo查询进行转账
//amount要查多少钱，from要查询的地址，txs:缓存中的交易列表,用于多笔交易
//返回，查到的金额和被花费的utxo
//更新数据库中指定地址的utxo数量  []*TxOutPut
func (bc *BlockChainDB)FindSpendableUTXO(amount int64,from string,txs []*Transaction) (int64,map[string][]int){
	//已被使用的utxo
	spendableUtxo:=make(map[string][]int)
	utxos:=bc.UTXOs(from,txs)
	var value int64
	for _,utxo:=range utxos{
		value+=utxo.TxOutput.Value
		hash:= hex.EncodeToString(utxo.TxHash)//交易的哈希转为字符串
		spendableUtxo[hash]=append(spendableUtxo[hash],utxo.index)
		if value>=amount{
			break
		}
	}
	if value<amount{
		fmt.Printf("%s 余额不足,总额：%d，需要：%d\n", from,value,amount)
		os.Exit(1)
	}
	return value,spendableUtxo
}
