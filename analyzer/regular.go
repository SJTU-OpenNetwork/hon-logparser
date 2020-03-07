package analyzer

/**
 * Define all the regular expression here.
 */
import "github.com/SJTU-OpenNetwork/hon-textile/stream"

var infoExprs = map[string]string {
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
	stream.TAG_WORKERSTART : `Stream ([\w]*), To ([\w]*).*`,
	stream.TAG_WORKEREND : `Stream ([\w]*), To ([\w]*),*`,
	stream.TAG_BLOCKSEND: `Block ([\w]*), Stream ([\w]*), Index ([0-9]*), To ([\w]*), Size ([0-9]*).*`,
	stream.TAG_BLOCKRECEIVE : `Block ([\w]*), Stream ([\w]*), From ([\w]*), Size ([0-9]*).*`,
	stream.TAG_STREAMREQUEST : `Stream ([\w]*), To ([\w]*).*`,
	stream.TAG_STREAMRESPONSE : `Stream ([\w]*), From ([\w]*).*`,
}

// basicExpr is used to
// Format of basic expression:
//		time level module file:location [TYPE] info
//		info can further match infoExprs
var basicExpr = `([\d -\.:]{26}) ([A-Z]*) ([a-z-_\.]*) ([a-z-_:\.0-9A-Z]*) \[([A-Z]*)\] (.*)`
var timeExpr = `([\d -\.:]{26}) ([A-Z]*) ([a-z-_\.]*) ([a-z-_:\.0-9A-Z]*) =====pic_cid:([\w]*) millis:([0-9]*) bytes:([0-9]*) bytePerMills:([0-9]*).*`
var allEventType = []string{"MSGRECV", "BLKRECV", "BLKSEND", "BLKCANCEL", "WANTRECV", "WANTSEND", "TKTRECV", "TKTREJECT", "TKTACCEPT", "TKTSEND", "ACKSEND", "ACKRECV"}