package main

import (
	"fmt"
	"regexp"
)

/**
 * Formatted Tag:
 * [MSGRECV] From <peerid>
 * [BLKRECV] Cid <cid>, From <peerid>
 * [BLKSEND] Cid <cid>, SendTo <peerid>
 * [BLKCANCEL] Cid <cid>, From <peerid>
 * [WANTRECV] Cid <cid>, From <peerid>
 * [WANTSEND] Cid <cid>, SendTo <peerid>	;peerid could be ALL if it is broadcast
 * [TKTRECV] Cid <cid>, Publisher <peerid>, Receiver <peerid>, TimeStamp <time>
 * [TKTREJECT] Cid <cid>, Publisher <peerid>, Receiver <peerid>, TimeStamp <time>
 * [TKTACCEPT] Cid <cid>, Publisher <peerid>, Receiver <peerid>, TimeStamp <time>
 * [TKTSEND] Cid <cid>, SendTo <peerid>, TimeStamp <time>
 * [ACKSEND] Cid <cid>, Publisher <peerid>, Receiver <peerid>	; The receiver of ack means receiver of the corresponding ticket.
 *																; In other words, it is the sender of ticket acks.
 *																; It is set in this way as receiver and publisher are used to index to specific ticket.
 * [ACKRECV] Cid <cid>, Publisher <peerid>, Receiver <peerid>, Type <ACCEPT|CANCEL>
 */

var (
	infoExprs = map[string]string {
		"MSGRECV": `From ([\w]*).*`,
		"BLKRECV": `Cid ([\w]*), From ([\w]*).*`,
		"BLKSEND": `Cid ([\w]*), SendTo ([\w]*).*`,
		"BLKCANCEL":`Cid ([\w]*), From ([\w]*).*`,
		"WANTRECV" : `Cid ([\w]*), From ([\w]*).*`,
		"WANTSEND" : `Cid ([\w]*), SendTo ([\w]*).*`,
		//"TKTRECV"  : `Cid ([\w]*), Publisher ([\w]*), Receiver ([\w]*), TimeStamp ([0-9]*).*`,
		//"TKTREJECT": `Cid ([\w]*), Publisher ([\w]*), Receiver ([\w]*), TimeStamp ([0-9]*).*`,
		//"TKTACCEPT": `Cid ([\w]*), Publisher ([\w]*), Receiver ([\w]*), TimeStamp ([0-9]*).*`,
		//"TKTSEND"  : `Cid ([\w]*), SendTo ([\w]*), TimeStamp ([0-9]*).*`,
		"TKTRECV"  : `Cid ([\w]*), Publisher ([\w]*), Receiver ([\w]*).*`,
		"TKTREJECT": `Cid ([\w]*), Publisher ([\w]*), Receiver ([\w]*).*`,
		"TKTACCEPT": `Cid ([\w]*), Publisher ([\w]*), Receiver ([\w]*).*`,
		"TKTSEND"  : `Cid ([\w]*), SendTo ([\w]*).*`,

		"ACKSEND"  : `Cid ([\w]*), Publisher ([\w]*), Receiver ([\w]*).*`,
		"ACKRECV"  : `Cid ([\w]*), Publisher ([\w]*), Receiver ([\w]*), Type ([A-Z]*).*`,
	}
	infoRegs  map[string]*regexp.Regexp
	basicExpr = `([\d -\.:]{26}) ([A-Z]*) ([a-z-_\.]*) ([a-z-_:\.0-9A-Z]*) \[([A-Z]*)\] (.*)`
	basicReg  *regexp.Regexp
	err       error
)
const timeFotmat = "2006-01-02 15:04:05.000000"
const SELF = "SELF"

func init(){
	fmt.Println("Initialize regulation expressions.")
	infoRegs = make(map[string]*regexp.Regexp)
	basicReg, err = regexp.Compile(basicExpr)
	for k, exp := range infoExprs {
		infoRegs[k], err = regexp.Compile(exp)
		if err != nil{
			infoRegs = nil
			fmt.Println(err.Error())
			break
		}
	}
}

func extractBasic(line string) (map[string]string, error) {
	//fmt.Println(line)
	params := basicReg.FindStringSubmatch(line)
	// original log; time; type; subsystem; location; contains
	//fmt.Println(params)
	//fmt.Println(len(params))
	if len(params) > 6 {
		return map[string] string{
			"origin": params[0],
			"time": params[1],
			"type": params[2],
			"system": params[3],
			"location": params[4],
			"event": params[5],
			"eventInfo": params[6],
		}, nil
	} else {
		return nil, &InvalidLogLine{}
	}
}

