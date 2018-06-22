package model

import "fmt"

type TaskFileCopy struct {
	ID        		uint64	`gorm:"primary_key"`
	TaskID 			string
	FileName 		string
	FileSize 		int64
	FileMd5			string
}

func (tFileCopy *TaskFileCopy) Insert() error {
	return DB.Create(tFileCopy).Error
}

func (tFileCopy *TaskFileCopy) Find() ( *TaskFileCopy,  error) {
	var taskFileCopy TaskFileCopy

	err := DB.Where("file_name = ?", tFileCopy.FileName).First(&taskFileCopy).Error

	return &taskFileCopy, err
}

func (tFileCopy *TaskFileCopy) Update() error {

	sql := fmt.Sprintf("update task_file_copy set file_md5='%s' where file_name='%s';",
		tFileCopy.FileMd5, tFileCopy.FileName)

	return DB.Exec(sql).Error
}