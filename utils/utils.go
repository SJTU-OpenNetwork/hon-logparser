package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const Version = "1.2.0"

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

/**
 * Traverse a directory and find all the log files.
 * Log file must have the format "uniqueId_index.log"
 */
func ListLogFiles(dirPath string, fileMap map[string][]string) map[string][]string {
	fstat, err := os.Stat(dirPath)
	if err != nil {
		fmt.Printf(err.Error())
		return fileMap
	}

	if fstat.IsDir() {
		files, err := ioutil.ReadDir(dirPath)
		if err != nil {
			return fileMap
		}
		for _, f := range files {
			fileMap = ListLogFiles(path.Join(dirPath, f.Name()), fileMap)
		}
	} else {
		fileName := fstat.Name()
		//filePath := path.Join(dirPath, fileName)
		//fmt.Printf("Traverse to %s\n", fileName)
		fileExtInfo := strings.Split(fileName, ".")
		if len(fileExtInfo) < 2 || fileExtInfo[1] != "log" {
			fmt.Printf("%s is not a log file\n", dirPath)
			return fileMap
		}

		fileNameInfo := strings.Split(fileExtInfo[0], "_")
		name := fileNameInfo[0]

		_, ok := fileMap[name]
		if ok {
			fileMap[name] = append(fileMap[name], dirPath)
		} else {
			fileMap[name] = []string{dirPath}
		}
	}

	return fileMap
}

// Input: /root/xxx/abcd_1.log
// Output: abcd
func GetLogName(filePath string) string{
	_, fileName := path.Split(filePath)
	return strings.Split(fileName, ".")[0]
}