package main

import (
	"os"
	"log"
	"io"
	"net"
	"fmt"
	"encoding/json"
	"github.com/dzhenquan/filesync/util"
	"github.com/dzhenquan/filesync/task"
	"github.com/dzhenquan/filesync/model"
	"strings"
	"time"
	"github.com/dzhenquan/filesync/fserver/schedule"
)

func main() {

	//Init Model
	db, err := model.InitDB()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	defer db.Close()

	listen, err := util.CreateSocketListen("", util.MSG_TRAN_PORT)
	if err != nil {
		log.Println("创建本地监听失败!")
		os.Exit(-1)
	}
	defer listen.Close()

	go schedule.HandleTaskSchedule()

	log.Println("Message Listen ...")

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal("接受新的连接请求失败, err: ", err)
			continue
		}
		log.Println("conn: ", conn.RemoteAddr())
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
		log.Println("Read client message is failure!")
		return
	}
	defer conn.Close()

	log.Println("recvBuf: ", string(recvBuf[:readLen]), "len:", readLen)

	respMsg.Status = handleTaskRequest(recvBuf[:readLen], respMsg)

	respJson := respMsg.GetRespMessageJson()
	_, err = conn.Write([]byte(respJson))
	if err != nil {
		log.Println("向客户端回复返回消息失败!")
	}
}

// handle client task request
//returns true/false
func handleTaskRequest(taskJson []byte, respMsg *util.RespMessage) bool {
	var returnValue bool

	taskInfo := &task.TaskInfo{}
	fileTask := &task.FileTask {
		Quit:make(chan bool, 1),
	}
	tFileInfo := &model.TaskFileInfo{}

	err := json.Unmarshal(taskJson, taskInfo)
	if err != nil {
		log.Println("解析客户端发送的JSON报文失败!")
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
				log.Printf("源文件路径[%s]不存在!", taskInfo.SrcPath)
				return false
			}
			if !taskInfo.SendTaskInfoJson(taskJson) {
				return false
			}
		}
	} else {
		if (!util.CheckIsDirByPath(taskInfo.SrcPath)) {
			log.Printf("源文件路径[%s]不存在!", taskInfo.SrcPath)
			return false
		}
	}

	taskType := taskInfo.TaskType
	fileTask.TaskID = taskInfo.TaskID
	fileTask.TaskInfo = taskInfo
	tFileInfo.TaskID = taskInfo.TaskID

	switch taskType {
	case util.TASK_CREATE:
		log.Println("Task Created.")

		returnValue = handleTaskCreate(tFileInfo, taskInfo)

	case util.TASK_START:
		fmt.Println("Task Start.")

		returnValue = handleTaskStart(tFileInfo, fileTask, isLocalDestIP, isLocalHost)

	case util.TASK_SROP:
		fmt.Println("Task Stop.")

		returnValue = handleTaskStop(tFileInfo, taskInfo, isLocalDestIP)

		fmt.Println("returnValue: ", returnValue)

	case util.TASK_UPDATE:
		fmt.Println("Task Update")

		returnValue = handleTaskUpdate(tFileInfo, taskInfo)

	case util.TASK_DELETE:
		fmt.Println("Task Delete.")

		returnValue = handleTaskDelete(tFileInfo)
	}

	return returnValue
}

// handle task create and insert to db
func handleTaskCreate(tFileInfo *model.TaskFileInfo, taskInfo *task.TaskInfo) bool {
	/*从数据库中查找该任务节点*/
	fTask, err := tFileInfo.Find()
	if fTask != nil && err == nil {
		log.Printf("该任务[%s]已存在!", taskInfo.TaskID)
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
		log.Printf("新建任务[%s]失败!", taskInfo.TaskID)
		return false
	}

	return true
}

