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
	filter []string
	//peerIndex map[string]int
}

type peerName struct{
	names map[string]string
	peerIndex map[string]int
	indexPeer map[int]string
	number int
}

func CreateCSVAnalyzer(outputDir string, recorder *Recorder, filter []string) *CSVAnalyzer {
	// Initialize directory
	ok, err := PathExists(outputDir)
	if err != nil {
		panic(err)
	}

	if !ok{
		err = os.Mkdir(outputDir, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
		}
	}


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
		filter: filter,
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

func (a *CSVAnalyzer) AnalyzeAll(){
	for e := a.recorder.eventList.Front(); e != nil; e = e.Next() {
		event :=  e.Value.(*BitswapEvent)
		if Contains(a.filter, event.Type){
			a.names.Add(event.GetPeer(event.Direction[0]))
			a.names.Add(event.GetPeer(event.Direction[1]))
			a.eventList.PushBack(event)
		}
	}
	a.writeAllCSV()
	a.writeNamePeer()
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

func (a *CSVAnalyzer) csvHeaderAll() string {
	res := "Time, Cid, Event"
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

func (a *CSVAnalyzer) csvAllLine(e *BitswapEvent) string{
	//fmt.Println("new line")
	cid := e.Info["Cid"].(string)
	result := fmt.Sprintf("%s, %s, %s", e.Time.String(), cid, e.Type)
	publisher := e.GetPeer(e.Direction[0])
	receiver := e.GetPeer(e.Direction[1])
	publisherInd := a.names.peerIndex[publisher]
	result = result + a.csvContain(a.names.number, publisherInd, a.names.Get(receiver))
	return result
}

//func (a *CSVAnalyzer) csvAllLine(e *BitswapEvent) string{
//	cid := e.Info["Cid"].(string)
//	result := fmt.Sprintf("%s, %s, %s", e.Time.String(), cid, e.Type)
//	var publisher string
//	var receiver string
//	switch e.Type{
//	case "BLKRECV", "BLKCANCEL", "WANTRECV":
//		// [BLKRECV] Cid <cid>, From <peerid>
//		// [BLKCANCEL] Cid <cid>, From <peerid>
//		// [WANTRECV] Cid <cid>, From <peerid>
//		publisher = e.Info["From"].(string)
//		receiver = e.Peer
//
//	case "BLKSEND", "WANTSEND":
//		// [BLKSEND] Cid <cid>, SendTo <peerid>
//		// [WANTSEND] Cid <cid>, SendTo <peerid>
//		publisher = e.Peer
//		receiver = e.Info["SendTo"].(string)
//
//	case "TKTSEND":
//		// [TKTSEND] Cid <cid>, SendTo <peerid>, TimeStamp <time>
//		publisher = e.Peer
//		receiver = e.Info["SendTo"].(string)
//
//	case "ACKSEND":
//		// [ACKSEND] Cid <cid>, Publisher <peerid>, Receiver <peerid>
//		publisher = e.Info["Receiver"].(string)
//		receiver = e.Info["Publisher"].(string)
//		if publisher != e.Peer{
//			panic(&UnMatchedSelfPeer{
//				selfpeer: e.Peer,
//				peer:     publisher,
//			})
//		}
//
//	case "TKTRECV":
//		// [TKTRECV] Cid <cid>, Publisher <peerid>, Receiver <peerid>, TimeStamp <time>
//		publisher = e.Info["Publisher"].(string)
//		receiver = e.Info["Receiver"].(string)
//
//	case "TKTREJECT", "TKTACCEPT":
//		// [TKTREJECT] Cid <cid>, Publisher <peerid>, Receiver <peerid>, TimeStamp <time>
//		// [TKTACCEPT] Cid <cid>, Publisher <peerid>, Receiver <peerid>, TimeStamp <time>
//		publisher = e.Info["Receiver"].(string)
//		receiver = e.Info["Publisher"].(string)
//		if publisher != e.Peer{
//			panic(&UnMatchedSelfPeer{
//				selfpeer: e.Peer,
//				peer:     publisher,
//			})
//		}
//
//	case "ACKRECV":
//		// [ACKRECV] Cid <cid>, Publisher <peerid>, Receiver <peerid>, Type <ACCEPT|CANCEL>
//		publisher = e.Info["Receiver"].(string)
//		receiver = e.Info["Publisher"].(string)
//		if receiver != e.Peer{
//			panic(&UnMatchedSelfPeer{
//				selfpeer: e.Peer,
//				peer:     publisher,
//			})
//		}
//	default:
//		panic(&WrongEventType{eventType:e.Type})
//	}
//	publisherInd := a.names.peerIndex[publisher]
//	result = result + a.csvContain(a.names.number, publisherInd, a.names.Get(receiver))
//	return result
//}

func (a *CSVAnalyzer) writeAllCSV(){
	csvPath := path.Join(a.outputDir, "All.csv")
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
	header := a.csvHeaderAll()
	if _, err := w.Write([]byte(header)); err != nil{
		panic(err)
	}

	for e := a.eventList.Front(); e != nil; e = e.Next(){
		event := e.Value.(*BitswapEvent)
		if Contains(a.filter, event.Type) {
			line := a.csvAllLine(event)
			if _, err := w.Write([]byte(line)); err != nil {
				panic(err)
			}
		}else{
			//fmt.Println(event.Type)
			panic("NO WAY")
		}
	}

	if err = w.Flush(); err != nil {
		panic(err)
	}
}

func (a *CSVAnalyzer) writeNamePeer(){
	txtPath := path.Join(a.outputDir, "name_peer.txt")
	fo, err := os.Create(txtPath)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := fo.Close(); err != nil{
			panic(err)
		}
	}()

	w := bufio.NewWriter(fo)

	for name, peerid := range a.names.names{
		line := fmt.Sprintf("%s - %s\n", name, peerid)
		if _, err := w.Write([]byte(line)); err != nil {
			panic(err)
		}else{
			panic("NO WAY")
		}
	}
	if err = w.Flush(); err != nil{
		panic(err)
	}
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

//=======================================

func (a *CSVAnalyzer) AnalyzerRECVTree(){
	// Add sub

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
			}
			a.names.Add(event.GetPeer(event.Direction[0]))
			a.names.Add(event.GetPeer(event.Direction[1]))
			a.eventList.PushBack(event)
			l.PushBack(event)
		}
	}

	// Write to file
	for cid, l := range a.cidMap{
		a.writeRECVTree(a.outputDir, cid, l)
	}
}

