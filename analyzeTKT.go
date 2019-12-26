package main

import (
    "bufio"
    "container/list"
    "fmt"
    "path"
)

func (a *CSVAnalyzer) AnalyzeTKT(){
    // Initialize directory
    outDir = path.Join(a.outputDir, "tickets")
    ok, err := PathExists(outDir)
    if err != nil{
        panic(err)
    }
    if !ok {
        err = os.Mkdir(outDir, os.ModePerm)
        if err != nil {
            fmt.Println(err.Error())
        }
    }

    // Traverse events
    blkMap := make(map[string] *list.List)
    tktMap := make(map[string]map[string] *list.List)
    for e := a.recorder.eventList.Front(); e != nil; e = e.Next(){
        event := e.Value.(*BitswapEvent)
        switch event.Type{
        case "BLKRECV":
            cid := event.Info["Cid"].(string)
            l, ok := cidMap[cid]
            if !ok {
                l = list.New()
                blkMap[cid] = l
            }
            a.names.GetandAdd(event.GetPeer(event.Direction[0]))
            a.names.GetandAdd(event.GetPeer(event.Direction[1]))
            l.PushBack(event)
        case "TKTSEND":

        }
    }

    for cid, l := range blkMap{
    
    }

}


// State:
//      Send -> Accept
//           -> Reject
//           -> (UNKNOWN)
// Directiona:
//      Name from Analyzer.names instead of peerId
// Time:
//      Time of event
// 
type simpleTktEvent struct{
    State string
    Direction []string
    
}