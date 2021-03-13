package BLC

import (
	"time"
)
//区块结构
type Block struct {
	Height        int64  //区块高度
	NTime         int64  //时间戳，代表区块时间
	BlockHash     []byte //区块哈希值，可以用字符串表示，这里用字节数组.使用数组可以将数据不以字节数来计算长度，而是以字符数计算长度
	HashPrevBlock []byte //上一区块哈希值
	/*交易数据data，先设计为[]byte,后期用transaction实现*/
	//Data  []byte
	Txs 		  []*Transaction
	Nonce int64
}
//新建区块
func NewBlock(height int64,hashPrevBlock []byte,txs []*Transaction) *Block{
	var block Block

	block=Block{
		Height:        height,
		BlockHash:     nil,//因为区块的哈希是算出来的，所以没算之前应该赋值为空，当有计算哈希的函数时才将其更改
		NTime:         time.Now().Unix(),
		HashPrevBlock: hashPrevBlock,
		Txs:          txs,//交易数据
		Nonce:         0,
	}
	pow:=CreateNewPowOfWork(&block)
	nonce,hash:=pow.Run()
	block.BlockHash=hash[:]
	block.Nonce =nonce

	return &block
}
//计算区块哈希值
/*func (block *Block) SetHash(){
	heightByte:=IntToHex(block.Height)
	timeByte:=IntToHex(block.nTime)
	hashData:=bytes.Join([][]byte{
		block.BlockHash,
		block.HashPrevBlock,
		block.Data,
		heightByte,
		timeByte,
	},[]byte{})//后面的一维数组用于传入分隔符
	hash:=sha256.Sum256(hashData)//返回切片
	block.BlockHash=hash[:]//将切片转为数组
}*/

//生成创世区块，不属于任何一个类
func CreateGensisBlock(txs []*Transaction)*Block{
	return NewBlock(0,make([] byte,32,32),txs)

}

