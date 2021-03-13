package BLC

import (
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
) // flag包提供了一系列解析命令行参数的功能接口



//命令行交互接口，可以手动添加区块和打印区块链
//相当于设计一个交互的前端
type CLI struct {
	//BC *BlockChainDB//也可以不要，写一个获取区块链对象的函数来代替
}

//获取区块链对象的函数
func blockchainObject()*BlockChainDB{
	var bc *BlockChainDB
	db,err:=bolt.Open(dbName, 0600, nil)
	if err!=nil{
		log.Panic(err)
	}
	err=db.View(func(tx *bolt.Tx) error {
		table:=tx.Bucket([]byte(bkName))
		if table!=nil{
			blockHash:=table.Get([]byte("latest"))
			bc=&BlockChainDB{DB: db,BlockHash: blockHash}
		}

		return nil
	})
	if err!=nil{
		log.Panic(err)
	}
	return bc
}

func (cli *CLI)createBlockChain(addr string){
	DBCreateBlockChain(addr)
}

//获取指定地址的余额
func(cli *CLI)getBalance(addr string){
	//获取区块链对象
	bc:=blockchainObject()
	defer bc.DB.Close()//关闭实例对象
	amount:=bc.getBalance(addr)
	fmt.Printf("\t地址[%s]的余额：[%d]\n",addr,amount)
}


//flag.Parse() 把用户传递的命令行参数解析为对应变量的值
//使用flag 获取命令行输出参数
func (cli *CLI)Run(){
	//检查命令行输入参数是否正确
	isValidArgs()
	//新建命令
	//addBlockCmd:=flag.NewFlagSet("addblock",flag.ExitOnError)//给命令取名
	printBlockCmd:=flag.NewFlagSet("printchain",flag.ExitOnError)//给命令取名
	createBlockChainCmd:=flag.NewFlagSet("createblockchain",flag.ExitOnError)//给命令取名
	sendTxCmd:=flag.NewFlagSet("send",flag.ExitOnError)
	getBalanceCmd:=flag.NewFlagSet("getbalance",flag.ExitOnError)

	//设置命令行参数
	//flagAddBlockArg:=addBlockCmd.String("data","send 100 BTC to everyone","交易数据")
	flagCreateBlockArg:=createBlockChainCmd.String("address","send 10 BTC to everyone","指定接受系统奖励的矿工地址")
	flagSendFromArg:=sendTxCmd.String("from","","转账源地址")
	flagSendToArg:=sendTxCmd.String("to","","接受转账地址")
	flagSendAmountArg:=sendTxCmd.String("amount","","转账金额")
	flagGetBalanceArg:=getBalanceCmd.String("address","","需要进行查询余额的地址")
	//解析命令
	switch os.Args[1] {
	case "createblockchain":
		err:=createBlockChainCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panicf("parse cmd of create block chain failed! %v\n", err)
		}
	case "send":
		err:=sendTxCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panicf("parse cmd of send Tx failed! %v\n", err)
		}
	case "printchain":
		err:=printBlockCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panicf("parse cmd of printchain failed! %v\n", err)
		}
	case "getbalance":
		err:=getBalanceCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panicf("parse cmd of getbalance failed! %v\n", err)
		}
	default:
		PrintUsage()
		os.Exit(1)
	/*case "addblock":
	err:=addBlockCmd.Parse(os.Args[2:])
	if err!=nil{
		log.Panicf("parse cmd of add block failed! %v\n", err)
	}*/
	}

	//根据解析的命令执行相应操作
/*	if addBlockCmd.Parsed(){//返回是否f.Parse已经被调用过
		if *flagAddBlockArg==""{//给了命令，没给相应参数
			PrintUsage()
			os.Exit(1)
		}
		cli.addBlock(*flagAddBlockArg)
	}*/

	if createBlockChainCmd.Parsed(){//返回是否f.Parse已经被调用过
		if *flagCreateBlockArg==""{//给了命令，没给相应参数
			PrintUsage()
			os.Exit(1)
		}
		cli.createBlockChain(*flagCreateBlockArg)
	}
	if sendTxCmd.Parsed(){
		if *flagSendFromArg==""{
			println("转账地址不能为空")
			PrintUsage()
			os.Exit(1)
		}
		if *flagSendToArg==""{
			println("接收转账不能为空")
			PrintUsage()
			os.Exit(1)
		}
		if *flagSendAmountArg==""{
			println("转账金额不能为空")
			PrintUsage()
			os.Exit(1)
		}
		fmt.Printf("\tFROM:[%s]\n",JsonToSlice(*flagSendFromArg))
		fmt.Printf("\tTo:[%s]\n",JsonToSlice(*flagSendToArg))
		fmt.Printf("\tAmount:[%s]\n",JsonToSlice(*flagSendAmountArg))
		cli.send(JsonToSlice(*flagSendFromArg),JsonToSlice(*flagSendToArg),JsonToSlice(*flagSendAmountArg))
	}
	if printBlockCmd.Parsed(){
		cli.printChain()
	}
	if getBalanceCmd.Parsed(){
		if *flagGetBalanceArg==""{
			PrintUsage()
			os.Exit(1)
		}
		cli.getBalance(*flagGetBalanceArg)
	}
}