package main

import (
	"bufio"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"os"
	"fmt"
)
var (
	input = kingpin.Flag("input", "Input log file").Short('i').Required().String()
	output = kingpin.Flag("output", "Output directory").Short('o').String()
)

type tmp struct{
	aa string
	bb int
	cc string
}

func main(){
	if basicReg == nil || infoRegs == nil{
		fmt.Println("Regulation initialization faild.")
		return
	}


	kingpin.Parse()
	if err != nil {
		err.Error()
		return
	}

	// default value of unset Flag (*output for eg.) is ""
	// default value of string inside struct is ""

	//Begin parse
	f, err := os.Open(*input)
	if err != nil{
		fmt.Println("Open input failed.")
		return
	}

	defer f.Close()
	reader := bufio.NewReader(f)
	recorder := CreateRecorder()
	recorder.AddMapCounter()
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			fmt.Println("log end")
			break
		}
		event, err := ParseLine(string(line))
		if err == nil {
			recorder.AddEvent(event)
		}else{
			switch err.(type){
			case *InvalidLogLine:
				//Do nothing
			case *ParseFailed, *UnknownReg:
				fmt.Println(err.Error())
			default:
				fmt.Println(err.Error())
			}
		}
	}
	recorder.PrintCounter()
	recorder.CheckSelf()
	ok := recorder.SetEventsPeer()
	if !ok {
		fmt.Println("Set peer failed")
	}
	fmt.Println(recorder.selfPeer)

	if *output != ""{
		analyzer := CreateCSVAnalyzer(*output, recorder)
		analyzer.AnalyzeBLK()
	}

	//test peer name
	//peername := &peerName{
	//	names:  make(map[string]string),
	//}
	//for i:=0; i<300; i++{
	//	peername.Add(string(i))
	//	peername.Add(string(i))
	//}
	//fmt.Println(Stringmap2json(peername.names))
}
