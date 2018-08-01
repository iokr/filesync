package task

import (
	"os"
	"io"
	"net"
	"sync"
	"time"
	"strings"
	"github.com/dzhenquan/filesync/config"
	"github.com/dzhenquan/filesync/model"
	"github.com/dzhenquan/filesync/util"
	. "github.com/dzhenquan/filesync/tlog"
	"context"
)

type FileTask struct {
	Status		int
	TaskID		string
	TaskInfo 	*TaskInfo
	mutex 		sync.Mutex
}


var (
	FileTasks []*FileTask
	taskMutex	sync.Mutex
)


func (fileTask *FileTask) NewFileTask() {

}

// 创建文件传输服务器
func (fileTask *FileTask) CreateFileTranServer(ctx context.Context) {
	listen, err := util.CreateSocketListen("", fileTask.TaskInfo.FilePort)
	if err != nil {
		//log.Println("创建本地监听失败!")
		return
	}
	defer listen.Close()

	Tlog.Printf("创建[%s]文件传输服务器成功!\n", fileTask.TaskID)

	for {
		select {
		case <-ctx.Done():
			Tlog.Printf("文件传输服务器[%s]退出!\n", fileTask.TaskID)
			return
		default:
			conn, err := listen.Accept()
			if err != nil {
				Tlog.Errorln("接受新的连接请求失败,file err: ", err)
				continue
			}
			//log.Println("conn: ", conn.RemoteAddr())
			// 为每一个数据连接开一个协程去处理数据
			go fileTask.handleDataConn(ctx, conn)
		}
	}
}

// 开启任务传输请求
func (fileTask *FileTask) HandleTaskStartRequest(ctx context.Context)  {
	select {
	case <-ctx.Done():
		Tlog.Printf("开始退出文件传输[%d]....... \n", fileTask.Status)
		return
	default:
		// 获取源目录文件列表
		filelist, filedir, err := util.GetCurrentFileList(fileTask.TaskInfo.SrcPath)
		if err != nil {
			return
		}
		// 发送目录
		fileTask.handleFileGoSchedule(ctx, filedir, util.TRAN_DIR)

		// 发送文件
		fileTask.handleFileGoSchedule(ctx, filelist, util.TRAN_FILE)

		// 任务完成,更改数据库任务完成时间
		fileTask.handleTaskFinishUpdateStatusTime(ctx)
	}
	Tlog.Println("Exit HandleTaskStartRequest.")
}

// 处理数据连接
func (fileTask *FileTask) handleDataConn(ctx context.Context, conn net.Conn) {
	for {
		dataBuf := make([]byte, util.MAX_MESSAGE_LEN)

		recvSize, err := conn.Read(dataBuf)
		if err != nil {
			if err == io.EOF {
				return
			}
			Tlog.Errorln("读取文件属性信息失败!")
			continue
		}
		defer conn.Close()

		tFileInfo, err := MarshalJsonToStruct(dataBuf[:recvSize])
		if err != nil {
			Tlog.Errorln("err:", err)
			continue
		}

		// 接收json报文成功,返回成功消息
		conn.Write([]byte("ok"))

		select {
		case <-ctx.Done():
			Tlog.Printf("文件传输服务器[%s]退出!\n", fileTask.TaskID)
			return
		default:
			filename := fileTask.TaskInfo.DestPath + tFileInfo.FilePath
			if strings.Compare(tFileInfo.FileType, "dir") == 0 {
				// 处理文件夹数据连接
				fileTask.handleDirDataConn(ctx, conn, filename)
			} else if strings.Compare(tFileInfo.FileType, "file") == 0 {
				// 处理文件数据连接
				fileTask.handleFileDataConn(ctx, conn, filename, tFileInfo.FileSize)
			}
		}
	}
}

