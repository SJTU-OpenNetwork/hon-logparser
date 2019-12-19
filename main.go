package main

import (
	"bufio"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"os"
	"fmt"
)
var (
	input = kingpin.Flag("input", "Input log file").Short('i').String()
	output = kingpin.Flag("output", "Output directory").Short('o').String()
)

func main(){
	if basicReg == nil || infoRegs == nil{
		fmt.Println("Regulation initialization faild.")
		return
	}


	kingpin.Parse()
	//defer (*input).Close() //打开文件出错处理
	fmt.Printf("input file: %s, output to %s\n", *input, *output)

	//Initialize regular expression
	//String that begin with time stamp. Extract the timestamp.
	//basicReg, err = regexp.Compile(`([\d -\.:]{26}) ([A-Z]*) ([a-z-_\.]*) ([a-z-_:\.0-9]*) (.*)`)
	if err != nil {
		err.Error()
		return
	}


	//Create output directory
	exist, err := PathExists(*output)
	if err != nil{
		err.Error()
		return
	}
	if !exist {
		err = os.Mkdir(*output, os.ModePerm)
		if err != nil{
			err.Error()
			return
		}
	}

	//Begin parse
	f, err := os.Open(*input)
	if err != nil{
		fmt.Println("Open input failed.")
		return
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			fmt.Println("log end")
			break
		}
		//fmt.Println(line)
		//fmt.Println("!!!!!!!!")
		info, err := extractBasic(string(line))
		if err == nil{
			testParse(info)
		}
	}
	//extractBasic("2019-12-18 01:25:45.289748 INFO tex-service service.go:482: pubsub service listener started for /textile/threads/2.0.0/Thread/12D3KooWSBLLrCiAidDzPwcZs2ak1u7ZXhC3geU8RAg61eRzAWm7")
}
