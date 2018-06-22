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
		time.Sleep(500*time.Millisecond)
	}
}