func extractInfo(event string, info string) ([]string, error) {
	reg, ok := infoRegs[event]
	if !ok {
		return nil, &UnknownReg{reg:event}
	}
	params := reg.FindStringSubmatch(info)
	if len(params) > 0 {
		return params, nil
	}
	return nil, &ParseFailed{
		expr: infoExprs[event],
		str:  info,
	}
}

func parseInfo(info map[string]string) (*BitswapEvent, error) {

	tmpTime, _ := parseTimestamp(info["time"])
	params, err := extractInfo(info["event"], info["eventInfo"])
	if err != nil {
		return nil, err
	}

	switch info["event"]{
	case "MSGRECV":
		// [MSGRECV] From <peerid>
		return &BitswapEvent{
			Type: info["event"],
			Time: tmpTime,
			Direction: []string{
				params[1], SELF,
			},
			Info: map[string]interface{}{
				"From": params[1],
			},
		}, nil
	case "BLKRECV", "BLKCANCEL", "WANTRECV":
		// [BLKRECV] Cid <cid>, From <peerid>
		// [BLKCANCEL] Cid <cid>, From <peerid>
		// [WANTRECV] Cid <cid>, From <peerid>
		return &BitswapEvent{
			Type: info["event"],
			Time: tmpTime,
			Direction:[]string{
				params[2], SELF,
			},
			Info: map[string]interface{}{
				"Cid": params[1],
				"From": params[2],
			},
		}, nil
	case "BLKSEND", "WANTSEND":
		// [BLKSEND] Cid <cid>, SendTo <peerid>
		// [WANTSEND] Cid <cid>, SendTo <peerid>
		return &BitswapEvent{
			Type: info["event"],
			Time: tmpTime,
			Direction: []string{
				SELF, params[2],
			},
			Info: map[string]interface{}{
				"Cid": params[1],
				"SendTo": params[2],
			},
		},nil
	case "TKTSEND":
		// [TKTSEND] Cid <cid>, SendTo <peerid>, TimeStamp <time>
		return &BitswapEvent{
			Type: info["event"],
			Time: tmpTime,
			Direction: []string{
				SELF, params[2],
			},
			Info: map[string]interface{}{
				"Cid": params[1],
				"SendTo": params[2],
				//"TimeStamp": params[3],
			},
		},nil
	case "ACKSEND":
		// [ACKSEND] Cid <cid>, Publisher <peerid>, Receiver <peerid>
		return &BitswapEvent{
			Type: info["event"],
			Time: tmpTime,
			Direction: []string{
				SELF, params[2],
			},
			Info: map[string]interface{}{
				"Cid": params[1],
				"Publisher": params[2],
				"Receiver": params[3],
			},
		}, nil
	case "TKTRECV":
		// [TKTRECV] Cid <cid>, Publisher <peerid>, Receiver <peerid>, TimeStamp <time>
		return &BitswapEvent{
			Type: info["event"],
			Time: tmpTime,
			Direction: []string{
				params[2], SELF,
			},
			Info: map[string]interface{}{
				"Cid": params[1],
				"Publisher": params[2],
				"Receiver": params[3],
				//"TimeStamp": params[4],
			},
		}, nil
	case "TKTREJECT", "TKTACCEPT":
		// [TKTREJECT] Cid <cid>, Publisher <peerid>, Receiver <peerid>, TimeStamp <time>
		// [TKTACCEPT] Cid <cid>, Publisher <peerid>, Receiver <peerid>, TimeStamp <time>
		return &BitswapEvent{
			Type: info["event"],
			Time: tmpTime,
			Direction: []string{
				SELF, params[2],
			},
			Info: map[string]interface{}{
				"Cid": params[1],
				"Publisher": params[2],
				"Receiver": params[3],
				//"TimeStamp": params[4],
			},
		}, nil
	case "ACKRECV":
		// [ACKRECV] Cid <cid>, Publisher <peerid>, Receiver <peerid>, Type <ACCEPT|CANCEL>
		return &BitswapEvent{
			Type: info["event"],
			Time: tmpTime,
			Direction: []string{
				params[3], SELF,
			},
			Info: map[string]interface{}{
				"Cid" : params[1],
				"Publisher": params[2],
				"Receiver" : params[3],
				"Type" : params[4],
			},
		}, nil
	}

	return nil, nil
}


func testParse(info map[string]string){
	//fmt.Println(Stringmap2json(info))
	event, err := parseInfo(info)
	if err != nil{
		fmt.Println(err.Error())
		return
	}
	if event == nil {
		return
	}
	//fmt.Println(*event)
	/*
	switch info["system"] {
	case "tex-ipfs": fmt.Println(info["contains"])
	}
	 */
}

func ParseLine(line string) (*BitswapEvent, error) {
	info, err := extractBasic(line)
	if err != nil{
		return nil, err
	}

	event, err := parseInfo(info)
	if err != nil{
		return nil, err
	}

	return event, nil
}

