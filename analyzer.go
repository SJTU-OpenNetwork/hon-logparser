package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"path"
)

type CSVAnalyzer struct{
	outputDir string
	recorder *Recorder
	names *peerName					//handle peer rename
	eventList *list.List			//eventList that store only the filtered events
	cidMap map[string] *list.List	//event indexed by cid
	//peerIndex map[string]int
}

type peerName struct{
	names map[string]string
	peerIndex map[string]int
	indexPeer map[int]string
	number int
}

func CreateCSVAnalyzer(outputDir string, recorder *Recorder) *CSVAnalyzer {
	return &CSVAnalyzer{
		outputDir: outputDir,
		recorder:  recorder,
		names:  &peerName{
			names: make(map[string]string),
			peerIndex : make(map[string]int),
			indexPeer : make(map[int]string),
		},
		eventList : list.New(),
		cidMap : make(map[string] *list.List),
	}
}



func (p *peerName) Add(peerId string){
	p.add(peerId)
}

func (p *peerName) Get(peerId string) string{
	//p.add(peerId)
	return p.names[peerId]
}

func (p *peerName) add(peerId string) {
	_, ok := p.names[peerId]
	var tmp  = p.number
	var res string
	if !ok {
		if p.number == 0{
			res = "A"
		} else {
			for tmp > 0 {
				r := tmp % 26
				res = res + string('A' + r)
				tmp /= 26
			}
		}
		p.names[peerId] = reverseString(res)
		p.peerIndex[peerId] = p.number
		p.indexPeer[p.number] = peerId
		p.number++
	}
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}


/**
 *
 */
func (a *CSVAnalyzer) AnalyzeBLK(){
	// Initialize directory
	ok, err := PathExists(a.outputDir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if !ok{
		err = os.Mkdir(a.outputDir, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	// Analyze block
	for e := a.recorder.eventList.Front(); e != nil; e = e.Next(){
		switch e.Value.(*BitswapEvent).Type{
		case "BLKRECV":
			// [BLKRECV] Cid <cid>, From <peerid>
			event := e.Value.(*BitswapEvent)
			//fmt.Println(Map2json(event.Info))
			cid := event.Info["Cid"].(string)
			l, ok := a.cidMap[cid]
			if !ok {
				l = list.New()
				a.cidMap[cid] = l
				a.names.Add(event.Peer)
				a.names.Add(event.Info["From"].(string))
				//a.cidList = append(a.cidList, cid)
			}
			l.PushBack(event)
			a.eventList.PushBack(event)
		}
	}

	// Write to file
	for cid, l := range a.cidMap{
		a.writeBLKCSV(a.outputDir, cid, l)
	}
}

/**
 * Get the header for this analyzer
 * (It may contains different peers.)
 * Format:
 * 		Time, Cid, p1, p2, ..., pn \n
 */
func (a *CSVAnalyzer) csvHeader()string{
	res := "Time, Cid"
	for i:=0; i<len(a.names.indexPeer); i++{
		res = res + fmt.Sprintf(", %s", a.names.Get(a.names.indexPeer[i]))
	}
	res = res + "\n"
	return res
}

/**
 * Get the contains of csv line.
 */
func (a *CSVAnalyzer) csvContain(length int, index int, contains string) string{
	result := ""
	for i:=0; i<length; i++{
		if i == index{
			result = result + ", " + contains
		} else {
			result = result + ", "
		}
	}
	result = result + "\n"
	return result
}

func (a *CSVAnalyzer) csvBLKLine(e *BitswapEvent) string{
	if e.Type != "BLKRECV"{
		panic(&WrongEventType{eventType: e.Type})
	}

	receiver := e.Peer
	publisher := e.Info["From"].(string)
	cid := e.Info["Cid"].(string)
	result := fmt.Sprintf("%s, %s",e.Time.String(), cid)
	publisherInd := a.names.peerIndex[publisher]
	result = result + a.csvContain(a.names.number, publisherInd, a.names.Get(receiver))
	return result
}

func (a *CSVAnalyzer) writeBLKCSV(outDir string, cid string, l *list.List){
	csvPath := path.Join(outDir, "BLKRECV.csv")
	fo, err := os.Create(csvPath)
	if err != nil {
		panic(err)
	}

	defer func(){
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	w:= bufio.NewWriter(fo)
	header := a.csvHeader()
	if _, err := w.Write([]byte(header)); err != nil{
		panic(err)
	}

	for e := a.eventList.Front(); e != nil; e = e.Next(){
		line := a.csvBLKLine(e.Value.(*BitswapEvent))
		if _, err := w.Write([]byte(line)); err != nil{
			panic(err)
		}
	}

	if err = w.Flush(); err != nil {
		panic(err)
	}
}