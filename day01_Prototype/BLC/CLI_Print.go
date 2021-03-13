package BLC

import (
	"fmt"
	"os"
)
// 展示用法
func PrintUsage() {
	fmt.Println("Usage:")
	fmt.Printf("\tcreateblockchain -address Address -- 创建区块链.\n")
	//fmt.Printf("\taddblock -data DATA -- 交易数据\n")
	fmt.Printf("\tprintchain -- 输出区块链的信息\n")
	//通过命令转账
	fmt.Println("\tsend -from From -to To -amount Amount -- 发起转账")//send -from "[\"alice\"]" -to "[\"bob\"]" -amount "[\"1\"]"
	fmt.Println("\t\t-from FROM --转账源地址")
	fmt.Println("\t\t-to To --转账目的地址")
	fmt.Println("\t\t-amount Amount --转账金额")
	//查询指定地址的余额
	fmt.Printf("\tgetbalance -address Address --查询指定账户余额\n")
}

func (cli *CLI)printChain(){
	if !dbExists(dbName){
		fmt.Println(dbName,"数据库不存在")
		os.Exit(1)
	}
	bc:=blockchainObject()
	bc.PrintChain()
	defer bc.DB.Close()
}