package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

func TestJson(outDir string) error {
	fmt.Printf("Test to load and write json file to %s\n", outDir)
	testMap := make(map[string]interface{})
	ok, err := PathExists(outDir);
	if err != nil {
		return err
	}
	if !ok {
		fmt.Printf("Make directory %s\n", outDir)
		err = os.Mkdir(outDir, os.ModePerm);
		if err != nil {
			return err
		}
	}

	testMap["key_int"] = 1
	testMap["key_string"] = "aaa"
	testMap["key_slice"] = []string{"a", "b", "c"}

	js, err := json.MarshalIndent(testMap, "", "  ")
	jsFilePath := path.Join(outDir, "testmap.json")
	_, err = WriteBytes(jsFilePath, js);
	if err != nil {
		return err
	}

	data, err := ReadBytes(jsFilePath);
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", string(data))
	return nil
}
