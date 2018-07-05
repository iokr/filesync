package main

import (
	"os"
	"io"
	"net"
	"time"
	"strings"
	"encoding/json"
	. "github.com/dzhenquan/filesync/tlog"
	"github.com/dzhenquan/filesync/util"
	"github.com/dzhenquan/filesync/task"
	"github.com/dzhenquan/filesync/model"
	"github.com/dzhenquan/filesync/config"
	"github.com/dzhenquan/filesync/fserver/schedule"
)

func main() {

	//Init Model
	db, err := model.InitDB()
	if err != nil {
		Tlog.Errorln(err.Error())
		os.Exit(-1)
	}
	defer db.Close()

	// 开启多核CPU传输文件
	//numCpu := runtime.NumCPU()
	//runtime.GOMAXPROCS(numCpu)

	listen, err := util.CreateSocketListen(config.ServerConfig.FServerHost,
		config.ServerConfig.FServerPort)
	if err != nil {
		Tlog.Errorln("创建本地监听失败!")
		os.Exit(-1)
	}
	defer listen.Close()

	//从数据库获取任务,对未完成任务进行调度
	go schedule.HandleTaskSchedule()

	Tlog.Println("Message Listen ...")

	for {
		conn, err := listen.Accept()
		if err != nil {
			Tlog.Errorf("接受新的连接请求失败, err: ", err)
			continue
		}
		//log.Println("conn: ", conn.RemoteAddr())
		go handleMessageConn(conn)
	}
}

// handle client send json package
func handleMessageConn(conn net.Conn) {
	respMsg := &util.RespMessage{
		Status: true,
	}

	recvBuf := make([]byte, util.MAX_MESSAGE_LEN)

	readLen, err := conn.Read(recvBuf)
	if err != nil {
		if err == io.EOF {
			return
		}
		Tlog.Errorln("Read client message is failure!")
		return
	}
	defer conn.Close()

	//处理任务请求
	respMsg.Status = handleTaskRequest(recvBuf[:readLen], respMsg)

	//获取返回消息JSON报文
	respJson := respMsg.GetRespMessageJson()
	_, err = conn.Write([]byte(respJson))
	if err != nil {
		Tlog.Println("向客户端回复返回消息失败!")
	}
}

// handle client task request
// returns true/false
func handleTaskRequest(taskJson []byte, respMsg *util.RespMessage) bool {
	var returnValue bool
	var fileTask 	*task.FileTask

	taskInfo := &task.TaskInfo{}
	tFileInfo := &model.TaskFileInfo{}

	// 解析json报文
	err := json.Unmarshal(taskJson, taskInfo)
	if err != nil {
		respMsg.TaskType = util.TASK_UNABLE
		return false
	}

	respMsg.TaskType = taskInfo.TaskType

	isLocalDestIP := true
	isLocalHost := true

	if strings.Compare(taskInfo.SrcHost, taskInfo.DestHost) != 0 {
		isLocalHost = false

		if !util.CheckIPIsLocalIP(taskInfo.DestHost) {
			isLocalDestIP = false

			if (!util.CheckIsDirByPath(taskInfo.SrcPath)) {
				Tlog.Printf("源文件路径[%s]不存在!", taskInfo.SrcPath)
				return false
			}
			if !taskInfo.SendTaskInfoJson(taskJson) {
				return false
			}
		}
	} else {
		if (!util.CheckIsDirByPath(taskInfo.SrcPath)) {
			Tlog.Printf("源文件路径[%s]不存在!", taskInfo.SrcPath)
			return false
		}
	}

	// 如果任务链表中存在该任务,则不创建
	fTask := task.FindFileTaskByTaskIDFromList(taskInfo.TaskID)
	if fTask == nil {
		fileTask = &task.FileTask {}
	} else {
		fileTask = fTask
	}

	taskType := taskInfo.TaskType
	fileTask.TaskID = taskInfo.TaskID
	fileTask.TaskInfo = taskInfo
	tFileInfo.TaskID = taskInfo.TaskID

	switch taskType {
	case util.TASK_CREATE:			//任务创建
		Tlog.Debugf("Task [%s] Created!\n", taskInfo.TaskID)
		returnValue = handleTaskCreate(tFileInfo, taskInfo)

	case util.TASK_START:			//任务启动
		Tlog.Debugf("Task [%s] Start!\n", taskInfo.TaskID)
		returnValue = handleTaskStart(tFileInfo, fileTask, isLocalDestIP, isLocalHost)

	case util.TASK_SROP:			//任务暂停
		Tlog.Debugf("Task [%s] Stop!\n", taskInfo.TaskID)
		returnValue = handleTaskStop(tFileInfo, taskInfo, isLocalDestIP)

	case util.TASK_UPDATE:			//任务修改
		Tlog.Debugf("Task [%s] Update!\n", taskInfo.TaskID)
		returnValue = handleTaskUpdate(tFileInfo, taskInfo)

	case util.TASK_DELETE:			//任务删除
		Tlog.Debugf("Task [%s] Delete!\n", taskInfo.TaskID)
		returnValue = handleTaskDelete(tFileInfo)
	}

	return returnValue
}

