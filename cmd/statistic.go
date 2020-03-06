package cmd

import (
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-logparser/analyzer"
	"github.com/SJTU-OpenNetwork/hon-logparser/utils"
	"path"
	"strings"

	//"io/ioutil"
	"os"
	//"path"
)

func statistic(filePath string, outDir string, maintain bool, cidFilterPath string) error {
	parser, err := analyzer.NewParser(); if err != nil {return err}
	if cidFilterPath != "" {
		cidFilter, err := utils.CidFilterFromFile(cidFilterPath); if err != nil {return err}
		//cidFilter.PrintOut()
		parser.SetCidFilter(cidFilter)
	}

	fmt.Printf("Make directory for %s", outDir)
	err = os.MkdirAll(outDir, os.ModePerm)
	if err != nil {
		return err
	}
	fstat, err := os.Stat(filePath); if err != nil {return err}

	if fstat.IsDir() {
		allStatistic := analyzer.NewEmpryStatistic("ALL")
		fmt.Printf("List all log files in %s\n", filePath)
		fileMap := utils.ListLogFiles(filePath, make(map[string][]string))
		for _, v := range fileMap {
			// For files of each peer
			staForOnePeer :=  analyzer.NewEmpryStatistic("")
			var savePath string
			if len(v) > 0{
				savePath, err = getStatisticFilePath(outDir, v[0]); if err != nil {return err}
			}else{
				continue
			}

			for _, f := range v {
				// For each file
				sta, err := analyzer.CountForFile(parser, f); if err != nil {return err}
				allStatistic = analyzer.MergeTwoStatistics(allStatistic, sta)
				staForOnePeer = analyzer.MergeTwoStatistics(staForOnePeer, sta)
			}
			// Save
			if maintain {
				err = staForOnePeer.SaveToDiskFile(savePath); if err != nil {return err	}
			}else {
				err = staForOnePeer.SaveToDisk(outDir);	if err != nil {	return err}
			}
		}
		err = allStatistic.SaveToDiskFile(path.Join(outDir, "ALL.json")); if err != nil {return err	}

	} else {
		sta, err := analyzer.CountForFile(parser, filePath); if err != nil {return err}
		// Save statistic file
		if maintain {
			savePath, err := getStatisticFilePath(outDir, filePath)
			if err != nil {
				return err
			}
			err = sta.SaveToDiskFile(savePath)
			if err != nil {
				return err
			}
		} else
		{
			err = sta.SaveToDisk(outDir)
			if err != nil {
				return err
			}
		}
		// End save
	}

	return nil
}

// Get the file path to save statistic file.
// This is useful when you want maintain the file name of log file as the file name of statistic file.
func getStatisticFilePath(outDir string, logPath string) (string, error) {
	_, logFile := path.Split(logPath)
	logFileName := strings.Split(logFile, ".")[0]
	//logFileBaseName := strings.Split(logFileName)
	return path.Join(outDir, logFileName + ".json"), nil
}

