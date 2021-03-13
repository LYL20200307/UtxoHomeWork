package BLC

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

/*
持久化用到的功能，写、读、查询
读写事务用DB.Update(),View()
哈希值是索引，区块信息是值
*/
const dbName="blockchain.db"//数据库名称
const bkName="blocks"//桶名也就是表

type BlockChainDB struct{
	DB *bolt.DB//数据库对象
	BlockHash []byte//记录最新区块的哈希值(原先的是用切片来记录，因为改成数据库了，所以用最新区块的哈希值作为一个索引，来回溯之前所有的区块)
}
//创建或打开数据库
func CreateDB() *bolt.DB{
	db,err:=bolt.Open(dbName, 0600, nil)
	if err!=nil{
		log.Panicf("create db %s failed %v\n",dbName,err)
	}
	return db
}


/*数据库内容: 32-byte 该块哈希 -> 序列化后区块信息
			‘l’ -> 链中最后一个区块的hash值*/
//使用数据库存储区块链,创建一个创世区块
func DBCreateBlockChain(address string) *BlockChainDB{

	if dbExists(dbName){//文件存在
		fmt.Println("已存在一个区块链")
		db,err:=bolt.Open(dbName, 0600, nil)
		if err!=nil{
			log.Panicf("open db %s failed %v\n",dbName,err)
		}
		var bc *BlockChainDB
		err=db.View(func(tx *bolt.Tx) error {//以只读方式打开数据库，并获取数据库里的最新区块值
			bcTable:=tx.Bucket([]byte(bkName))//打开表
			latestHash:=bcTable.Get([]byte("latest"))//------重要,得到当前区块链的最新区块，然后将这个最新区块的信息共享
			bc=&BlockChainDB{db,latestHash}
			return nil
		})
		if nil!=err{
			log.Panic(err)
		}
		return bc
	} else{
		db,err:=bolt.Open(dbName, 0600, nil)
		if err!=nil{
			log.Panicf("create db %s failed %v\n",dbName,err)
		}
		var blockhash []byte
		//创建表(桶)
		err=db.Update(func(tx *bolt.Tx)error{//形参是个函数
			bcTable:=tx.Bucket([]byte(bkName))//查找桶
			if bcTable==nil{//没找到表
				bcTable,err=tx.CreateBucket([]byte(bkName))
				if err!=nil{
					log.Panicf("create bucket %s failed %v\n",bkName,err)
				}
			}
			if bcTable!=nil{//如果有表则添加数据
				txCoinBase:=CoinBaseTransaction(address)
				genesisBlock:=CreateGensisBlock([]*Transaction{txCoinBase})//创建创世区块
				err:=bcTable.Put(genesisBlock.BlockHash,genesisBlock.Serialize())//将块数据存进表中
				if err!=nil{
					log.Panic(err)
				}
				//记录下最新区块的哈希值，因为是创建新的链，所以最新区块是创世区块
				err=bcTable.Put([]byte("latest"),genesisBlock.BlockHash)
				if err!=nil{
					log.Panic(err)
				}
				blockhash=genesisBlock.BlockHash
			}
			return nil
		})
		if err!=nil{
			log.Panic(err)
		}

		return &BlockChainDB{db,blockhash}//返回数据库和最新区块的哈希值
	}


}

//使用数据库添加区块
/*
1.从桶中找到最后一个区块
2.对区块进行反序列化，也就是把byte转为block
3.根据取出的信息，创建新的区块
4.把新区块序列化存进数据库
*/
func (blockchain *BlockChainDB)AddBlockDB(txs []*Transaction){
	err:=blockchain.DB.Update(func (tx *bolt.Tx)error{
		table:=tx.Bucket([]byte(bkName))
		if table!=nil{
			blockByte:=table.Get(blockchain.BlockHash)//得到最新区块的区块信息([]byte类型)
			latestblock:=Deserialize(blockByte)

			newblock:=NewBlock(latestblock.Height+1,latestblock.BlockHash,txs)
			err:=table.Put(newblock.BlockHash,newblock.Serialize())
			if err!=nil{
				log.Panic("put the data of new block into Dbfailed! %v\n", err)
			}
			err=table.Put([]byte("latest"),newblock.BlockHash)
			if err!=nil{
				log.Panic("put the hash of the newest block into Dbfailed! %v\n", err)
			}
			blockchain.BlockHash=newblock.BlockHash
		}
		return nil
	})
	if err!=nil{
		log.Panic(3,err)
	}
}

