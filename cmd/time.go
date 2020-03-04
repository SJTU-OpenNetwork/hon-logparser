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
	var timeInfos []*analyzer.TimeInfo
	if err != nil {
		return err
	}
	fstat, err := os.Stat(filePath); if err != nil {return err}

	if fstat.IsDir() {
		fmt.Printf("List all log files in %s\n", filePath)
		fileMap := utils.ListLogFiles(filePath, make(map[string][]string))
		for _, v := range fileMap {
			for _, f := range v {
				timeInfos, err = analyzer.TimeFromFile(parser, f); if (err != nil) {return err}
			}
		}

	} else {
		timeInfos, err = analyzer.TimeFromFile(parser, filePath); if (err != nil) {return err}
	}

	if timeInfos != nil {
		for _, ti := range timeInfos {
			ti.PrintOut()
		}
	}

	return nil
}