// handle task create and insert to db
func handleTaskCreate(tFileInfo *model.TaskFileInfo, taskInfo *task.TaskInfo) bool {
	// 从数据库中查找该任务节点
	fTask, err := tFileInfo.Find()
	if fTask != nil && err == nil {
		Tlog.Printf("该任务[%s]已存在!\n", taskInfo.TaskID)
		return false
	}

	tFileInfo.SrcHost = taskInfo.SrcHost
	tFileInfo.DestHost = taskInfo.DestHost
	tFileInfo.SrcPath = taskInfo.SrcPath
	tFileInfo.DestPath = taskInfo.DestPath
	tFileInfo.FilePort = taskInfo.FilePort
	tFileInfo.Status = util.TASK_IS_STOP
	tFileInfo.TranType = taskInfo.TranType
	tFileInfo.ScheduleTime = taskInfo.ScheduleTime
	tFileInfo.LastFinishTime = time.Now().Unix()

	err = tFileInfo.Insert()
	if err != nil {
		Tlog.Errorf("新建任务[%s]失败!\n", taskInfo.TaskID)
		return false
	}

	return true
}

// handle task start
func handleTaskStart(tFileInfo *model.TaskFileInfo, fileTask *task.FileTask, isLocalDestIP bool, isLocalHost bool) bool {
	// 从数据库获取任务
	taskFileInfo, err := tFileInfo.Find()
	if err != nil {
		Tlog.Errorf("从数据库获取任务[%s]失败!\n", tFileInfo.TaskID)
		return false
	}

	// 从任务链表中查找该任务
	fTask := task.FindFileTaskByTaskIDFromList(tFileInfo.TaskID)
	if fTask != nil {

		if (fTask.Status == util.TASK_IS_RUNNING) ||
			(taskFileInfo.Status == util.TASK_IS_RUNNING) {
			Tlog.Printf("任务[%s]正在运行,不做处理!\n", tFileInfo.TaskID)
			return true
		}
	}

	// 本地或者本地是目标端传输文件
	if isLocalHost || isLocalDestIP {
		// 检查目标路径是否存在,不存在则创建
		if !util.MkdirAllByPath(fileTask.TaskInfo.DestPath) {
			Tlog.Errorf("目标路径[%s]创建失败!\n", fileTask.TaskInfo.DestPath)
			return false
		}
	}

	// 更新数据库中任务的状态
	tFileInfo.Status = util.TASK_IS_RUNNING
	err = tFileInfo.UpdateTaskStatus()
	if err != nil {
		Tlog.Errorf("更新任务[%s]状态失败!\n", tFileInfo.TaskID)
		return false
	}

	fileTask.Status = util.TASK_IS_RUNNING

	if fTask == nil {
		task.FileTasks = append(task.FileTasks, fileTask)
	}

	// 如果源主机和目标主机是同一个主机
	if isLocalHost {
		go fileTask.CreateFileTranServer()

		// 在这里睡眠是为了保证服务器先准备好接收文件操作
		time.Sleep(5*time.Millisecond)

		go fileTask.HandleTaskStartRequest()

		return true
	}

	// 如果本机是目标主机
	if isLocalDestIP {
		go fileTask.CreateFileTranServer()
	} else {
		go fileTask.HandleTaskStartRequest()
	}

	return true
}

