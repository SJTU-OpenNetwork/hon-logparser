package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path"
	"time"
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
		//test.TestJson()
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
		fmt.Printf("Make directory %s\n", outDir)
		err = os.Mkdir(outDir, os.ModePerm); if err!=nil {return err}
	}

	testMap["key_int"] = 1
	testMap["key_string"] = "aaa"
	testMap["key_slice"] = []string{"a", "b", "c"}

	js, err := json.MarshalIndent(testMap, "", "  ")
	jsFilePath := path.Join(outDir, "testmap.json")
	_, err = WriteBytes(jsFilePath, js); if err !=nil {return err}

	data, err := ReadBytes(jsFilePath); if err != nil {return err}
	fmt.Printf("%s\n", string(data))

	fmt.Printf("Test to load and write recorder file")
	rec := CreateRecorder()
	evt1 := &BitswapEvent{
		Peer:      "peerid",
		Type:      "eventType",
		Time:      time.Now(),
		Direction: []string{"From", "To"},
		Info:      map[string]interface{}{
			"key_int": 1,
			"key_string": "aa",
			"key_slice": []string{"v1","v2"},
		},
	}

	evt2 := &BitswapEvent{
		Peer:      "peerid",
		Type:      "eventType",
		Time:      time.Now(),
		Direction: []string{"From", "To"},
		Info:      map[string]interface{}{
			"key_int": 1,
			"key_string": "aa",
			"key_slice": []string{"v1","v2"},
		},
	}

	rec.AddEvent(evt1)
	rec.AddEvent(evt2)

	js2, err := json.MarshalIndent(rec, "", "  ")
	jsFilePath2 := path.Join(outDir, "testRecorder.json")
	_, err = WriteBytes(jsFilePath2, js2); if err !=nil {return err}

	data2, err := ReadBytes(jsFilePath2); if err != nil {return err}
	fmt.Printf("%s\n", string(data2))

	return nil
}