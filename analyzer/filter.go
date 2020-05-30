package analyzer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

// Filter is used to filter lines in a file and output the lines matched some regulars to output file.
type Filter struct {
	reg *regexp.Regexp
}

func NewFilterFromString(regStr string) (*Filter, error) {
	return &Filter{reg: regexp.MustCompile(regStr)}, nil
}

func (filter *Filter) isMatch(line string) bool {
	// This return "" if no match (or matches "", which means regular expression is wrong).
	match := filter.reg.FindString(line)
	if match != "" {
		return true
	} else {
		return false
	}
}

func (filter *Filter) FilterFile(input *os.File, output *os.File) error {
	reader := bufio.NewReader(input)
	writer := bufio.NewWriter(output)
	defer func(){
		err := writer.Flush()
		if err != nil {
			fmt.Println(err)
		}
	}()
	for {
		line, _, err := reader.ReadLine()
		//fmt.Printf(string(line))
		if err == io.EOF {
			//fmt.Println("Read to EOF\n")
			break
		}
		lineStr := string(line)
		if filter.isMatch(lineStr) {
			_, err = writer.WriteString(lineStr + "\n")
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}
