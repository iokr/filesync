package task

import "encoding/json"

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