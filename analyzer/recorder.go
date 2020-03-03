package analyzer

import (
	"bufio"
	"container/list"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

/**
 * Build an recorder from a log file.
 * 	- filePath: path of file.
 *	- filter: recorder will only contain events with type in filter
 */
func RecorderFromFile(filePath string, parser *Parser, filter []string) (*Recorder, error) {
	fmt.Printf("Create filter for recorder\n")
	//fmt.Printf(filter)
	mapFilter := make(map[string]interface{})
	if filter == nil || len(filter)==0 {
		filter = allEventType
	}
	for _, str := range filter {
		mapFilter[str] = struct{}{}
	}

	fmt.Printf("Initialize recorder\n")
	rec := CreateRecorder()
	rec.AddMapCounter()
	rec.SetFilter(filter)

	fmt.Printf("Begin parse for %s\n", filePath)
	f, err := os.Open(filePath)
	if err != nil{
		fmt.Printf("Cannot open %s\n", filePath)
		return nil, err
	}
	fmt.Printf("File opened\n")
	defer f.Close()

	reader := bufio.NewReader(f)
	lineNum := 0
	for {
		line, _, err := reader.ReadLine()
		fmt.Printf(string(line))
		if err == io.EOF {
			fmt.Println("Read to EOF\n")
			break
		}
		fmt.Printf("Try to parse line\n")
		event, err := parser.ParseLineWithFilter(string(line), mapFilter)
		if err == nil  && event != nil{	//Note that event would be nil if the line can not be parsed.
			rec.AddEvent(event)
		}else{
			fmt.Printf("Error: %s\n", err.Error())
		}
		lineNum++
		if lineNum % 10 == 0 {
			fmt.Printf("line: %d\r", lineNum)
		}
	}

	return rec, nil
}

type Recorder struct{
	eventFilter  []string		// filter within recorder have nothing to do with how to parse a file.
								// It just contains info about the event type in a recorder.
	selfPeer     string			// "" means it is a recorder merged from several sub-recorders.
	eventList    *list.List	// Store the events
	// We use list instead of a simple slice
	eventCounter Counter
}

type Event struct{
	Peer string		// Peer is meaningless when it belongs to the recorder of one log files.
	// But it is used to store the event owner when we merge several recorders.
	Type string
	Time time.Time
	Direction []string
	Info map[string]interface{}
}

func (e *Event)GetPeer(peer string) string{
	if peer == SELF && e.Peer!=""{
		return e.Peer
	} else {
		return peer
	}
}

/**
 * When we create
 */
//func copyEvent(event *Event) *Event{
//
//}

func CreateRecorder() *Recorder{
	return &Recorder{
		selfPeer:  "",
		eventList: list.New(),
		eventCounter: nil,
	}
}

func (r *Recorder) AddMapCounter(){
	if r.eventCounter != nil {
		fmt.Println("Recorder already have a counter")
	}else {

		r.eventCounter = CreateMapCounter()
	}
}

func (r *Recorder) SetFilter(filter []string) {
	r.eventFilter = filter
}

func (r *Recorder) AddEvent(event *Event){
	r.eventList.PushBack(event)
	if r.eventCounter != nil {
		r.eventCounter.Count(event)
	}
}

func (r *Recorder) PrintCounter(){
	if r.eventCounter == nil {
		fmt.Println("Recorder do not have a counter.")
	} else {
		fmt.Println(r.eventCounter.String())
	}
}

func (r *Recorder) SaveCounter(outPath string){
	if r.eventCounter == nil {
		fmt.Println("Recorder do not have a counter.")
	} else {
		fo, err := os.Create(outPath)
		if err != nil {
			panic(err)
		}

		defer func(){
			if err := fo.Close(); err != nil {
				panic(err)
			}
		}()

		w:= bufio.NewWriter(fo)

		w.Write([]byte(r.eventCounter.String()))
		if err = w.Flush(); err != nil {
			panic(err)
		}
	}
}

func (r *Recorder) checkPerSelf(peer string){
	if r.selfPeer == "" {
		r.selfPeer = peer
	} else {
		if r.selfPeer != peer{
			fmt.Println("Un matched self peer "+r.selfPeer + " " + peer)
			r.selfPeer = peer
		}
	}
}

// CheckSelf check and set the value of selfpeer
func (r *Recorder) CheckSelf(){
	for e:= r.eventList.Front(); e != nil; e = e.Next(){
		event := e.Value.(*Event)
		switch event.Type{
		case "ACKSEND":
			// [ACKSEND] Cid <cid>, Publisher <peerid>, Receiver <peerid>
			r.checkPerSelf(event.Info["Receiver"].(string))
		case "TKTRECV":
			// [TKTRECV] Cid <cid>, Publisher <peerid>, Receiver <peerid>, TimeStamp <time>
			r.checkPerSelf(event.Info["Receiver"].(string))

		case "TKTREJECT", "TKTACCEPT":
			// [TKTREJECT] Cid <cid>, Publisher <peerid>, Receiver <peerid>, TimeStamp <time>
			// [TKTACCEPT] Cid <cid>, Publisher <peerid>, Receiver <peerid>, TimeStamp <time>
			r.checkPerSelf(event.Info["Receiver"].(string))

		case "ACKRECV":
			// [ACKRECV] Cid <cid>, Publisher <peerid>, Receiver <peerid>, Type <ACCEPT|CANCEL>
			r.checkPerSelf(event.Info["Publisher"].(string))

		}
	}
}

func (r *Recorder) SetEventsPeer() bool {
	if r.selfPeer == ""{
		fmt.Println("Recorder contains no self peer info!!")
		return false
	}else {
		for e := r.eventList.Front(); e != nil; e = e.Next() {
			e.Value.(*Event).Peer = r.selfPeer
		}
		return true
	}
}

func MergeRecorders(rs []*Recorder) *Recorder{
	rch := make(chan *Recorder, len(rs))
	//resultChan := make(chan *Recorder)	//channel used to hold output
	//rch2 := make(chan *Recorder, len(rs))
	for _,r := range rs{
		rch <- r
	}


	//There should be at least one recorder in channel
	for {
		//fmt.Println(len(rch))
		r1 := <-rch
		select {
		case r2 := <-rch:
			rch <- mergeTwoRecorders(r1, r2)
		default:
			return r1
		}
	}
}

func mergeTwoRecorders(r1 *Recorder, r2 *Recorder) *Recorder {
	//ech1 := make(chan *Event, r1.eventList.Len())
	//ech2 := make(chan *Event, r2.eventList.Len())
	result := CreateRecorder()
	e1 := r1.eventList.Front()
	e2 := r2.eventList.Front()
	for{
		if e1 == nil && e2 == nil{
			break
		}

		if e1 == nil{
			for ; e2 != nil; e2 = e2.Next(){
				result.eventList.PushBack(e2.Value)
			}
		} else if e2 ==nil {
			for ; e1 != nil; e1 = e1.Next(){
				result.eventList.PushBack(e1.Value)
			}
		} else {
			if e1.Value.(*Event).Time.Before(e2.Value.(*Event).Time){
				result.AddEvent(e1.Value.(*Event))
				e1 = e1.Next()
			} else {
				result.AddEvent(e2.Value.(*Event))
				e2 = e2.Next()
			}
		}
	}
	return result
}

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

