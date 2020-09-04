package cmd

import (
	"bufio"
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-logparser/analyzer"
	"io"
	"io/ioutil"
	"os"
	"path"
)

func setRecorderSelf(recorder *analyzer.Recorder){
	recorder.CheckSelf()
	if recorder.GetSelfId() == analyzer.SELF {
		fmt.Println("Cannot get peer id through check self")
	}
	ok := recorder.SetEventsPeer()
	if !ok {
		fmt.Println("Set peer failed")
	}
}

func parseFile(filePath string) *analyzer.Recorder{
	f, err := os.Open(filePath)
	if err != nil{
		panic(err)
	}

	defer f.Close()
	reader := bufio.NewReader(f)
	parser, err := analyzer.NewParser()
	if err != nil {
		fmt.Println("Error when create parser: ", err)
		return nil
	}
	recorder := analyzer.CreateRecorder()
	recorder.AddMapCounter()
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			fmt.Println("log end")
			break
		}
		event, err := parser.ParseLine(string(line))
		if err == nil && event != nil {
			recorder.AddEvent(event)
		}else{
			/*
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
			 */
		}
	}
	//recorder.PrintCounter()

	return recorder
}

func parseRecursiveDir(dir string) *analyzer.Recorder{
	files, err := ioutil.ReadDir(dir)
	recorders := make([]*analyzer.Recorder, 0)

	var isLeaf = false

	if err != nil{
		panic(err)
	}

	for _, f := range files{
		var tmpRecorder *analyzer.Recorder
		tmpPath := path.Join(dir, f.Name())
		//fmt.Println(tmpPath)
		fordir, _ := os.Stat(tmpPath)
		if fordir.IsDir() {
			if isLeaf {
				panic("Invalid directory structure.")
			}
			tmpRecorder = parseRecursiveDir(tmpPath)
		} else {
			fmt.Println("Parse "+ tmpPath)
			isLeaf = true
			tmpRecorder = parseFile(tmpPath)
		}
		if tmpRecorder != nil{
			recorders = append(recorders, tmpRecorder)
		} else {
			fmt.Println(fmt.Sprintf("Get nil recorder when parse %s", tmpPath))
		}
	}

	if len(recorders) == 0{
		return nil
	} else {
		mergedRecorder := analyzer.MergeRecorders(recorders)
		if isLeaf {
			setRecorderSelf(mergedRecorder)
		}
		return mergedRecorder
	}
}
