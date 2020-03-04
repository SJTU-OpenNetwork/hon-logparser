package analyzer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)



func TimeFromFile(parser *Parser, filePath string) ([]*TimeInfo, error){
	f, err := os.Open(filePath)
	if err != nil{
		fmt.Printf("Cannot open %s\n", filePath)
		return nil, err
	}
	//fmt.Printf("File opened\n")
	defer f.Close()

	reader := bufio.NewReader(f)
	res := make([]*TimeInfo, 0)
	lineNum := 0
	fmt.Printf("Begin parse for %s\n", filePath)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			fmt.Println("Read to EOF\n")
			break
		}
		info := parser.ParseLineForTime(string(line))
		if info != nil {
			res = append(res, info)
		}
		lineNum++
		if lineNum % 10 == 0 {
			fmt.Printf("line: %d\r", lineNum)
		}
	}
	return res, nil
}

// TimeInfo contains basic info extract from one line in log file
// It may be useless to define a specific struct for it if we only use the time info as is.
// But it would be helpful if we want to do further analysing about time info.
type TimeInfo struct {
	Cid string
	Ms int
	Bytes int
	BytePerMs int
	PeerId string
}

func (ti *TimeInfo) PrintOut(){
	js, err := json.Marshal(ti)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	} else {
		fmt.Printf("%s\n", string(js))
	}
}

func (ti *TimeInfo) ToCSVLine() []byte {
	str := fmt.Sprintf("%s, %d, %d, %d\n", ti.PeerId, ti.Ms, ti.Bytes, ti.BytePerMs)
	return []byte(str)
}