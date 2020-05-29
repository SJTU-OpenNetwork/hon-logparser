package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// CidFilter is used to filter the result corresponding to cid.
// It can be generated from or wrote to file.
type CidFilter struct {
	data map[string]interface{}
}

func EmptyCidFilter() *CidFilter {
	return &CidFilter{data: make(map[string]interface{})}
}

// Generate a cid filter from file.
// Format of file contains:
//		pic_cid cid1 cid2 cid3 ....
func CidFilterFromFile(filePath string) (*CidFilter, error){
	f, err := os.Open(filePath)
	if err != nil{
		fmt.Printf("Cannot open %s\n", filePath)
		return nil, err
	}
	//fmt.Printf("File opened\n")
	defer f.Close()

	res := &CidFilter{data:make(map[string]interface{})}
	reader := bufio.NewReader(f)
	for {
		line, _, err := reader.ReadLine()
		//fmt.Printf(string(line))
		if err == io.EOF {
			fmt.Println("Read to EOF\n")
			break
		}

		cids := strings.Split(string(line), " ")
		for _, cid := range cids {
			res.Add(cid)
		}
	}
	return res, nil
}

func MergeTwoCidFilter(f1 *CidFilter, f2 *CidFilter) *CidFilter {
	f3 := &CidFilter{make(map[string]interface{})}
	for k,v := range f1.data {
		f3.data[k] = v
	}
	for k, v := range f2.data {
		f3.data[k] = v
	}
	return f3
}

func (filter *CidFilter) ToFile(filePath string) error {

	fo, err := os.Create(filePath); if err != nil {return err}
	w:= bufio.NewWriter(fo)

	for k, _ := range filter.data {
		_, err = w.Write([]byte(k+" ")); if err != nil {return err}
	}

	err = w.Flush(); if err != nil {return err}
	err = fo.Close(); if err != nil {return err}
	return nil
}

func (filter *CidFilter) Has(cid string) bool {
	_, ok := filter.data[cid]
	if ok {
		//fmt.Printf("Has %s\n", cid)
	} else {
		//fmt.Printf("Not has %s\n", cid)
	}
	return ok
}

func (filter *CidFilter) Add(cid string) {
	if !filter.Has(cid){
		filter.data[cid] = struct{}{}
	}
}

func (filter *CidFilter) PrintOut() {
	for k, _:= range filter.data {
		fmt.Printf("%s\n", k)
	}
}