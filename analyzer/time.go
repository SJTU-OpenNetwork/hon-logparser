package analyzer

import (
	"bufio"
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
	//lineNum := 0
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
	}
	return res, nil
}
