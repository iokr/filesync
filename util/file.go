package util

import (
	"os"
	"net"
	"io/ioutil"
	"io"
)

type FileInfo struct {
	Name 	string 		`json:"name"`
	Size 	uint64 		`json:"size"`
}

// Send file by conn and file name
// returns send file size or err
func SendFile(conn net.Conn, filename string) (sendLen int, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return -1, err
	}
	defer f.Close()

	fileContent, err := ioutil.ReadAll(f)
	if err != nil {
		return -1, err
	}

	sendSize, err := conn.Write(fileContent)
	if err != nil {
		return -1, err
	}

	return sendSize, nil
}


// Recv file by conn, file name and file size
// returns true/false or err
func RecvFile(conn net.Conn, filename string, filesize uint64) (flag bool, err error) {
	f, err := os.Create(filename)
	if err != nil {
		f.Close()
		return false, err
	}
	defer f.Close()

	fileBuf := make([]byte, MAX_FILE_DATA_LEN)

	for filesize > 0{
		readLen, err := conn.Read(fileBuf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return false, err
			}
		}
		f.Write(fileBuf[:readLen])

		filesize = filesize-uint64(readLen)
	}
	return true, nil
}


// check path is dir
// returns true/false
func CheckIsDirByPath(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}