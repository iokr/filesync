package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/dzhenquan/filesync/model"
	"fmt"
	"github.com/dzhenquan/filesync/util"
	"strconv"
	"errors"
	"strings"
)

func AdminTaskIndexGet(c *gin.Context) {

	taskFiles, err := model.FindAllTaskQuery()
	if err != nil {
		taskFiles = nil
	}

	c.HTML(http.StatusOK, "admin/task.html", gin.H{
		"TaskFiles": taskFiles,
	})
}

func AdminTaskPublish(c *gin.Context) {
	var returnValue bool
	id := c.Param("id")

	tFileInfo, err := model.FindTaskByID(id)
	if err != nil {
		fmt.Println("Find task file info is failure!")
		return
	}

	if tFileInfo.Status == util.TASK_IS_STOP {
		returnValue = tFileInfo.SendTaskInfoToLocal(util.TASK_START)
	} else {
		returnValue = tFileInfo.SendTaskInfoToLocal(util.TASK_SROP)
	}

	c.JSON(http.StatusOK, gin.H{
		"succeed": returnValue,
	})
}

func AdminTaskCreateGet(c *gin.Context) {
	c.HTML(http.StatusOK, "task/new.html", nil)
}

func AdminTaskCreatePost(c *gin.Context) {
	taskID 		:= c.PostForm("taskID")
	srcHost		:= c.PostForm("srcHost")
	srcPath 	:= c.PostForm("srcPath")
	port 		:= c.PostForm("filePort")
	destHost 	:= c.PostForm("destHost")
	destPath 	:= c.PostForm("destPath")
	tType 		:= c.PostForm("tranType")
	scheTime 	:= c.PostForm("scheduleTime")

	tranType := 0
	if strings.Compare(tType, "cut") == 0 {
		tranType = 1
	}

	filePort, _ 	:= strconv.Atoi(port)
	scheduleTime, _ := strconv.Atoi(scheTime)

	if scheduleTime < 1 {
		scheduleTime = 1
	}

	tFileInfo := &model.TaskFileInfo{
		TaskID:taskID,
		SrcHost:srcHost,
		SrcPath:srcPath,
		FilePort:filePort,
		DestHost:destHost,
		DestPath:destPath,
		TranType:tranType,
		ScheduleTime:int64(scheduleTime),
	}

	returnValue := tFileInfo.SendTaskInfoToLocal(util.TASK_CREATE)

	c.JSON(http.StatusOK, gin.H{
		"succeed": returnValue,
	})

	return
}

func AdminTaskDisplay(c *gin.Context) {
	id := c.Param("id")

	tFileInfo, err := model.FindTaskByID(id)
	if err != nil {
		fmt.Println("TaskID is not exists!")
		return
	}

	c.HTML(http.StatusOK, "task/display.html", gin.H{
		"TaskInfo": tFileInfo,
	})
	return
}

func AdminTaskEditGet(c *gin.Context) {
	id := c.Param("id")

	tFileInfo, err := model.FindTaskByID(id)
	if err != nil {
		fmt.Println("TaskID is not exists!")
		return
	}

	c.HTML(http.StatusOK, "task/modify.html", gin.H{
		"TaskInfo": tFileInfo,
	})
	return
}

func AdminTaskEditPost(c *gin.Context) {
	taskID 			:= c.Param("id")

	//taskID 		:= c.PostForm("taskID")
	srcHost		:= c.PostForm("srcHost")
	srcPath 	:= c.PostForm("srcPath")
	port 		:= c.PostForm("filePort")
	destHost 	:= c.PostForm("destHost")
	destPath 	:= c.PostForm("destPath")
	tType 		:= c.PostForm("tranType")
	scheTime	:= c.PostForm("scheduleTime")

	tranType := 0
	if strings.Compare(tType, "cut") == 0 {
		tranType = 1
	}

	filePort, _ := strconv.Atoi(port)
	scheduleTime,_ := strconv.Atoi(scheTime)

	if scheduleTime < 1 {
		scheduleTime = 1
	}

	tFileInfo := &model.TaskFileInfo{
		TaskID:taskID,
		SrcHost:srcHost,
		SrcPath:srcPath,
		FilePort:filePort,
		DestHost:destHost,
		DestPath:destPath,
		TranType:tranType,
		ScheduleTime:int64(scheduleTime),
	}

	returnValue := tFileInfo.SendTaskInfoToLocal(util.TASK_UPDATE)

	c.JSON(http.StatusOK, gin.H{
		"succeed": returnValue,
	})

	return
}

func AdminTaskDelete(c *gin.Context) {
	id := c.Param("id")

	tFileInfo, err := model.FindTaskByID(id)
	if err == nil {
		returnValue := tFileInfo.SendTaskInfoToLocal(util.TASK_DELETE)
		if returnValue {
			err = tFileInfo.Delete()
		} else {
			err = errors.New("Delete is failure!")
		}
	} else {
		err = errors.New("TaskID is not exists!")
	}

	c.JSON(http.StatusOK, gin.H{
		"succeed": err == nil,
	})
	return
}