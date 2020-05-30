package cmd

import (
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-logparser/analyzer"
	"github.com/SJTU-OpenNetwork/hon-logparser/utils"
	"os"
	"path"
)

func filter(input string, output string, reg string) error {
	_, err := os.Stat(input)
	if os.IsNotExist(err) {
		return fmt.Errorf("No such file %s", input)
	}

	filter, err := analyzer.NewFilterFromString(reg)
	if err != nil {
		return err
	}

	err = utils.CheckOrCreateDir(output)
	if err != nil {
		return err
	}

	fileMap := utils.ListLogFiles(input, make(map[string][]string))
	for logname, loglist := range fileMap {
		for i, f := range loglist {
			outPath := path.Join(output, fmt.Sprintf("%s_%d.log", logname, i))
			fmt.Printf("Filter %s\n", f)
			inFile, err := os.Open(f)
			if err != nil {
				fmt.Println(err)
				continue
			}
			outFile, err := os.Open(outPath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = filter.FilterFile(inFile, outFile)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = inFile.Close(); if err!= nil {fmt.Println(err)}
			err = outFile.Close(); if err!= nil {fmt.Println(err)}
		}
	}
	return nil
}