// 处理文件夹传输数据连接
func (fileTask *FileTask) handleDirDataConn(ctx context.Context, conn net.Conn, dirPath string) {
	select {
	case <-ctx.Done():
		Tlog.Printf("文件传输服务器[%s]退出!\n", fileTask.TaskID)
		return
	default:
		// 检查本地是否存在该文件夹,不存在则创建
		if !util.CheckIsDirByPath(dirPath) {
			util.MkdirAllByPath(dirPath)
		}
	}
}

// 处理文件传输数据连接
func (fileTask *FileTask) handleFileDataConn(ctx context.Context, conn net.Conn, filename string, filesize int64) {
	var transResult string

	startTime := time.Now().Format("2006-01-02 15:04:05")
	tFileLog := &model.TaskFileLog{
		TaskID:fileTask.TaskID,
		SrcHost:fileTask.TaskInfo.SrcHost,
		DestHost:fileTask.TaskInfo.DestHost,
		FileName: filename,
		FileSize: filesize,
		FileStartTime: startTime,
	}

	select {
	case <-ctx.Done():
		Tlog.Printf("文件传输服务器[%s]退出!\n", fileTask.TaskID)
		return
	default:
		if _, err := util.RecvFile(ctx, conn, filename, uint64(filesize)); err == nil {
			Tlog.Printf("文件[%s]接收完毕,TaskID:[%s]!\n", filename,fileTask.TaskID)

			transResult = "文件接收成功"
			if fileTask.TaskInfo.TranType == util.FILE_CUT{
				conn.Write([]byte("ok"))
			}
			//hash, _ := HashFile(filename)
			//fmt.Println("md5:", hash)
		} else {
			transResult = "文件接收失败"
			Tlog.Printf("文件[%s]接收失败,TaskID:[%s]!\n", filename,fileTask.TaskID)
		}
	}

	finishTime := time.Now().Format("2006-01-02 15:04:05")
	tFileLog.FileEndTime = finishTime
	tFileLog.TransResult = transResult

	tFileLog.Insert()
}

// 根据每次传输文件的个数,对所有文件列表进行分页调度
func (fileTask *FileTask) handleFileGoSchedule(ctx context.Context, filelist []string, tranType int) {
	fileTotalCount := len(filelist)
	if fileTotalCount <= 0 {
		return
	}

	// 根据每次传输文件的个数,对所有文件列表进行分页调度
	if fileTotalCount > config.ServerConfig.MaxFtsNum {
		fileTotalPage := fileTotalCount / config.ServerConfig.MaxFtsNum
		if fileTotalCount % config.ServerConfig.MaxFtsNum > 0 {
			fileTotalPage++
		}

		var endFile int
		for i := 0; i < fileTotalPage; i++ {
			startFile := i*config.ServerConfig.MaxFtsNum
			curfileCount := fileTotalCount - startFile
			if curfileCount >= config.ServerConfig.MaxFtsNum {
				endFile = startFile + config.ServerConfig.MaxFtsNum
			} else {
				endFile = startFile + curfileCount
			}

			select {
			case <-ctx.Done():
				Tlog.Printf("开始退出文件传输[%d]....... \n", fileTask.Status)
				return
			default:
				fileTask.handleMaxFileTransNums(ctx, tranType, filelist[startFile:endFile])
			}
		}
	} else {
		fileTask.handleMaxFileTransNums(ctx, tranType, filelist[:fileTotalCount])
	}
}

