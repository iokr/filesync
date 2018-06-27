package task

import (
	"encoding/json"
	"github.com/dzhenquan/filesync/util"
	"github.com/dzhenquan/filesync/config"
)

type TaskInfo struct {
	FilePort		int 		`json:"filePort"`		// 文件传输端口
	TaskID			string		`json:"taskID"'`		// 任务ID
	TaskType		string		`json:"taskType"`		// 任务类型(创建,开始,停止,更新,删除)
	SrcHost 		string		`json:"srcHost"`		// 源主机
	DestHost		string		`json:"destHost"`		// 目标主机
	SrcPath			string		`json:"srcPath"`		// 源路径
	DestPath		string		`json:"destPath"`		// 目的路径
	TranType		int			`json:"tranType"`		// 传输方式(0-复制, 1-移动)
	ScheduleTime	int64		`json:"scheduleTime"`	// 调度时间
}


func (taskInfo *TaskInfo) SendTaskInfoJson(taskJson []byte) bool {
	respJson := make([]byte, util.MAX_MESSAGE_LEN)
	respMsg := &util.RespMessage{}

	//连接到destHost
	otherConn, err := util.CreateSocketConnect(taskInfo.DestHost, config.ServerConfig.FServerPort)
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