// handle task start
func handleTaskStart(tFileInfo *model.TaskFileInfo, fileTask *task.FileTask, isLocalDestIP bool, isLocalHost bool) bool {
	//从数据库获取任务
	taskFileInfo, err := tFileInfo.Find()
	if err != nil {
		log.Println("从数据库获取任务失败!")
		return false
	}

	//从任务链表中查找该任务
	fTask := task.FindFileTaskByTaskIDFromList(tFileInfo.TaskID)
	if fTask != nil {
		if (fTask.Status == util.TASK_IS_RUNNING) ||
			(taskFileInfo.Status == util.TASK_IS_RUNNING) {
			log.Println("任务正在运行,不做处理!")
			return true
		}
	} else {
		task.FileTasks = append(task.FileTasks, fileTask)
	}

	//修改数据库中任务的状态
	tFileInfo.Status = util.TASK_IS_RUNNING
	err = tFileInfo.UpdateTaskStatus()
	if err != nil {
		fmt.Printf("更新任务[%s]状态失败!", tFileInfo.TaskID)
		return false
	}

	fileTask.Status = util.TASK_IS_RUNNING

	if isLocalHost {
		// 检查目标路径是否存在,不存在则创建
		if (!util.CheckIsDirByPath(fileTask.TaskInfo.DestPath)) {
			err := os.MkdirAll(fileTask.TaskInfo.DestPath, os.ModePerm)
			if err != nil {
				log.Printf("目标路径[%s]创建失败!", fileTask.TaskInfo.DestPath)
				return false
			}
		}

		go fileTask.CreateFileTranServer()

		// 在这里睡眠是为了保证服务器先准备好接收文件操作
		time.Sleep(5*time.Millisecond)

		go fileTask.HandleTaskStartRequest()

		return true
	}

	if isLocalDestIP {
		// 检查目标路径是否存在,不存在则创建
		if (!util.CheckIsDirByPath(fileTask.TaskInfo.DestPath)) {
			err := os.MkdirAll(fileTask.TaskInfo.DestPath, os.ModePerm)
			if err != nil {
				log.Printf("目标路径[%s]创建失败!", fileTask.TaskInfo.DestPath)
				return false
			}
		}

		go fileTask.CreateFileTranServer()
	} else {
		go fileTask.HandleTaskStartRequest()
	}

	return true
}

// handle task stop
func handleTaskStop(tFileInfo *model.TaskFileInfo, taskInfo *task.TaskInfo, isLocalDestIP bool) bool {
	//从任务链表中查找该任务
	fTask := task.FindFileTaskByTaskIDFromList(taskInfo.TaskID)
	if fTask == nil {
		log.Printf("任务[%s]不存在!", taskInfo.TaskID)
		return false
	}

	if fTask.Status == util.TASK_IS_STOP {
		log.Printf("任务[%s]已经停止,不做处理!", taskInfo.TaskID)
		return true
	}

	fTask.Quit <- true
	fTask.SetFileTaskStatus(util.TASK_IS_STOP)
	//fTask.Status = util.TASK_IS_STOP

	tFileInfo.Status = util.TASK_IS_STOP

	err := tFileInfo.UpdateTaskStatus()
	if err != nil {
		log.Printf("更新任务[%s]状态失败!", taskInfo.TaskID)
		return false
	}

	if isLocalDestIP {
		conn, err := util.CreateSocketConnect(taskInfo.DestHost, taskInfo.FilePort)
		if err != nil {
			/*从任务链表中删除该任务*/
			task.FileTasks = task.RemoveFileTaskFromList(task.FileTasks, fTask)
			return true
		}
		conn.Close()
	}

	/*从任务链表中删除该任务*/
	task.FileTasks = task.RemoveFileTaskFromList(task.FileTasks, fTask)

	return true
}

// handle task update
func handleTaskUpdate(tFileInfo *model.TaskFileInfo, taskInfo *task.TaskInfo) bool {
	/*从任务链表中查找该任务*/
	fTask := task.FindFileTaskByTaskIDFromList(taskInfo.TaskID)
	if fTask != nil {
		if fTask.Status != util.TASK_IS_STOP {
			log.Printf("请停止任务[%s]后修改!", fTask.TaskID)
			return false
		}
	}

	tFileInfo.SrcHost = taskInfo.SrcHost
	tFileInfo.DestHost = taskInfo.DestHost
	tFileInfo.FilePort = taskInfo.FilePort
	tFileInfo.SrcPath = taskInfo.SrcPath
	tFileInfo.DestPath = taskInfo.DestPath
	tFileInfo.TranType = taskInfo.TranType
	tFileInfo.ScheduleTime = taskInfo.ScheduleTime

	err := tFileInfo.Update()
	if err != nil {
		log.Printf("任务[%s]修改失败!", taskInfo.TaskID)
		return false
	}

	return true
}

// handle task delete
func handleTaskDelete(tFileInfo *model.TaskFileInfo) bool {
	/*从任务链表中查找该任务*/
	fTask := task.FindFileTaskByTaskIDFromList(tFileInfo.TaskID)
	if fTask != nil {
		if fTask.Status != util.TASK_IS_STOP {
			log.Printf("请停止任务[%s]后删除!", fTask.TaskID)
			return false
		}
	}

	/*从数据库中查找该任务节点*/
	tFileInfo, err := tFileInfo.Find()
	if err != nil {
		log.Printf("该任务[%s]不存在!", tFileInfo.TaskID)
		return false
	}

	err = tFileInfo.Delete()
	if err != nil {
		log.Printf("任务[%s]删除失败!", tFileInfo.TaskID)
		return false
	}

	return true
}