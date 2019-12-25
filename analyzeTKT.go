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
            l, ok := cidMap[]
        }
    }

}


// State:
//      Send -> Accept
//           -> Reject
//           -> (UNKNOWN)
// Direction:
//      Name from Analyzer.names instead of peerId
type simpleTktEvent struct{
    State string
    Direction []string

}
