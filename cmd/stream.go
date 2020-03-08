package cmd

import (
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-logparser/analyzer"
	"github.com/SJTU-OpenNetwork/hon-logparser/utils"
	"os"
	"path"
	"github.com/SJTU-OpenNetwork/hon-textile/stream"
)

func streamCmdRecord(filePath string, outDir string) error {
	parser, err := analyzer.NewParser(); if err != nil {return err}
	res := &utils.CidFilter{make(map[string]interface{})}

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
				// For each file
				recorder, err := analyzer.RecorderFromFile(f, parser, []string{stream.TAG_BLOCKSEND, stream.TAG_BLOCKRECEIVE,})
				if err != nil {return err}
				filter := recorder.GetCidFilter()
				res = utils.MergeTwoCidFilter(res, filter)
			}
		}

	} else {
		recorder, err := analyzer.RecorderFromFile(filePath, parser, []string{stream.TAG_BLOCKSEND, stream.TAG_BLOCKRECEIVE,})
		if err != nil {
			return err
		}
		filter := recorder.GetCidFilter()
		res = utils.MergeTwoCidFilter(res, filter)
		// End save
	}

	err = res.ToFile(path.Join(outDir, "cids.txt")); if err != nil {return err}
	return nil
}