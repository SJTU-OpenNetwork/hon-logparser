package main

import (
	"encoding/json"
	"fmt"
)

type Counter interface{
	Count(*Event)
	String() string
	//SaveCounter(savepath string) error
}

type MapCounter struct{
	datastore map[string]int
}

func CreateMapCounter() *MapCounter{
	return &MapCounter{
		datastore: make(map[string]int),
	}
}

func (c *MapCounter) Count(event *Event){
	_, ok := c.datastore[event.Type]
	if ok {
		c.datastore[event.Type] += 1
	} else {
		c.datastore[event.Type] = 1
	}
}

func (c *MapCounter) String() string{
	jsonString, err := json.MarshalIndent(c.datastore, "", "\t")
	if err != nil{
		fmt.Println(err.Error())
		return ""
	}
	return string(jsonString)
}

//func (c *MapCounter) SaveCounter(savepath string) error{
//
//}