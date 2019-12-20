package main

import (
	"bufio"
	"path"

	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"io/ioutil"
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

func parseFile(filePath string) *Recorder{
	f, err := os.Open(filePath)
	if err != nil{
		panic(err)
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
				//fmt.Println(err.Error())
				panic(err)
			default:
				//fmt.Println(err.Error())
				panic(err)
			}
		}
	}
	recorder.PrintCounter()
	recorder.CheckSelf()
	ok := recorder.SetEventsPeer()
	if !ok {
		fmt.Println("Set peer failed")
	}
	return recorder
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
	var recorder *Recorder
	fordir, _ := os.Stat(*input)
	if fordir.IsDir(){
		recorders := make([]*Recorder, 0)
		files, err := ioutil.ReadDir(*input)
		if err != nil{
			panic(err)
		}

		for _, f := range files{
			//if path.Ext(f.Name())=="log" {
				fmt.Println("parse "+f.Name())
				recorders = append(recorders, parseFile(path.Join(*input, f.Name())))
			//}else{
			//	fmt.Println("not a log file " + f.Name())
			//}
		}
		recorder = MergeRecorders(recorders)
	} else {
		recorder = parseFile(*input)
	}

	if *output != ""{
		analyzer := CreateCSVAnalyzer(*output, recorder,
			[]string{"BLKRECV", "BLKCANCEL", "WANTRECV", "BLKSEND",
				"WANTSEND","TKTSEND","ACKSEND","TKTRECV","TKTREJECT", "TKTACCEPT","ACKRECV"})
		analyzer.AnalyzeAll()
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