//遍历区块链中的区块
func (bc *BlockChainDB)PrintChain(){
	fmt.Println("——————————————打印区块链———————————————————————")
	var curBlock *Block//当前区块
	var curHash []byte=bc.BlockHash//当前区块哈希
	for{
		err:=bc.DB.View(func(tx *bolt.Tx) error {
			table:=tx.Bucket([]byte(bkName))
			if table!=nil{//有区块在链里,获取区块数据
				blockByte:=table.Get(curHash)
				curBlock=Deserialize(blockByte)
				fmt.Printf("\tHeigth : %d\n", curBlock.Height)
				fmt.Printf("\tTimeStamp : %v\n", curBlock.NTime)
				fmt.Printf("\tPrevBlockHash : %x\n", curBlock.HashPrevBlock)
				fmt.Printf("\tHash : %x\n", curBlock.BlockHash)
				fmt.Printf("\tNonce : %d\n", curBlock.Nonce)
				fmt.Printf("\tTxs : %v\n", curBlock.Txs)
				for _,tx :=range curBlock.Txs{
					fmt.Printf("\t\t tx-hash:%x\n",tx.TxHash)
					fmt.Printf("\t\t 这笔交易的输入信息\n")
					for _,vin:=range tx.Vins{
						fmt.Printf("\t\t\t vin-txhash:%x\n",vin.TxHash)//  若是%x格式输出，字符串会被转为0-255之间的数，然后变为十六进制
						fmt.Printf("\t\t\t vin-scriptSig:%v\n",vin.ScriptSig)
						fmt.Printf("\t\t\t vin-vout:%v\n",vin.Vout)
					}
					fmt.Printf("\t\t 这笔交易的输出信息\n")
					for _,vout :=range tx.Vouts{
						fmt.Printf("\t\t\t vout-value:%d\n",vout.Value)
						fmt.Printf("\t\t\t vout-scriptPubkey:%v\n",vout.ScriptPubkey)
					}
				}
				fmt.Println("—————————————————————————————————————————————")
			}else{
				fmt.Println("区块链为空")
			}
			return nil
		})
		if nil!=err{
			log.Panic(err)
		}
		if ArriveFirstBlock(curBlock)==true{//到达创世区块
			break
		}
		curHash=curBlock.HashPrevBlock
	}
}
//查询该地址为花费的utxo
//将所有与该地址相关的输出值加起来就是余额


//遍历整条区块链的所有交易，判断每一笔交易的输出是否符合：1.需要查找的地址2.是否没被花费(也就是没在输入里被引用)
//[]*TxOutPut
//实现一个区块包含多笔交易
//在得到所有已花费输出的基础上，遍历所有交易的输入来得到utxo
func (bc *BlockChainDB)UTXOs(addr string,txs []*Transaction) []*Utxo{//创世区块没有输入交易，因此需要对此增加判断
	var UTXOS  []*Utxo
	//遍历数据库，查找与地址相关的所有交易。利用迭代器,不断获取下一个区块
	bcit:=bc.Iterator()
	//获取指定地址所有已花费输出
	spentTXOutputs:=bc.spentOutPuts(addr)
	//后续要增加缓存和数据库迭代的判断

	for _,tx:=range txs{//缓存迭代,查找缓存中的已花费输出(因为交易还没被打包)
		//判断coinbase交易
		if !tx.IsCoinbaseTransaction(){
			for _,vin:=range tx.Vins{//遍历输入，得到已花费的钱
				//判断用户
				if vin.checkPubkey(addr){
					//因为被交易在输入中引用，所以属于被花费的钱
					key:=hex.EncodeToString(vin.TxHash)
					spentTXOutputs[key]=append(spentTXOutputs[key],vin.Vout)
				}
			}
		}
	}
	//获得缓存中的utxo(遍历输出，匹配spentTXOutputs得到utxo)
	for _,tx:=range txs{
		WorkCacheTx:
		for index,vout:=range tx.Vouts{
			if vout.checkPubkey(addr){
				if len(spentTXOutputs)!=0{//如果该账户有钱被花费
					var isUtxoTx bool//记录交易是否有被引用
					for txHash,indexArr:=range spentTXOutputs{
						txHashStr:=hex.EncodeToString(tx.TxHash)
						if txHash==txHashStr{//当前遍历到的交易有输出被其他交易的输入引用
							isUtxoTx=true
							var isSpentUtxo bool//用来记录索引号是否匹配上，索引号和交易哈希都匹配才是被花费
							for _,voutIndex:=range indexArr{
								if index==voutIndex{//输出交易的索引和引用输入的索引相同
									//输出被引用(该utxo被引用)
									isSpentUtxo=true
									continue WorkCacheTx
								}
							}
							if isSpentUtxo==false{//与某个地址相关的所有utxo都没有被花费
								utxo:=&Utxo{tx.TxHash,index,vout}
								UTXOS=append(UTXOS,utxo)
							}
						}
					}
					if isUtxoTx==false{//交易没被引用过
						utxo:=&Utxo{tx.TxHash,index,vout}
						UTXOS=append(UTXOS,utxo)
					}
				}
			}else{//没有被花费
				utxo:=&Utxo{tx.TxHash,index,vout}
				UTXOS=append(UTXOS,utxo)
			}
		}
	}

	//数据库迭代，不断获取下一区块
	for{
		block:=bcit.Next()
		for _,tx :=range block.Txs{//获取区块中的每一笔交易
			work://添加一个跳转
			for index,vout:=range tx.Vouts{//获取每笔交易里的输出和其对应的索引值
				if vout.checkPubkey(addr){
					if len(spentTXOutputs)!=0{
						var isSpentoutput bool//默认false
						for txhash,indexArray:=range spentTXOutputs{//获得被引用交易的哈希值以及索引号
							for _,i :=range indexArray{//得到每一个索引号
								if txhash==hex.EncodeToString(tx.TxHash) && index==i{//交易哈希相同并且索引号相同则说明这笔钱被花费了(当前输出被其他交易作为输入引用)
									isSpentoutput=true
									continue work//去判断下一个vout
								}
							}
						}
						if isSpentoutput==false{//遍历完所有被引用的输入后发现有的输出依然没有被花费
							utxo:=&Utxo{tx.TxHash,index,vout}
							UTXOS=append(UTXOS,utxo)
							//UTXOS=append(UTXOS,vout)
						}
					}else{//此时该地址没有任何交易被作为输入引用，因此当前地址的所有输出都可以添加到utxo中
						utxo:=&Utxo{tx.TxHash,index,vout}
						UTXOS=append(UTXOS,utxo)
						//UTXOS=append(UTXOS,vout)
					}
				}

			}
		}
		if ArriveFirstBlock(block)==true{//到达创世区块
			break
		}
	}
	return UTXOS
}

