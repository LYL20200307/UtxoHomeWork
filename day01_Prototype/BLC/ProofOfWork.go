package BLC
import(
	"fmt"
	"bytes"
	"crypto/sha256"
	"math/big"
)
const targetBit=16//前四位哈希值为0
type ProofOfWork struct{
	//要进行共识的区块
	block *Block
	//共识算法的上界
	target big.Int
}
//创建新的pow对象
func CreateNewPowOfWork(block *Block) *ProofOfWork{
	//设置算法难度
	target:=big.NewInt(1)
	target=target.Lsh(target,256-targetBit)//将target左移240位，其中16表示2byte
	return &ProofOfWork{block,*target}
}

//运行pow算法，包括将随机数添加入区块头生成哈希,返回哈希值和随机数
func (pow *ProofOfWork)Run()(int64,[]byte){
	var nonce int64=0
	hashInt:=new(big.Int)//用于将字节数组转为int，用于确定哈希值是否符合难度，var hashInt big.Int
	var hash [32]byte//存储哈希值的字节数组

	for{
		data:=pow.prepareData(nonce)

		hash=sha256.Sum256(data)
		fmt.Printf("nonce=%v,hash=%x\n",nonce,hash)
		hashInt.SetBytes(hash[:])//传入切片，方便哈希值与target比较
		if pow.target.Cmp(hashInt)==1{//哈希值小于目标值时结束寻找随机数
			break
		}
		nonce++
	}
	fmt.Printf("求哈希的次数为%d\n",nonce)
	return nonce,hash[:]
}
//拼接区块头获取字节数组
func (pow *ProofOfWork)prepareData(nonce int64) []byte{
	data:= bytes.Join([][]byte{
		IntToHex(pow.block.Height),
		IntToHex(pow.block.NTime),
		pow.block.HashPrevBlock,
		pow.block.HashTransactions(),//pow.block.data
		IntToHex(nonce),
		IntToHex(targetBit),
	},[]byte{})//二维转一维

	return data
}

//因为准备数据中的tx是结构体而不是字节数组，所以需要转为字节类型
func (block *Block)HashTransactions() []byte{
	var txHashes [][]byte

	for _,tx:=range block.Txs{
		txHashes=append(txHashes,tx.TxHash)//将每个交易的哈希值加入到这个切片中
	}
	txhash:=sha256.Sum256(bytes.Join(txHashes,[]byte{}))
	return txhash[:]
}