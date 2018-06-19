package model

import "fmt"

type TaskFileLog struct {
	ID        		uint64	`gorm:"primary_key"`
	TaskID 			string
	SrcHost 		string
	DestHost		string
	FileName		string	`json:"fileName"`
	FileSize		string 	`json:"fileSize"`
	FileStartTime	string
	FileEndTime		string
	TransResult 	string
}

func (tFilelog *TaskFileLog) Insert() error {
	return DB.Create(tFilelog).Error
}

func GetAllTranLogCount() uint64 {
	var totalCount uint64

	row := DB.Raw("select count(*) from task_file_log;").Row()

	err := row.Scan(&totalCount)
	if err != nil {
		totalCount = 0
	}

	return totalCount
}

func FindAllTranLogQuery(curOffset int) ([]*TaskFileLog, error) {
	//select * from article ORDER BY updated_at desc LIMIT ? offset ?
	sql := fmt.Sprintf("select * from task_file_log order by id desc LIMIT 15 offset %d;", curOffset)

	rows, err := DB.Raw(sql).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tranLogs []*TaskFileLog
	for rows.Next() {
		var tranLog TaskFileLog
		err := DB.ScanRows(rows, &tranLog)
		if err != nil {
			return nil, err
		}
		tranLogs = append(tranLogs, &tranLog)
	}
	return tranLogs, nil
}