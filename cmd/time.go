package cmd

import (
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-logparser/analyzer"
	"github.com/SJTU-OpenNetwork/hon-logparser/utils"
	"os"
)

// get receiving time for all transforming files in all log files
// write the result into outDir
func time(filePath string, outDir string) error {
	parser, err := analyzer.NewParser(); if err != nil {return err;}

	// list all log files
	// get receiving time for each file
	// write them to files. one peer each file.

	fmt.Printf("Make directory for %s", outDir)
	err = os.MkdirAll(outDir, os.ModePerm)
	if err != nil {
		return err
	}
	fstat, err := os.Stat(filePath); if err != nil {return err}

	if fstat.IsDir() {
		fmt.Printf("List all log files in %s\n", filePath)
		fileMap := utils.ListLogFiles(filePath, make(map[string][]string))
		for _, v := range fileMap {
			for _, f := range v {
				_, err= analyzer.TimeFromFile(parser, f); if (err != nil) {return err}
			}
		}

	} else {
		_, err = analyzer.TimeFromFile(parser, filePath)
	}

	return nil
}
