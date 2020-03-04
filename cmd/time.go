package cmd

import (
	"bufio"
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-logparser/analyzer"
	"github.com/SJTU-OpenNetwork/hon-logparser/utils"
	"os"
	"path"
)


// get receiving time for all transforming files in all log files
// write the result into outDir
func time(filePath string, outDir string) error {
	parser, err := analyzer.NewParser(); if err != nil {return err;}

	// list all log files
	// get receiving time for each file
	// write them to files. one peer each file.

	fmt.Printf("Make directory for %s", outDir)
	exists, err := utils.PathExists(outDir); if err != nil {return err}

	if !exists {
		err = os.MkdirAll(outDir, os.ModePerm); if err != nil {return err}
	}

	fstat, err := os.Stat(filePath); if err != nil {return err}

	totalInfos := InitInfoStore()	// Contains all the timeinfos
	if fstat.IsDir() {
		fmt.Printf("List all log files in %s\n", filePath)
		fileMap := utils.ListLogFiles(filePath, make(map[string][]string))
		for name, v := range fileMap {
			for _, f := range v {
				timeInfos, err := analyzer.TimeFromFile(parser, f); if (err != nil) {return err}
				totalInfos.add(timeInfos, name)
			}
		}

	} else {
		name := utils.GetLogName(filePath)
		timeInfos, err := analyzer.TimeFromFile(parser, filePath); if (err != nil) {return err}
		totalInfos.add(timeInfos, name)
	}

	//if timeInfos != nil {
	//	for _, ti := range timeInfos {
	//		ti.PrintOut()
	//	}
	//}
	err = totalInfos.toCSV(outDir); if err != nil {return err}

	return nil
}

type infoStore struct {
	// data[cid] [info from peer 1, info from peer 2, ...]
	// peer is distinguished by the name get from log file name.
	data map[string] []*analyzer.TimeInfo
}

func InitInfoStore() *infoStore{
	return &infoStore{
		data: make(map[string][]*analyzer.TimeInfo),
	}
}

func (store *infoStore) get(cid string) ([]*analyzer.TimeInfo, bool) {
	tmplist, ok := store.data[cid]
	return tmplist, ok
}

func (store *infoStore) set(cid string, infos []*analyzer.TimeInfo) {
	store.data[cid] = infos
}

func (store *infoStore) add(infos []*analyzer.TimeInfo, peerId string) {
	if (infos != nil) {
		for _, info := range infos {
			info.PeerId = peerId // Note that peerId may not be the true peerId.
								 // It depends on how user named the log files.
			cid := info.Cid
			tmpInfoSlice, ok := store.get(cid)
			if !ok {
				store.set(cid, []*analyzer.TimeInfo{info})
			} else {
				store.set(cid, append(tmpInfoSlice, info))
			}
		}
	}
}

// Write to several csv files in outDir
func (store *infoStore) toCSV(outDir string) error {
	for cid, infos := range store.data {
		csvPath := path.Join(outDir, cid+".csv")
		//csvPath := path.Join(a.outputDir, "All.csv")
		fo, err := os.Create(csvPath); if err != nil {return err}
		w:= bufio.NewWriter(fo)
		_, err = w.Write(store.getCSVHeader()); if err != nil {return err}

		for _, info := range infos {
			_, err = w.Write(info.ToCSVLine()); if err != nil {return err}
		}

		err = w.Flush(); if err != nil {return err}
		err = fo.Close(); if err != nil {return err}
	}
	return nil
}

func (store *infoStore) getCSVHeader() []byte {
	return []byte(", time/ms, bytes, bytesPerMs\n")
}
