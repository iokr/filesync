package task

import (
	"encoding/json"
	"github.com/dzhenquan/filesync/util"
)

type TaskInfo struct {
	FilePort		int 		`json:"filePort"`
	TaskID			string		`json:"taskID"'`
	TaskType		string		`json:"taskType"`
	SrcHost 		string		`json:"srcHost"`
	DestHost		string		`json:"destHost"`
	SrcPath			string		`json:"srcPath"`
	DestPath		string		`json:"destPath"`
	TranType		int			`json:"tranType"`
	ScheduleTime	int64		`json:"scheduleTime"`
}


func (taskInfo *TaskInfo) SendTaskInfoJson(taskJson []byte) bool {
	respJson := make([]byte, util.MAX_MESSAGE_LEN)
	respMsg := &util.RespMessage{}

	//连接到destHost
	otherConn, err := util.CreateSocketConnect(taskInfo.DestHost, util.MSG_TRAN_PORT)
	if err != nil {
		return false
	}
	defer otherConn.Close()

	_, err = otherConn.Write(taskJson)
	if err != nil {
		return false
	}

	recvLen, err := otherConn.Read(respJson)
	if err != nil {
		return false
	}

	err = json.Unmarshal(respJson[:recvLen], respMsg)
	if err != nil {
		return false
	}

	return respMsg.Status
}

func (taskInfo *TaskInfo) FetchTaskInfoFromDB() {

}