// 根据分页调度创建传输文件的最大协程数
func (fileTask *FileTask) handleMaxFileTransNums(ctx context.Context, tranType int, filelist []string) {
	fileThreadCount := len(filelist)

	var waitGroup sync.WaitGroup
	for i := 0; i < fileThreadCount; i++ {
		select {
		case <-ctx.Done():
			Tlog.Printf("开始退出文件传输[%d]....... \n", fileTask.Status)
			return
		default:
			// 将windows下文件路径中反斜杠进行转化
			filePath := strings.Replace(filelist[i], `\`,`/`,-1)

			// 如果是文件复制且传输的是文件
			if (fileTask.TaskInfo.TranType == util.FILE_COPY) &&
				(tranType == util.TRAN_FILE) {
				// 从数据库中检查该文件是否更改
				if fileTask.checkIsExistsFileCopyFromDB(filePath) {
					continue
				}
			}

			waitGroup.Add(1)
			go fileTask.handleDataTran(ctx, filePath, &waitGroup)
		}
	}
	waitGroup.Wait()
}

// 处理数据传输,目录传输和文件传输
func (fileTask *FileTask) handleDataTran(ctx context.Context, filePath string, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	// 获取文件属性
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		Tlog.Errorf("File [%s] is not exists!\n", filePath)
		return
	}

	// 去掉源文件路径前缀
	newPath := strings.TrimPrefix(filePath, fileTask.TaskInfo.SrcPath)

	// 连接到文件传输服务器
	conn, err := util.CreateSocketConnect(fileTask.TaskInfo.DestHost, fileTask.TaskInfo.FilePort)
	if err != nil {
		return
	}
	defer conn.Close()

	// 获取发送报文
	dataByte, _ := GetMarshalToJson(newPath, fileInfo)

	// 发送报文
	_, err = conn.Write([]byte(dataByte))
	if err != nil {
		Tlog.Errorln("Send data string is failure!")
		return
	}

	// 接收返回消息
	dataBuf := make([]byte, util.MAX_MESSAGE_LEN)
	recvSize, err := conn.Read(dataBuf)
	if err != nil {
		return
	}

	recvOk := dataBuf[:recvSize]
	if strings.Compare(string(recvOk), "ok") == 0 {
		// 如果不是目录,处理文件传输
		if !fileInfo.IsDir() {
			fileTask.handleFileDataTran(ctx, conn, filePath, fileInfo)

			// 如果是文件移动策略,则删除
			fileTask.handleDataCutTran(ctx, conn, filePath)
		}
	}
	return
}

// 处理数据传输方式为移动
func (fileTask *FileTask) handleDataCutTran(ctx context.Context, conn net.Conn, filePath string) bool {

	select {
	case <-ctx.Done():
		return false
	default:
		// 如果是文件移动策略,则删除
		if fileTask.TaskInfo.TranType == util.FILE_CUT {
			bufDelete := make([]byte, util.MAX_MESSAGE_LEN)

			recvSize, err := conn.Read(bufDelete)
			if err != nil {
				return false
			}
			recvDelete := bufDelete[:recvSize]
			if strings.Compare(string(recvDelete), "ok") == 0 {
				os.Remove(filePath)
			}
		}
	}
	return true
}

// 处理文件数据传输
func (fileTask *FileTask) handleFileDataTran(ctx context.Context, conn net.Conn, filePath string, fileInfo os.FileInfo) {
	var transResult string

	// 文件传输日志记录到数据库
	startTime := time.Now().Format("2006-01-02 15:04:05")
	tFileLog := &model.TaskFileLog{
		TaskID: fileTask.TaskID,
		SrcHost:fileTask.TaskInfo.SrcHost,
		DestHost:fileTask.TaskInfo.DestHost,
		FileName: fileInfo.Name(),
		FileSize: fileInfo.Size(),
		FileStartTime: startTime,
	}

	if sendSize, err := util.SendFile(ctx, conn, filePath); err != nil {
		Tlog.Printf("文件[%s]发送失败,TaskID:[%s],Len:[%d]!\n",filePath, fileTask.TaskID ,sendSize)
		transResult = "文件发送失败"
	} else {
		Tlog.Printf("文件[%s]发送完毕,TaskID:[%s],Len:[%d]!\n",filePath, fileTask.TaskID, sendSize)
		transResult = "文件发送成功"

		if fileTask.TaskInfo.TranType == util.FILE_COPY {
			// 文件Copy日志写入数据库
			fileTask.copyFileLogToDB(filePath, sendSize)
		}
	}

	finishTime := time.Now().Format("2006-01-02 15:04:05")
	tFileLog.FileEndTime = finishTime
	tFileLog.TransResult = transResult

	// 文件传输日志插入到数据库
	tFileLog.Insert()
}

// 处理任务完成后续工作
func (fileTask *FileTask) handleTaskFinishUpdateStatusTime(ctx context.Context) error {

	time.Sleep(5*time.Millisecond)

	nowTime := time.Now().Unix()

	select {
	case <-ctx.Done():
		Tlog.Println("开始退出文件传输,....... ", fileTask.Status)
		return nil
	default:
		fileTask.SetFileTaskStatus(util.TASK_IS_RUNED)

		// 更新数据库中任务状态
		tFileInfo := model.TaskFileInfo{
			TaskID:fileTask.TaskID,
			Status:fileTask.Status,
			LastFinishTime:nowTime,
		}
		return tFileInfo.UpdateTaskStatusTime()
	}
}

// 从复制数据库中查找文件
func (fileTask *FileTask) findCopyFileFromDB(filename, filemd5 string) int {
	taskFileCopy := &model.TaskFileCopy{
		FileName:filename,
	}

	// 不存在该文件
	fileCopy,err := taskFileCopy.Find()
	if err != nil {
		return 0
	}

	// 存在该文件且MD5值相等
	if strings.Compare(fileCopy.FileMd5, filemd5) == 0 {
		return -1
	}

	//存在该文件且MD5值不相等
	return 1
}

// 将复制方式传输文件日志存入数据库
func (fileTask *FileTask) copyFileLogToDB(filename string, filesize int64) bool {
	// 根据文件全路径获取文件md5
	fileMd5, err := util.HashFile(filename)
	if err != nil {
		Tlog.Errorf("获取文件[%s]MD5失败!\n", filename)
		return false
	}

	tFileCopy := &model.TaskFileCopy{
		TaskID:fileTask.TaskID,
		FileName:filename,
		FileSize:filesize,
		FileMd5:fileMd5,
	}

	nReturnValue := fileTask.findCopyFileFromDB(filename, fileMd5)
	if nReturnValue == -1 {
		return true
	} else if nReturnValue == 1 {
		err := tFileCopy.Update()
		if err != nil {
			return false
		}
	} else {
		err := tFileCopy.Insert()
		if err != nil {
			return false
		}
	}
	return true
}

// 检查数据库中是否存在该文件
func (fileTask *FileTask) checkIsExistsFileCopyFromDB(filename string) bool {
	// 根据文件全路径获取文件md5
	fileMd5, err := util.HashFile(filename)
	if err != nil {
		Tlog.Errorf("获取文件[%s]MD5失败!\n", filename)
		return true
	}

	nReturnValue := fileTask.findCopyFileFromDB(filename, fileMd5)
	if nReturnValue == -1 {
		return true
	}

	return false
}

// 设置文件传输任务状态
func (fileTask *FileTask) SetFileTaskStatus(status int) {
	fileTask.mutex.Lock()
	defer fileTask.mutex.Unlock()

	fileTask.Status = status
}

// 从任务切片链表中查找任务
func FindFileTaskByTaskIDFromList(taskID string) (*FileTask) {
	taskMutex.Lock()
	defer taskMutex.Unlock()

	if len(taskID) == 0 {
		return nil
	}

	// 从任务链表中查找该任务
	//for i, _ := range FileTasks {
	//	if strings.Compare(FileTasks[i].TaskID, taskID) == 0 {
	//		return FileTasks[i]
	//	}
	//}
	for i := 0; i < len(FileTasks); i++ {
		if strings.Compare(FileTasks[i].TaskID, taskID) == 0 {
			return FileTasks[i]
		}
	}

	return nil
}

// 从任务切片链表中删除任务
func RemoveFileTaskFromList(slice []*FileTask, elems ...*FileTask) []*FileTask {
	taskMutex.Lock()
	defer taskMutex.Unlock()

	isInElems := make(map[*FileTask]bool)
	for _, elem := range elems {
		isInElems[elem] = true
	}
	w := 0
	for _, elem := range slice {
		if !isInElems[elem] {
			slice[w] = elem
			w += 1
		}
	}
	return slice[:w]
}