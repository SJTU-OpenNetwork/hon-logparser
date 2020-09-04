package cmd

import (
	//"encoding/json"
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-logparser/analyzer"
	"github.com/SJTU-OpenNetwork/hon-logparser/utils"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	//"path"
	//"time"
)

type cmdsMap map[string]func() error

var (
	appCmd =kingpin.New("parser", "Parser is a tool to analyse the log files of hon-textile project")
)

func Run() error {
	fmt.Printf("Initialize commands\n")
	cmds := make(cmdsMap)
	appCmd.UsageTemplate(kingpin.CompactUsageTemplate)

	// Test for json file load and write
	testCmd := appCmd.Command("test", "Test the load and write of json files.")
	testDir := testCmd.Arg("ourdir", "Output directory.").Required().String()
	cmds[testCmd.FullCommand()] = func() error {
		return utils.TestJson(*testDir)
	}

	// For statistic infos
	statisticCmd := appCmd.Command("statistic", "Do some statistical analysis for several log files.")
	statisticInputDir := statisticCmd.Arg("input", "Input directory or file. " +
		"Format file name as \"uniqueId_index.extension\" to distinguish log from different peers.").Required().String()
	statisticOutputDir := statisticCmd.Arg("output", "Output directory for result. A new directory would be created if not exists.").Required().String()
	//statisticMaintainName := statisticCmd.Flag("maintain", "Whether to maintain the input file to output file name." +
	//	"This is pretty usefull when you use file name to distinguish different kinds of log file.").Bool()
	statisticCidFilter := statisticCmd.Flag("cidFilter", "File or directory path of cid filter. " +
		"If given, logparser would only extract block information of cids within filter.").String()
	cmds[statisticCmd.FullCommand()] = func() error {
		return statistic(*statisticInputDir, *statisticOutputDir,  *statisticCidFilter)
	}

	// For time infos
	timeCmd := appCmd.Command("time", "Extract the receiving time info from log files.")
	timeInputDir := timeCmd.Arg("input", "Input directory or file. " +
		"Format file name as \"uniqueId_index.extension\" to distinguish log from different peers.").Required().String()
	timeOutputDir := timeCmd.Arg("output", "Output directory for result. A new directory would be created if not exists.").Required().String()
	cmds[timeCmd.FullCommand()] = func () error {
		return time(*timeInputDir, *timeOutputDir)
	}

	// For stream analyze
	streamCmd :=  appCmd.Command("stream", "Analyzer for stream protocol.")
	streamCidRecordCmd := streamCmd.Command("recordCid", "Record block cids transferred by stream protocol.")
	streamCidRecordInputDir := streamCidRecordCmd.Arg("input", "Input directory or file. ").Required().String()
	streamCidRecordOutputDir := streamCidRecordCmd.Arg("output", "Output directory for result. A new directory would be created if not exists.").Required().String()
	cmds[streamCidRecordCmd.FullCommand()] = func () error {
		return streamCidRecord(*streamCidRecordInputDir, *streamCidRecordOutputDir)
	}

	versionCmd := appCmd.Command("version", "Version of hon-logparser.")
	cmds[versionCmd.FullCommand()] = func () error {
		fmt.Printf("Version: %s\n", utils.Version)
		return nil
	}

	// For regular expression filter.
	filterCmd := appCmd.Command("filter", "Filter lines in input and output the matched lines to files in output directory")
	filterInput := filterCmd.Arg("input", "Input directory or file.").Required().String()
	filterOutput := filterCmd.Arg("output", "Output directory. A new directory would be created if no one exists").Required().String()
	filterRegular := filterCmd.Arg("regular", "Regular expression usd for filter (Can be a simple substring)").Required().String()
	cmds[filterCmd.FullCommand()] = func () error {
		return filter(*filterInput, *filterOutput, *filterRegular)
	}

	// For tree analyze
	treeCmd := appCmd.Command("tree", "Analyzer the distribution tree for ticket based method.")
	treeInput := treeCmd.Arg("input", "Input directory or file.").Required().String()
	treeOutput := treeCmd.Arg("output", "Output directory.").Required().String()
	cmds[treeCmd.FullCommand()] = func() error {
		fordir, err := os.Stat(*treeInput)
		if err != nil {
			return err
		}
		var recorder *analyzer.Recorder
		if fordir.IsDir() {
			recorder = parseRecursiveDir(*treeInput)
		} else {
			//recorder.SaveCounter(path.Join(*output, "counters", recorder.selfPeer+ ".json"))
			recorder = parseFile(*treeInput)
		}


		fmt.Println("Do Parse.")
		CsvAnalyzer := analyzer.CreateCSVAnalyzer(*treeOutput, recorder,
			[]string{"BLKRECV", "BLKCANCEL", "WANTRECV", "BLKSEND",
				"WANTSEND", "TKTSEND", "ACKSEND", "TKTRECV", "TKTREJECT", "TKTACCEPT", "ACKRECV"})
		CsvAnalyzer.AnalyzeAll()
		CsvAnalyzer.AnalyzerRECVTree()
		return nil
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
/*

func testJson(outDir string) error{
	fmt.Printf("Test to load and write json file to %s\n", outDir)
	testMap := make(map[string]interface{})
	ok, err := main.PathExists(outDir); if err != nil {return err}
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

	js2, err := json.MarshalIndent(*rec, "", "  ")
	jsFilePath2 := path.Join(outDir, "testRecorder.json")
	_, err = WriteBytes(jsFilePath2, js2); if err !=nil {return err}

	data2, err := ReadBytes(jsFilePath2); if err != nil {return err}
	fmt.Printf("%s\n", string(data2))

	return nil
}
*/
