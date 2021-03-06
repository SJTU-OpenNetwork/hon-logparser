package analyzer

import (
	"fmt"
	"github.com/SJTU-OpenNetwork/hon-logparser/utils"
	"os"
	"path"
)

type Statistic struct {
	PeerId		string
	NumBlockSend int
	NumBlockRecv int
	NumDupBlock int
}

func NewEmpryStatistic(peerId string) *Statistic {
	return &Statistic{
		PeerId:       peerId,
		NumBlockSend: 0,
		NumBlockRecv: 0,
		NumDupBlock:  0,
	}
}

func CountForFile(parser *Parser, filePath string) (*Statistic, error){
	if parser == nil {
		fmt.Printf("Nil parser !!!!!!\n")
	}
	// Make sure to set a filter that the recorder can check self peer.
	recorder, err := RecorderFromFile(filePath, parser, []string{"BLKSEND", "BLKRECV", "TKTRECV"})
	if err != nil {
		return nil, err
	}
	recorder.CheckSelf()
	receivedBlk := make(map[string]int)
	result := &Statistic{
		PeerId: recorder.selfPeer,
		NumBlockSend: 0,
		NumBlockRecv: 0,
		NumDupBlock: 0,
	}
	for e:= recorder.eventList.Front(); e != nil; e = e.Next(){
		event := e.Value.(*Event)
		switch event.Type{
		case "BLKSEND":
			result.NumBlockSend += 1
		case "BLKRECV":
			result.NumBlockRecv += 1
			cid := event.Info["Cid"].(string)
			_, ok := receivedBlk[cid]
			if !ok {
				receivedBlk[cid] = 1
			} else {
				receivedBlk[cid] = receivedBlk[cid] + 1
				result.NumDupBlock += 1
			}
		default:
			// Do nothing
		}
	}
	return result, nil
}

// SaveToDiskFile would save statistic object as an json file with path outPath.
func (s *Statistic) SaveToDiskFile(outPath string) error {
	outMap := map[string]interface{}{
		"PeerId": s.PeerId,
		"NumBlockSend": s.NumBlockSend,
		"NumBlockRecv":s.NumBlockRecv,
		"NumDupBlock": s.NumDupBlock,
	}
	js := utils.Map2json(outMap)
	_, err := utils.WriteBytes(outPath, []byte(js))
	if err != nil {
		return err
	}
	return nil
}

// SaveToDisk would save statistic object as an json file with path outDir/s.PeerId.json
func (s *Statistic) SaveToDisk(outDir string) error {
	ok, err := utils.PathExists(outDir)
	if err != nil {
		return err
	}
	if !ok {
		err = os.Mkdir(outDir, os.ModePerm); if err != nil {return err}
	}

	filePath := path.Join(outDir, s.PeerId + ".json")
	err = s.SaveToDiskFile(filePath)
	if err != nil {
		return err
	}
	return nil
}

// MergeTwoStatistics merges two statistics into one.
// Note that the new statistic will keep their PeerId if two statistics have the same PeerId.
// Otherwise its PeerId would be "ALL".
func MergeTwoStatistics(s1 *Statistic, s2 *Statistic) *Statistic {
	var tmpPeerId string
	if s1.PeerId == s2.PeerId {
		tmpPeerId = s1.PeerId
	} else if s1.PeerId == "" {
		tmpPeerId = s2.PeerId
	} else if s2.PeerId == ""{
		tmpPeerId = s1.PeerId
	} else {
		tmpPeerId = "ALL"
	}
	return &Statistic{
		PeerId:	  tmpPeerId,
		NumBlockSend: s1.NumBlockSend + s2.NumBlockSend,
		NumBlockRecv: s1.NumBlockRecv + s2.NumBlockRecv,
		NumDupBlock: s1.NumDupBlock + s2.NumDupBlock,
	}
}

/**
 * StatisticStore is used to manage statistics from different peers.
 * One peer may have multi log files. StatisticStore can merge the statistic result from the same peer to one statistic result.
 */
type StatisticStore struct {
	data map[string]*Statistic
}

func NewStatisticStore() *StatisticStore{
	return &StatisticStore{data:make(map[string]*Statistic)}
}

// Note that if statistic contains no peerId info, it would be named as Unknown
func (store *StatisticStore) Add(statistic *Statistic) {
	if statistic.PeerId == "" {
		fmt.Printf("Unknown peerId\n")
		statistic.PeerId = "Unknown"
	}
	_, ok := store.data[statistic.PeerId]
	if !ok {
		store.data[statistic.PeerId] = statistic
	} else {
		store.data[statistic.PeerId] = MergeTwoStatistics(store.data[statistic.PeerId], statistic)
	}
}

func (store *StatisticStore) SaveToDisk(outDir string) error {
	for _, sta := range store.data{
		err := sta.SaveToDisk(outDir); if err != nil {return err}
	}
	return nil
}

