package main

import (
	"bufio"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"io/ioutil"
	"os"
	"path"
)
var (
	app	  = kingpin.New("hon-logparser", "A command-line tool used to parse log of HON project.")
	input = app.Flag("input", "Input log file").Short('i').Required().String()
	output = app.Flag("output", "Output directory").Short('o').Required().String()

	parseCmd = app.Command("parse", "Parse the log infos")
)

type tmp struct{
	aa string
	bb int
	cc string
}

func initDir(){
	// Initialize output directory
	ok, err := PathExists(*output)
	if err != nil {
		panic(err)
	}

	if !ok{
		err = os.Mkdir(*output, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	counterDir := path.Join(*output, "counters")
	ok, err = PathExists(counterDir)
	if err != nil {
		panic(err)
	}

	if !ok{
		err = os.Mkdir(counterDir, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
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
				// Happends when it is not a passerable line.
				// The format of passerable line:
				// - 20xx-xx-xx xx:xx:xx.xxx DEBUG tex-core xxx.go:111: [MSGRECV] xxxx
				// - ([\d -\.:]{26}) ([A-Z]*) ([a-z-_\.]*) ([a-z-_:\.0-9A-Z]*) \[([A-Z]*)\] (.*)
			case *ParseFailed:
				// Happends when the event info cannot be pasered
				// It always implies wrong regulation expression so we panic it directly.
				panic(err)
			case *UnknownReg:
				// Happends when there occurrs some unknown event type.
				// Such as [XXXSEND]
				//fmt.Println(err.Error())

			default:
				//fmt.Println(err.Error())
				panic(err)
			}
		}
	}
	//recorder.PrintCounter()
	recorder.CheckSelf()
	if recorder.selfPeer == SELF {
		fmt.Println("Cannot get peer id through check self")
	}
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

	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))
	initDir()

	// Create and merge recorder
	fordir, _ := os.Stat(*input)
	var recorder *Recorder
	if fordir.IsDir(){
		recorders := make([]*Recorder, 0)
		files, err := ioutil.ReadDir(*input)
		if err != nil{
			panic(err)
		}

		for _, f := range files{
			fmt.Println("parse "+f.Name())
			tmpRecorder := parseFile(path.Join(*input, f.Name()))
			tmpRecorder.SaveCounter(path.Join(*output, "counters", tmpRecorder.selfPeer+".json"))
			recorders = append(recorders, parseFile(path.Join(*input, f.Name())))
		}
		recorder = MergeRecorders(recorders)
	} else {
		recorder.SaveCounter(path.Join(*output, "counters", recorder.selfPeer+ ".json"))
		recorder = parseFile(*input)
	}

	//switch kingpin.MustParse(app.Parse(os.Args[1:])){
	switch cmd{
	case parseCmd.FullCommand():
		fmt.Println("Do Parse.")
		analyzer := CreateCSVAnalyzer(*output, recorder,
			[]string{"BLKRECV", "BLKCANCEL", "WANTRECV", "BLKSEND",
				"WANTSEND","TKTSEND","ACKSEND","TKTRECV","TKTREJECT", "TKTACCEPT","ACKRECV"})
		analyzer.AnalyzeAll()
		analyzer.AnalyzerRECVTree()
	}

	// default value of unset Flag (*output for eg.) is ""
	// default value of string inside struct is ""

}
