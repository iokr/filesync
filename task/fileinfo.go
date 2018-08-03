package task

import (
	"encoding/json"
	"os"
)

type TFileInfo struct {
	FilePath	string		`json:"filePath"`
	FileType 	string 		`json:"fileType"`
	FileSize 	int64		`json:"fileSize"`
}

func (tFileInfo *TFileInfo) MarshalToJson() ([]byte, error) {
	tJson, err := json.Marshal(&tFileInfo)
	if err != nil {
		return nil, err
	}
	return tJson, nil
}

func MarshalJsonToStruct(jsonByte []byte) (*TFileInfo, error) {
	taskFileInfo := TFileInfo{}

	err := json.Unmarshal(jsonByte, &taskFileInfo)

	return &taskFileInfo, err
}

func GetMarshalToJson(filePath string, fileInfo os.FileInfo) ([]byte, error) {
	tFileInfo := &TFileInfo{
		FilePath: filePath,
		FileSize:fileInfo.Size(),
	}

	if fileInfo.IsDir() {
		tFileInfo.FileType = "dir"
	} else {
		tFileInfo.FileType = "file"
	}

	// 将结构体转化为json报文
	dataByte, err := tFileInfo.MarshalToJson()
	if err != nil {
		return nil, err
	}
	return dataByte, nil
}