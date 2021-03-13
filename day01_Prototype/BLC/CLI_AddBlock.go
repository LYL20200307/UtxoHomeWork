package BLC

import (
	"fmt"
	"os"
)

//定义命令行接口
func (cli *CLI) addBlock(Txs []*Transaction){//添加区块
	if !dbExists(dbName){
		fmt.Println(dbName,"数据库不存在")
		os.Exit(1)
	}
	bc:=blockchainObject()
	bc.AddBlockDB(Txs)
	defer bc.DB.Close()
}

