package main

import (
	"container/list"
	"time"
)

type Recorder struct{
	selfPeer     string			// "" means it is a recorder merged from several sub-recorders.
	eventList    *list.List	// Store the events
	eventCounter Counter
}

type BitswapEvent struct{
	Peer string		// Peer is meaningless when it belongs to the recorder of one log files.
					// But it is used to store the event owner when we merge several recorders.
	Type string
	Time time.Time
	Info map[string]interface{}
}

func CreateRecorder() *Recorder{
	return &Recorder{
		selfPeer:  "",
		eventList: list.New(),
		eventCounter: nil,
	}
}

func (r *Recorder) AddEvent(event *BitswapEvent){

}