// handle task stop
func handleTaskStop(tFileInfo *model.TaskFileInfo, taskInfo *task.TaskInfo, isLocalDestIP bool) bool {
	// 从任务链表中查找该任务
	fTask := task.FindFileTaskByTaskIDFromList(taskInfo.TaskID)
	if fTask == nil {
		Tlog.Printf("任务[%s]不存在!\n", taskInfo.TaskID)
		return false
	}

	if fTask.Status == util.TASK_IS_STOP {
		Tlog.Printf("任务[%s]已经停止,不做处理!\n", taskInfo.TaskID)
		return true
	}

	fTask.SetFileTaskStatus(util.TASK_IS_STOP)

	tFileInfo.Status = util.TASK_IS_STOP

	// 更新数据库中任务状态
	err := tFileInfo.UpdateTaskStatus()
	if err != nil {
		Tlog.Errorf("更新任务[%s]状态失败!\n", taskInfo.TaskID)
		return false
	}

	// 向文件接收服务器发送停止信号
	if isLocalDestIP {
		conn, err := util.CreateSocketConnect(taskInfo.DestHost, taskInfo.FilePort)
		if err != nil {
			/*从任务链表中删除该任务*/
			task.FileTasks = task.RemoveFileTaskFromList(task.FileTasks, fTask)
			return true
		}
		conn.Close()
	}

	// 从任务链表中删除该任务
	task.FileTasks = task.RemoveFileTaskFromList(task.FileTasks, fTask)

	return true
}

// handle task update
func handleTaskUpdate(tFileInfo *model.TaskFileInfo, taskInfo *task.TaskInfo) bool {
	// 从任务链表中查找该任务
	fTask := task.FindFileTaskByTaskIDFromList(taskInfo.TaskID)
	if fTask != nil {
		if fTask.Status != util.TASK_IS_STOP {
			Tlog.Printf("请停止任务[%s]后修改!\n", fTask.TaskID)
			return false
		} else {
			// 从任务链表中删除该任务
			task.FileTasks = task.RemoveFileTaskFromList(task.FileTasks, fTask)
		}
	}

	tFileInfo.SrcHost = taskInfo.SrcHost
	tFileInfo.DestHost = taskInfo.DestHost
	tFileInfo.FilePort = taskInfo.FilePort
	tFileInfo.SrcPath = taskInfo.SrcPath
	tFileInfo.DestPath = taskInfo.DestPath
	tFileInfo.TranType = taskInfo.TranType
	tFileInfo.ScheduleTime = taskInfo.ScheduleTime

	// 从数据库中修改该任务节点
	err := tFileInfo.Update()
	if err != nil {
		Tlog.Errorf("任务[%s]修改失败!\n", taskInfo.TaskID)
		return false
	}

	return true
}

// handle task delete
func handleTaskDelete(tFileInfo *model.TaskFileInfo) bool {
	// 从任务链表中查找该任务
	fTask := task.FindFileTaskByTaskIDFromList(tFileInfo.TaskID)
	if fTask != nil {
		if fTask.Status != util.TASK_IS_STOP {
			Tlog.Printf("请停止任务[%s]后删除!\n", fTask.TaskID)
			return false
		} else {
			// 从任务链表中删除该任务
			task.FileTasks = task.RemoveFileTaskFromList(task.FileTasks, fTask)
		}
	}

	// 从数据库中查找该任务节点
	tFileInfo, err := tFileInfo.Find()
	if err != nil {
		Tlog.Printf("该任务[%s]不存在!\n", tFileInfo.TaskID)
		return false
	}

	// 从数据库中删除该任务节点
	err = tFileInfo.Delete()
	if err != nil {
		Tlog.Errorf("任务[%s]删除失败!\n", tFileInfo.TaskID)
		return false
	}

	return true
}