package BLC

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
	"math/big"
	"os"
)

//将int64转byte[],封装为工具类
func IntToHex(num int64) []byte {
	buffer:=new(bytes.Buffer)
	err:=binary.Write(buffer,binary.BigEndian,num)//将num的值以大端模式写入buffer,如果没有错误返回nil
	if err!=nil{
		log.Panicf("int transact to []byte failed %v\n",err)
	}
	return buffer.Bytes()
}


//json格式转切片,将json字符串转成字符串切片
func JsonToSlice(jsonstring string)[]string{
	var sArr []string

	if err := json.Unmarshal([]byte(jsonstring),&sArr);err != nil{
		log.Panic(err)
	}
	return sArr
}


//确定输入的命令行参数是否符合要求
func isValidArgs(){
	if len(os.Args)<2{//os.args  获得命令行输入的参数，第一个参数通常是文件名，此时os.args的长度是1
		PrintUsage()
		os.Exit(1)//退出程序
	}
}

func ArriveFirstBlock(block *Block)bool{
	var hashInt big.Int
	hashInt.SetBytes(block.HashPrevBlock)
	if hashInt.Cmp(big.NewInt(0))==0{
		return true
	}
	return false

}