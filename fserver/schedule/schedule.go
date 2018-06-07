package schedule

import (
	"github.com/dzhenquan/filesync/util"
	"time"
	"github.com/dzhenquan/filesync/model"
	"log"
)



func HandleTaskSchedule() {
	for {
		// 从数据库中加载所有任务
		tFileInfos, err := model.FindAllTaskQuery()
		if err != nil {
			log.Println("从数据库获取所有任务失败!")
			continue
		}

		for _, tFileInfo := range tFileInfos {
			if tFileInfo.Status != util.TASK_IS_STOP {
				nowTime := time.Now().Unix()
				deltaTime := nowTime-tFileInfo.LastFinishTime
				if deltaTime >= tFileInfo.ScheduleTime {
					tFileInfo.SendTaskInfoToLocal(util.TASK_START)
				}
			}
		}
		time.Sleep(2*time.Second)
	}
}

/*
func FindFileTaskByTaskID(taskID string) bool {

	if len(taskID) == 0 {
		return true
	}

	for _, fileTask := range task.ScheduleTasks {
		if strings.Compare(fileTask.TaskID, taskID) == 0 {
			return true
		}
	}
	return false
}

func RemoveFileTaskByTaskID(slice []*task.FileTask, elems ...*task.FileTask) []*task.FileTask {
	isInElems := make(map[*task.FileTask]bool)
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
*/