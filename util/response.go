package util

import (
	"fmt"
	"encoding/json"
)

type RespMessage struct {
	TaskType 	string
	Status 		bool
}

// Response json pack
func (this *RespMessage) GetRespMessageJson() string {
	respJson := fmt.Sprintf("{\"taskType\":\"%s\",\"status\":%v}", this.TaskType, this.Status)

	return respJson
}

func (this *RespMessage) CheckRespIsTureOrFalse(respJson []byte) bool {
	var respMsg RespMessage

	err := json.Unmarshal(respJson, &respMsg)
	if err != nil {
		return false
	}

	fmt.Println("status: ", respMsg.Status)

	if !respMsg.Status {
		return false
	}

	return true
}