//func (a *CSVAnalyzer) buildTree
type recvTreeNode struct{
	successors []*recvTreeNode
	precursor *recvTreeNode
	peerName string
	//duplicateSend int
	duplicateRecv int
}

type recvTree struct {
	cid             string
	root            *recvTreeNode
	//duplicatedSends int
	duplicateRecv int
	nameNode        map[string]*recvTreeNode
}

func newTreeNode(name string) *recvTreeNode{
	return &recvTreeNode{
		successors:    make([]*recvTreeNode,0),
		precursor:     nil,
		peerName:      name,
		duplicateRecv: 0,
	}
}

func (n *recvTreeNode) addSuccessor(s *recvTreeNode){
	n.successors = append(n.successors, s)
}

func (a *CSVAnalyzer) buildRecvTree(cid string, l *list.List) *recvTree{
	result := &recvTree{
		cid:             cid,
		root:            nil,
		//duplicatedSends: 0,
		duplicateRecv: 0,
		nameNode:		make(map[string]*recvTreeNode),
	}

	for e := l.Front(); e != nil; e = e.Next(){
		event := e.Value.(*BitswapEvent)
		switch event.Type{
		case "BLKRECV":
			publisher := a.names.names[event.GetPeer(event.Direction[0])]
			receiver := a.names.names[event.GetPeer(event.Direction[1])]

			// Get receive node
			rNode, ok := result.nameNode[receiver]
			if !ok {
				rNode = newTreeNode(receiver)
				//rNode.duplicateRecv = 1
			}
			if rNode.duplicateRecv > 0{
				rNode.duplicateRecv += 1
				fmt.Println(fmt.Sprintf("Duplicate block receive %s : %s -> %s", cid, publisher, receiver))
				break
			}
			rNode.duplicateRecv = 1

			// Get publish node
			pNode, ok := result.nameNode[publisher]
			if !ok{
				pNode = newTreeNode(publisher)
			}

			pNode.addSuccessor(rNode)
			rNode.precursor = pNode

			// Reset root node
			if result.root == nil || result.root == rNode {
				result.root = pNode
			}
		}
	}

	return result
}

func dfsRECVTree(node *recvTreeNode) int {
	var result = 0
	//result = 0
	for _, n := range node.successors{
		result = result + dfsRECVTree(n)
	}
	return result + 1
}

func (a *CSVAnalyzer) writeRECVTree(outDir string, cid string, l *list.List){
	tree := a.buildRecvTree(cid, l)
	//var counter int
	counted := dfsRECVTree(tree.root)
	fmt.Println(fmt.Sprintf("Counted nodes %d, Total nodes %d", counted, len(tree.nameNode)))
}