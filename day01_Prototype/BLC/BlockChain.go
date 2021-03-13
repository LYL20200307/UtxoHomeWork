package BLC

type BlockChain struct{
	Blocks []*Block//含有block指针的切片
}



func CreateBlockChainWithGensisBlock(txs []*Transaction) *BlockChain{
	block:=CreateGensisBlock(txs)
	return &BlockChain{Blocks: []*Block{block}}
}

//在当前区块链中添加区块
func (blockchain *BlockChain)AddBlock(height int64,hashPrevBlock []byte,txs []*Transaction){
	//var newblock *Block
	newBlock:=NewBlock(height,hashPrevBlock,txs)
	blockchain.Blocks=append(blockchain.Blocks,newBlock)

}