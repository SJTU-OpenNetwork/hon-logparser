package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)


/**
 * PathExists return true if the path exists.
 */
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}


func Map2json(info map[string]interface{}) string {
	jsonString, err := json.MarshalIndent(info, "", "\t")
	if err != nil{
		fmt.Println(err.Error())
		return ""
	}
	return string(jsonString)
}

func Stringmap2json(info map[string]string) string{
	jsonString, err := json.MarshalIndent(info, "", "\t")
	if err != nil{
		fmt.Println(err.Error())
		return ""
	}
	return string(jsonString)
}

func Contains(sl []string, v string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

/**
 * WriteBytes write data into a file
 */
func WriteBytes(filePath string, b []byte) (int, error) {
	//os.MkdirAll(path.Dir(filePath), os.ModePerm)
	fw, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer fw.Close()
	return fw.Write(b)
}

/**
 * ReadBytes read data from a file
 */
func ReadBytes(filePath string) ([]byte, error) {
	fw, err := os.Open(filePath); if err !=nil {return nil, err}
	defer fw.Close()
	return ioutil.ReadAll(fw)
}