package analyzer

import (
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-logparser/utils"
	"regexp"
	"time"

	//"github.com/SJTU-OpenNetwork/hon-textile/stream"
)

const TimeFotmat = "2006-01-02 15:04:05.000000"
const SELF = "SELF"

/**
 * Parser is used to parse the log line that match the basicExpr
 * Use Parser.ParseLine() to fetch an event from one line of log.
 */
type Parser struct {
	infoRegs  map[string]*regexp.Regexp
	basicReg  *regexp.Regexp
	timeReg   *regexp.Regexp
}

func NewParser() (*Parser, error) {
	fmt.Println("Initialize regulation expressions.")
	infoRegs := make(map[string]*regexp.Regexp)
	basicReg, err := regexp.Compile(basicExpr); if err != nil {return nil, err}
	for k, exp := range infoExprs {
		infoRegs[k], err = regexp.Compile(exp); if err != nil {return nil, err}
	}
	timeReg, err := regexp.Compile(timeExpr); if err != nil {return nil, err}
	return &Parser{
		infoRegs:infoRegs,
		basicReg:basicReg,
		timeReg:timeReg,
	}, nil
}

func parseTimestamp(str string) (time.Time, error){
	t, err := time.Parse(TimeFotmat, str)
	if err != nil {
		//fmt.Println(err.Error())
		return t, err
	}
	return t, nil
}

func (parser *Parser) ParseLineWithFilter(line string, filter map[string]interface{}) (*Event, error) {
	//fmt.Printf("extractBasic\n")
	info, err := parser.extractBasic(line)
	//fmt.Printf(filter)
	//fmt.Printf("basic extracted\n")
	if err != nil{
		return nil, err
	}
	if info == nil {	// mismatch basicReg
		return nil, nil
	}

	//fmt.Printf("Check filter\n")
	_, ok := filter[info["event"]]
	if !ok {
		//fmt.Printf("No such filter\n")
		return nil, nil
	}
	//fmt.Printf("Find filter\n")

	event, err := parser.extractInfo(info)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (parser *Parser) ParseLine(line string) (*Event, error){
	info, err := parser.extractBasic(line)
	if err != nil{
		return nil, err
	}

	if info == nil {	// mismatch basicReg
		return nil, nil
	}

	event, err := parser.extractInfo(info)
	if err != nil{
		return nil, err
	}

	return event, nil
}

/**
 * extractBasic parse a line according to basicReg
 */
func (parser *Parser)extractBasic(line string) (map[string]string, error) {
	params := parser.basicReg.FindStringSubmatch(line)
	// original log; time; type; subsystem; location; contains
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
		return nil, nil		// Do not raise error. As it is common the line mismatches basicReg.
	}
}

func (parser *Parser) parseInfo(event string, info string) ([]string, error) {
	reg, ok := parser.infoRegs[event]
	if !ok {
		return nil, &utils.UnknownReg{Reg:event}
	}
	params := reg.FindStringSubmatch(info)
	if len(params) > 0 {
		return params, nil
	}
	return nil, &utils.ParseFailed{
		Expr: infoExprs[event],
		Str:  info,
	}
}


func (parser *Parser)extractInfo(info map[string]string) (*Event, error) {

	tmpTime, _ := parseTimestamp(info["time"])
	params, err := parser.parseInfo(info["event"], info["eventInfo"])
	if err != nil {
		return nil, err
	}

	switch info["event"]{
	case "MSGRECV":
		// [MSGRECV] From <peerid>
		return &Event{
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
		return &Event{
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
		return &Event{
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
		return &Event{
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
		return &Event{
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
		return &Event{
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
		return &Event{
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
		return &Event{
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