//在获取utxo后添加一个专门用来查询余额的函数
func(bc *BlockChainDB)getBalance(addr string)int{
	var amount int
	utxos:=bc.UTXOs(addr,[]*Transaction{})
	for _,utxo:=range utxos{
		amount=amount+int(utxo.TxOutput.Value)
	}
	return amount
}


//获取指定地址所有已花费输出
func (bc *BlockChainDB)spentOutPuts(addr string)map[string][]int{//因为地址，交易哈希值，索引号确定一个引用，而现在已知地址了，所以只需要返回交易哈希值和索引号
	spentUtxo:=make(map[string][]int)
	bcit:=bc.Iterator()
	for{
		block:=bcit.Next()//每次获得一个块
		for _,tx:=range block.Txs{//遍历块中的每一个交易
			if !tx.IsCoinbaseTransaction(){//判断是否为创世区块，因为创世区块没有引用任何交易的输出
				for _,vin:=range tx.Vins{//遍历交易的每笔输入
					if vin.checkPubkey(addr){//在输入中找符合查询地址的的引用交易
						hashstr:=hex.EncodeToString(vin.TxHash)//将哈希值转为字符串
						//在交易的输入中被引用的交易就被添加到已花费的集合中
						spentUtxo[hashstr]=append(spentUtxo[hashstr],vin.Vout)//Vout是交易的索引号，会出现一个地址在一笔交易中有两笔钱的情况(多笔零钱合起来用)
					}
				}
			}
		}
		if ArriveFirstBlock(block)==true{//到达创世区块
			break
		}
	}
	return spentUtxo//返回该地址所有已经被引用的收入
}

func dbExists(name  string) bool {//判断是否已存在一条区块链，如果存在，被二次调用创建创世区块时，默认添加到已有的区块链后面
	if _,err:=os.Stat(name);os.IsNotExist(err){//os.stat函数返回文件属性和错误，另一个则用来判断文件是否存在(如果文件不存在则返回true)。这俩组合使用
		return false
	}
	return true//文件存在
}

//该数据库的key和value都只能存储字节数组，
//将block数据转为字节数组型叫做数据序列化，反之叫数据反序列化
func (block *Block)Serialize() []byte{
	var result bytes.Buffer//具有读写方法的字节大小可变的缓冲区
	encoder:=gob.NewEncoder(&result)//创建一个编码器
	//对block进行编码

	err:=encoder.Encode(block)
	if err!=nil{
		log.Panic(err)
	}
	return result.Bytes()
}

//反序列化,字节变区块
func Deserialize(blockByte []byte) *Block{
	var block Block
	var reader=bytes.NewReader(blockByte)
	decoder:=gob.NewDecoder(reader)//创建一个解码器
	//将区块解码
	err:=decoder.Decode(&block)

	if err!=nil{
		log.Panic(err)
	}
	return &block
}