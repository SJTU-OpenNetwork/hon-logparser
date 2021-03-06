package analyzer

import (
	//"encoding/json"
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-logparser/utils"
	"regexp"
	"strconv"
	"time"

	"github.com/SJTU-OpenNetwork/hon-textile/stream"
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
	cidFilter *utils.CidFilter
}

func NewParser() (*Parser, error) {
	fmt.Println("Initialize regulation expressions.\n")
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
		cidFilter:nil,
	}, nil
}

func (parser *Parser) SetCidFilter(filter *utils.CidFilter) {
	parser.cidFilter = filter
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
	info, err := parser.extractBasic(line)
	if err != nil{
		return nil, err
	}
	if info == nil {	// mismatch basicReg
		return nil, nil
	}

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

func (parser *Parser) ParseLineForTime(line string) *TimeInfo {
	// var timeExpr = `([\d -\.:]{26}) ([A-Z]*) ([a-z-_\.]*) ([a-z-_:\.0-9A-Z]*)
	// =====pic_cid:([\w]*) millis:([0-9]*) bytes:([0-9]*) bytePerMills:([0-9]*).*`
	params := parser.timeReg.FindStringSubmatch(line)
	if len(params) > 6 {
		res := &TimeInfo{PeerId: ""}
		res.Cid = params[5]
		res.Ms, _ = strconv.Atoi(params[6])
		res.Bytes, _ = strconv.Atoi(params[7])
		res.BytePerMs, _ = strconv.Atoi(params[8])
		return res
		//return map[string] string{
		//	"origin": params[0],
		//	"time": params[1],
		//	"type": params[2],
		//	"system": params[3],
		//	"location": params[4],
		//	"cid": params[5],
		//	"ms": params[6],
		//	"bytes": params[7],
		//	"bytePerMs": params[8],
		//}
	} else {
		return nil		// Do not raise error. As it is common the line mismatches basicReg.
	}
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

// Core function for parser.
// Note that this func may return nil, nil
func (parser *Parser)extractInfo(info map[string]string) (*Event, error) {

	tmpTime, _ := parseTimestamp(info["time"])
	params, err := parser.parseInfo(info["event"], info["eventInfo"])
	if err != nil {
		return nil, err
	}

	switch info["event"]{
	case "MSGRECV":
		// [MSGRECV] From <peerid>
		//if parser.cidFilter != nil && !parser.cidFilter.Has(params[1]){
		//	return nil, nil
		//}
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
		if parser.cidFilter != nil && !parser.cidFilter.Has(params[1]){
			return nil, nil
		}
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
		if parser.cidFilter != nil && !parser.cidFilter.Has(params[1]){
			return nil, nil
		}
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
		if parser.cidFilter != nil && !parser.cidFilter.Has(params[1]){
			return nil, nil
		}
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
		if parser.cidFilter != nil && !parser.cidFilter.Has(params[1]){
			return nil, nil
		}
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
		if parser.cidFilter != nil && !parser.cidFilter.Has(params[1]){
			return nil, nil
		}
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
		if parser.cidFilter != nil && !parser.cidFilter.Has(params[1]){
			return nil, nil
		}
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
		if parser.cidFilter != nil && !parser.cidFilter.Has(params[1]){
			return nil, nil
		}
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
	case stream.TAG_BLOCKSEND:
		//`Block ([\w]*), Stream ([\w]*), Index ([0-9]*), To ([\w]*), Size ([0-9]*).*`,
		if parser.cidFilter != nil && !parser.cidFilter.Has(params[1]){
			return nil, nil
		}
		return &Event{
			Type: info["event"],
			Time: tmpTime,
			Direction:[]string{
				SELF, params[4],
			},
			Info: map[string]interface{}{
				"Cid": params[1],
				"StreamId": params[2],
				"Index": params[3],
				"SendTo": params[4],
				"Size": params[5],
			},
		}, nil

	case stream.TAG_BLOCKRECEIVE:
		// `Block ([\w]*), Stream ([\w]*), From ([\w]*), Size ([0-9]*).*`
		if parser.cidFilter != nil && !parser.cidFilter.Has(params[1]){
			return nil, nil
		}
		return &Event{
			Type: info["event"],
			Time: tmpTime,
			Direction: []string{
				params[3], SELF,
			},
			Info: map[string]interface{}{
				"Cid": params[1],
				"StreamId": params[2],
				"From": params[3],
				"Size": params[4],
			},
		}, nil
	default:
		//Do nothing
	}

	return nil, nil
}
