package model

import (
	"github.com/dzhenquan/filesync/util"
	"strconv"
	"fmt"
)

type TaskFileInfo struct {
	BaseModel
	TaskID			string		`json:"taskID"'`
	SrcHost 		string		`json:"srcHost"`
	DestHost		string		`json:"destHost"`
	FilePort		int 		`json:"filePort"`
	SrcPath			string		`json:"srcPath"`
	DestPath		string		`json:"destPath"`
	Status			int			`json:"status"`
	TranType		int			`json:"tranType"`
	ScheduleTime	int64		`json:"scheduleTime"`
	LastFinishTime 	int64 		`json:"lastFinishTime"`
}


func (tFileInfo *TaskFileInfo) Insert() error {
	return DB.Create(tFileInfo).Error
}

func (tFileInfo *TaskFileInfo) Delete() error {
	return DB.Delete(tFileInfo).Error
}

func (tFileInfo *TaskFileInfo) Find() ( *TaskFileInfo,  error) {
	var taskFileInfo TaskFileInfo

	err := DB.Where("task_id = ?", tFileInfo.TaskID).First(&taskFileInfo).Error

	return &taskFileInfo, err
}

func (tFileInfo *TaskFileInfo) Update() error {
	tFileInfo.Status = util.TASK_IS_STOP

	sql := fmt.Sprintf("update task_file_info set src_host='%s',dest_host='%s'," +
		"file_port=%d,src_path='%s',dest_path='%s',status=%d,schedule_time=%d,tran_type=%d " +
			"where task_id='%s';",tFileInfo.SrcHost, tFileInfo.DestHost, tFileInfo.FilePort,
				tFileInfo.SrcPath,tFileInfo.DestPath,tFileInfo.Status,
					tFileInfo.ScheduleTime,tFileInfo.TranType, tFileInfo.TaskID)

	return DB.Exec(sql).Error
}

func (tFileInfo *TaskFileInfo) UpdateTaskStatusTime() error {
	sql := fmt.Sprintf("update task_file_info set last_finish_time=%d, status=%d where task_id='%s';",
		tFileInfo.LastFinishTime,tFileInfo.Status, tFileInfo.TaskID)

	return DB.Exec(sql).Error
}

func (tFileInfo *TaskFileInfo) UpdateTaskStatus() error {

	sql := fmt.Sprintf("update task_file_info set status=%d where task_id='%s';",
		tFileInfo.Status, tFileInfo.TaskID)

	return DB.Exec(sql).Error
}

func (tFileInfo *TaskFileInfo) SendTaskInfoToLocal(taskType string) bool {
	conn, err := util.CreateSocketConnect("127.0.0.1", util.MSG_TRAN_PORT)
	if err != nil {
		return false
	}
	defer conn.Close()

	taskJson := tFileInfo.GetTaskJsonPack(taskType)

	_, err = conn.Write([]byte(taskJson))
	if err != nil {
		return false
	}

	respJson := make([]byte, util.MAX_MESSAGE_LEN)
	recvLen, err := conn.Read(respJson)
	if err != nil {
		return false
	}

	respMsg := util.RespMessage{}

	return respMsg.CheckRespIsTureOrFalse(respJson[:recvLen])
}

func (tFileInfo *TaskFileInfo) GetTaskJsonPack(taskType string) string {

	taskJson := fmt.Sprintf("{\"taskID\":\"%s\",\"taskType\":\"%s\",\"filePort\":%d," +
		"\"srcHost\":\"%s\",\"destHost\":\"%s\",\"srcPath\":\"%s\",\"destPath\":\"%s\"," +
			"\"tranType\":%d,\"scheduleTime\":%d}",
				tFileInfo.TaskID, taskType, tFileInfo.FilePort,tFileInfo.SrcHost,
					tFileInfo.DestHost,tFileInfo.SrcPath,tFileInfo.DestPath,
						tFileInfo.TranType,tFileInfo.ScheduleTime)
	return taskJson
}

func FindTaskByID(id string) (*TaskFileInfo, error) {
	var tFileInfo TaskFileInfo

	tid, _ := strconv.Atoi(id)

	err := DB.First(&tFileInfo, "id = ?", tid).Error

	return &tFileInfo, err
}

func FindAllTaskQuery() ([]*TaskFileInfo, error) {
	rows, err := DB.Raw("select * from task_file_info;").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taskFiles []*TaskFileInfo
	for rows.Next() {
		var taskFile TaskFileInfo
		err := DB.ScanRows(rows, &taskFile)
		if err != nil {
			return nil, err
		}
		taskFiles = append(taskFiles, &taskFile)
	}
	return taskFiles, nil
}