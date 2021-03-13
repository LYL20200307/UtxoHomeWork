package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"github.com/boltdb/bolt"
	"strconv"
)

//交易结构
type Transaction struct{
	//交易的哈希值
	TxHash []byte
	//交易输入
	Vins []*TxInPut//指针构成的切片
	//交易输出
	Vouts []*TxOutPut
}

//交易输入
type TxInPut struct {
	//交易哈希
	TxHash []byte
	//引用上一笔交易的输出索引
	Vout int
	//解锁脚本,也就是签名
	ScriptSig string
}

//交易输出
type TxOutPut struct{
	//金额
	Value int64
	//收款人的地址
	ScriptPubkey string
}

//输入输出的验证功能
//验证当前的utxo是否属于指定地址
func (txoutput *TxOutPut)checkPubkey(addr string)bool{
	return txoutput.ScriptPubkey==addr
}

//同理，也要验证输入时引用的钱是否属于某个地址
func(txinput *TxInPut)checkPubkey(addr string)bool{
	return txinput.ScriptSig==addr
}

//生成交易哈希
func (tx *Transaction)HashTransaction(){
	var result bytes.Buffer
	encoder:=gob.NewEncoder(&result)//将编码后的数据写入result
	err:=encoder.Encode(tx)//对tx进行编码,也就是把交易进行序列化
	if nil != err {
		log.Panicf("tx hash generate failed! %v\n", err)
	}
	txhash:=sha256.Sum256(result.Bytes())
	tx.TxHash=txhash[:]
}

func CoinBaseTransaction(addr string) *Transaction{
	txinput:=&TxInPut{[]byte{},-1,"Gensis PData"}//创世区块没有需要引用的交易
	txoutput:=&TxOutPut{10,addr}
	txcoinbase:=&Transaction{nil,[]*TxInPut{txinput},[]*TxOutPut{txoutput}}
	txcoinbase.HashTransaction()//获得这笔交易的哈希值
	return txcoinbase
}

//生成转账交易
//将之前写死的交易输入变为真正的交易输入
//
func NewSimpleTransaction(from string,to string,amount int64,bc *BlockChainDB,txs []*Transaction)*Transaction{

	var txInputs []*TxInPut//输入列表
	var txOutputs []*TxOutPut//输出列表

	//调用utxo函数，说明会花费多少钱的utxo
	money,spendableUTXO:=bc.FindSpendableUTXO(amount,from,txs)
	fmt.Println("money:%v",money)

	for txhash,indexArr:=range spendableUTXO{
		txhashByte,err:=hex.DecodeString(txhash)//将字符串转为byte数组
		if err!=nil{
			log.Panic("字符串转换失败")
		}
		//遍历索引列表
		for _,index:=range indexArr{
			txInput:=&TxInPut{txhashByte,index,from}
			txInputs=append(txInputs,txInput)
		}
	}
	//输入(谁转的)
/*	txInput:=&TxInPut{[]byte("ed312bdea2e5df5e4ac1dd14683ecdd5" +
		"834700f39b336d6ab3b5855f6c0c08f8"),0,from}//引用上一个交易的哈希(暂时写一个固定的哈希)，上一笔交易里的第0个(这里引用的创世交易，所以默认是0)
	txInputs=append(txInputs,txInput)*/


	//输出(转给谁)
	txOutput:=&TxOutPut{amount,to}
	txOutputs=append(txOutputs,txOutput)

	//输出(找零)
	if money>=amount{
		txOutput =&TxOutPut{money-amount,from}
		txOutputs=append(txOutputs,txOutput)
	}else{
		log.Panicf("余额不足\n")
	}


	//生成交易的哈希值
	transaction:=Transaction{nil,txInputs,txOutputs}
	transaction.HashTransaction()

	return &transaction
}

//挖矿功能，通过挖矿生成新区块
//通过接受交易生成区块
func(bc *BlockChainDB)MineNewBlock(from,to,amount []string){//交易信息应该有多个，所以要使用切片
	//调用函数生成新交易
	var txs []*Transaction
	//遍历交易的参与者
	for index,addr:=range from{
		//将字符串类型转为int
		value,_:=strconv.Atoi(amount[index])
		tx:=NewSimpleTransaction(addr,to[index],int64(value),bc,txs)//因为只是简单的测试一下，所以不需要遍历切片
		txs=append(txs,tx)
	}



	//先写生成区块
	var newblock *Block
	err:=bc.DB.Update(func (tx *bolt.Tx)error{
		table:=tx.Bucket([]byte(bkName))
		if table!=nil{
			blockByte:=table.Get(bc.BlockHash)//得到最新区块的区块信息([]byte类型)
			latestblock:=Deserialize(blockByte)

			newblock=NewBlock(latestblock.Height+1,latestblock.BlockHash,txs)
			err:=table.Put(newblock.BlockHash,newblock.Serialize())
			if err!=nil{
				log.Panic("put the data of new block into Dbfailed! %v\n", err)
			}
			err=table.Put([]byte("latest"),newblock.BlockHash)
			if err!=nil{
				log.Panic("put the hash of the newest block into Dbfailed! %v\n", err)
			}
			bc.BlockHash=newblock.BlockHash
		}
		return nil
	})
	if err!=nil{
		log.Panic(3,err)
	}
}

//判断创世区块
func (tx *Transaction)IsCoinbaseTransaction() bool{
	return tx.Vins[0].Vout==-1&&len(tx.Vins[0].TxHash)==0
}