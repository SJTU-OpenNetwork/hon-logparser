package cmd

import (
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-logparser/analyzer"
	"os"
)

func statistic(filePath string, outDir string) error {
	fstat, err := os.Stat(filePath); if err != nil {return err}
	if fstat.IsDir() {
		fmt.Printf("Unimplement\n")
	} else {
		sta, err := statisticFile(filePath)
		if err != nil {
			return err
		}
		err = sta.SaveToDisk(outDir)
		if err != nil {
			return err
		}
	}
	//	var recorder *Recorder
	//	if fordir.IsDir(){
	//		recorder = parseRecursiveDir(*input)
	//	} else {
	//		//recorder.SaveCounter(path.Join(*output, "counters", recorder.selfPeer+ ".json"))
	//		recorder = parseFile(*input)
	//	}
	return nil
}

func statisticFile(filePath string) (*analyzer.Statistic, error) {
	// parse the whole file
	//
	parser, err := analyzer.NewParser(); if err != nil {return nil, err}
	statistics, err := analyzer.CountForFile(parser, filePath)
	return statistics, nil
}