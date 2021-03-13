package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)
//我们想要区块能够顺序打印，BoltDB 允许对一个 bucket 里面的所有 key 进行迭代，但是所有的 key 都以字节序进行存储，用于实现数据很大的时候，我们并不想让所有的数据都直接加载到内存中，需要一个迭代器，来一个一个读取区块
type BlockChainIterator struct{//当前迭代的块哈希
	DB *bolt.DB
	currentHash []byte
}

//创建迭代器对象
func (blc *BlockChainDB) Iterator() *BlockChainIterator{
	return &BlockChainIterator{blc.DB,blc.BlockHash}
}

func (bcit *BlockChainIterator)Next() *Block{//返回区块链中的下一个区块
	var block *Block
	err:=bcit.DB.View(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte(bkName))//打开表
		if b!=nil{
			blockByte:=b.Get(bcit.currentHash)
			block=Deserialize(blockByte)
			bcit.currentHash=block.HashPrevBlock
		}
		return nil
	})
	if err!=nil{
		log.Panic(err)
	}
	return block
}