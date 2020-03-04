package analyzer

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// TimePerFile contains information about file receiving time
// for each pic_cid in one log file.
type TimePerFile struct {

}

func TimeFromFile(parser *Parser, filePath string) (*TimePerFile, error){
	f, err := os.Open(filePath)
	if err != nil{
		fmt.Printf("Cannot open %s\n", filePath)
		return nil, err
	}
	//fmt.Printf("File opened\n")
	defer f.Close()

	reader := bufio.NewReader(f)
	//lineNum := 0
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			fmt.Println("Read to EOF\n")
			break
		}
		info := parser.ParseLineForTime(string(line))
		if info != nil {
			for k, v := range info {
				fmt.Printf("%s: %s\n", k, v)
			}
			fmt.Printf("\n")
		}
	}
	return nil, nil
}
