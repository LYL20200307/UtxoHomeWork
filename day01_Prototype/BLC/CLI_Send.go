package BLC

import (
	"fmt"
	"os"
)

//发起交易
//通过命令行send发起交易调用挖矿
func (cli *CLI)send(from,to,amount []string){
	//判断数据库是否存在
	if !dbExists(dbName){
		fmt.Println(dbName,"数据库不存在")
		os.Exit(1)
	}
	bc:=blockchainObject()
	defer bc.DB.Close()
	if len(from)!=len(to)||len(from)!=len(amount){
		fmt.Println("交易参数输入有误")
		os.Exit(1)
	}
	bc.MineNewBlock(from,to,amount)
}