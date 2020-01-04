package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"path"
	"time"
)

type CSVAnalyzer struct{
	outputDir string
	recorder *Recorder
	names *peerName					//handle peer rename
	eventList *list.List			//eventList that store only the filtered events
	//cidMap map[string] *list.List	//event indexed by cid
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
		// cidMap : make(map[string] *list.List),
		filter: filter,
	}
}



func (p *peerName) Add(peerId string){
	p.add(peerId)
}

func (p *peerName) GetandAdd(peerId string) string{
    p.add(peerId)
    return p.names[peerId]
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
//Analyze BLKRECV
//=======================================

func (a *CSVAnalyzer) AnalyzerRECVTree(){
	// Build directory
	outDir :=  path.Join(a.outputDir, "trees")
	ok, err := PathExists(outDir)
	if err != nil {
		panic(err)
	}
	if !ok{
		err = os.Mkdir(outDir, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	// Analyze block
	cidMap := make(map[string] *list.List)

	for e := a.recorder.eventList.Front(); e != nil; e = e.Next(){
		switch e.Value.(*BitswapEvent).Type{
		case "BLKRECV":
			// [BLKRECV] Cid <cid>, From <peerid>
			event := e.Value.(*BitswapEvent)
			//fmt.Println(Map2json(event.Info))
			cid := event.Info["Cid"].(string)
			l, ok := cidMap[cid]
			if !ok {
				l = list.New()
				cidMap[cid] = l
			}
			a.names.Add(event.GetPeer(event.Direction[0]))
			a.names.Add(event.GetPeer(event.Direction[1]))
			//a.eventList.PushBack(event)
			l.PushBack(event)
		}
	}

	// Write to file
	for cid, l := range cidMap{
		a.writeRECVTree(outDir, cid, l)
	}
}

//func (a *CSVAnalyzer) buildTree
type recvTreeNode struct{
	successors []*recvTreeNode
	precursor *recvTreeNode
	peerName string
	//duplicateSend int
	duplicateRecv int
	time time.Time
}

type recvTree struct {
	cid             string
	root            *recvTreeNode
	//duplicatedSends int
	duplicateRecv int
	nameNode        map[string]*recvTreeNode
	//nodeList 	*list.List
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

func (a *CSVAnalyzer) buildRecvTree(cid string, l *list.List) (*recvTree, []*BitswapEvent){
	result := &recvTree{
		cid:             cid,
		root:            nil,
		//duplicatedSends: 0,
		duplicateRecv: 0,
		nameNode:		make(map[string]*recvTreeNode),
	}

	//duplicateCounter := make([][]string, 0)
	dupEvent := make([]*BitswapEvent,0)

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
				result.nameNode[receiver] = rNode
				//rNode.duplicateRecv = 1
			}
			if rNode.duplicateRecv > 0{
				rNode.duplicateRecv += 1
				fmt.Println(fmt.Sprintf("Duplicate block receive %s : %s -> %s", cid, publisher, receiver))
				//tmpdup := []string{publisher, receiver,}
				//duplicateCounter = append(duplicateCounter, tmpdup)
				dupEvent = append(dupEvent, event)
				break
			}
			rNode.duplicateRecv = 1
			rNode.time = event.Time		// Only set time when receive event.


			// Get publish node
			pNode, ok := result.nameNode[publisher]
			if !ok{
				pNode = newTreeNode(publisher)
				result.nameNode[publisher] = pNode
			}

			pNode.addSuccessor(rNode)
			rNode.precursor = pNode

			// Reset root node
			if result.root == nil || result.root == rNode {
				result.root = pNode
			}
		}
	}

	return result, dupEvent
}

func dfsRECVTreePrefix(tree *recvTree, node *recvTreeNode, prefix string, isLast bool, isFirst bool, writer *bufio.Writer) int {
	var result = 0
	delete(tree.nameNode, node.peerName)

	millTimeFormat := "15:04:05.000"
	//if node.time != nil {
	tmpTime := node.time.Format(millTimeFormat)
	//}
	if isFirst {
		writer.Write([]byte(fmt.Sprintf("--[%s%3s]", tmpTime, node.peerName)))
	}else {
		writer.Write([]byte(fmt.Sprintf("%s|-[%s%3s]", prefix, tmpTime, node.peerName)))
	}
	if len(node.successors) == 0{
		writer.Write([]byte("\n"))
		return 1
	}
	var nextPrefix string
	if isLast {
		nextPrefix = prefix + "                   "
		//fmt.Println(nextPrefix)
	}else{
		nextPrefix = prefix + "|                  "
	}
	for i, n := range node.successors{
		nextIsLast := (i == len(node.successors)-1)
		nextIsFirst := i == 0
		//writer.Write([]byte(prefix))
		if i==0 {
			//writer.Write([]byte("--"))
			result = result + dfsRECVTreePrefix(tree, n, nextPrefix, nextIsLast, nextIsFirst, writer)

		}else{

			//writer.Write([]byte("|-"))
			result = result + dfsRECVTreePrefix(tree, n, nextPrefix, nextIsLast, nextIsFirst, writer)
		}
	}
	return result + 1
}

func getRoot(node *recvTreeNode) *recvTreeNode{
	for node.precursor != nil {
		node = node.precursor
	}
	return node
}

func getAnyNode(nameNode map[string] *recvTreeNode) *recvTreeNode{
	if len(nameNode) == 0{
		return nil
	}
	for _, node := range nameNode{
		return node
	}
	return nil
}

func (a *CSVAnalyzer) writeRECVTree(outDir string, cid string, l *list.List){
	millTimeFormat := "15:04:05.000"
	tree, dupEvent := a.buildRecvTree(cid, l)
	//var counter int

	txtPath := path.Join(outDir, fmt.Sprintf("%s.txt", cid))
	fo, err := os.Create(txtPath)
	if err != nil {
		panic(err)
	}

	defer func(){
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	w:= bufio.NewWriter(fo)

	//counted := dfsRECVTree(tree.root, 0, w)
	var counted int
	//total := len(tree.nameNode)
	for len(tree.nameNode) > 0 {
		node := getAnyNode(tree.nameNode)
		node = getRoot(node)
		counted = counted + dfsRECVTreePrefix(tree, node, "", true, true, w)
	}
	//fmt.Println(fmt.Sprintf("Counted nodes %d, Total nodes %d", counted, total))
	//w.Write([]byte("Duplicated blk recv\n"))
	for _, dup := range dupEvent{
		w.Write([]byte(fmt.Sprintf("%s %3s => %3s\n",
			dup.Time.Format(millTimeFormat), a.names.names[dup.GetPeer(dup.Direction[0])], a.names.names[dup.GetPeer(dup.Direction[1])])))
	}

	if err = w.Flush(); err != nil {
		panic(err)
	}
}

