package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path"
)

type cmdsMap map[string]func() error

var (
	appCmd =kingpin.New("parser", "Parser is a tool to analyse the log files of hon-textile project")
)

func run() error {
	fmt.Printf("Initialize commands\n")
	cmds := make(cmdsMap)
	appCmd.UsageTemplate(kingpin.CompactUsageTemplate)

	// Test for json file load and write
	testCmd := appCmd.Command("test", "Test the load and write of json files.")
	testDir := testCmd.Arg("ourdir", "Output directory.").Required().String()
	cmds[testCmd.FullCommand()] = func() error {
		return testJson(*testDir)
	}

	// commands
	cmd := kingpin.MustParse(appCmd.Parse(os.Args[1:]))
	for key, value := range cmds {
		if key == cmd {
			return value()
		}
	}
	return nil
}

func testJson(outDir string) error{
	fmt.Printf("Test to load and write json file to %s\n", outDir)
	testMap := make(map[string]interface{})
	ok, err := PathExists(outDir); if err != nil {return err}
	if !ok {
		fmt.Printf("Make directory %s", outDir)
		err = os.Mkdir(outDir, os.ModePerm); if err!=nil {return err}
	}

	testMap["key_int"] = 1
	testMap["key_string"] = "aaa"
	testMap["key_slice"] = []string{"a", "b", "c"}

	js, err := json.Marshal(testMap)
	jsFilePath := path.Join(outDir, "testmap.json")
	_, err = WriteBytes(jsFilePath, js); if err !=nil {return err}

	return nil
}