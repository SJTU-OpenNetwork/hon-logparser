package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)



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

/**
 * Throw InvalidLogLine when try to parse a line cannot be parsed.
 * It happends sometimes as logs may contain many unformatted lines.
 */
type InvalidLogLine struct{}
func (e *InvalidLogLine) Error() string{
	return "Invalid log line"
}

/**
 * Throw ParseFailed when when parse some string failed.
 * It always implies wrong regulation expression
 */
type ParseFailed struct{
	expr string
	str string
}

func (e *ParseFailed) Error() string {
	return fmt.Sprintf("Failed to parse %s use %s", e.str, e.expr)
}

/**
 *
 */
type UnknownReg struct{
	reg string
}
func (e *UnknownReg) Error() string{
	return fmt.Sprintf("Unknown regular expression %s.", e.reg)
}

type WrongEventType struct{
	eventType string
}
func (e *WrongEventType) Error() string{
	return fmt.Sprintf("Get wrong event type %s.", e.eventType)
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

func parseTimestamp(str string) (time.Time, error){
	t, err := time.Parse(timeFotmat, str)
	if err != nil {
		fmt.Println(err.Error())
		return t, err
	}
	return t, nil
}