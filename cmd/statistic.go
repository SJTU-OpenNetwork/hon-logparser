package cmd

import (
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-logparser/analyzer"
	"github.com/SJTU-OpenNetwork/hon-logparser/utils"

	//"io/ioutil"
	"os"
	//"path"
)

func statistic(filePath string, outDir string) error {
	fmt.Printf("Make directory for %s", outDir)
	err := os.MkdirAll(outDir, os.ModePerm)
	if err != nil {
		return err
	}
	fstat, err := os.Stat(filePath); if err != nil {return err}
	if fstat.IsDir() {
		//fmt.Printf("Unimplement\n")
		//files, err := ioutil.ReadDir(filePath)
		fmt.Printf("List all log files in %s\n", filePath)
		fileMap := utils.ListLogFiles(filePath, make(map[string][]string))
		for k, v := range fileMap {
			fmt.Printf("%s:", k)
			for _, f := range v {
				fmt.Printf(f)
			}
			fmt.Printf("\n")
		}

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

func statisticRecurisiveDir(dirPath string) ([]*analyzer.Statistic, error) {
	// find the path of all the files:

}

func statisticFile(filePath string) (*analyzer.Statistic, error) {
	// parse the whole file
	//
	parser, err := analyzer.NewParser(); if err != nil {return nil, err}
	statistics, err := analyzer.CountForFile(parser, filePath)
	return statistics, nil
}