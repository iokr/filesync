package util

import (
	"net"
	"strings"
	"path/filepath"
	"os"
	"errors"
)

// GetLocalIP returns the local ip address
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	bestIP := ""
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			//fmt.Println(ipnet.IP.String())
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return bestIP
}

// Check IP is local ip address
// returns true or false
func CheckIPIsLocalIP(ipaddr string) bool {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return false
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			//fmt.Println(ipnet.IP.To4().String())
			if strings.Compare(ipnet.IP.String(), ipaddr) == 0 {
				return true
			}
		}
	}
	return false
}


// Get current file list by file full path
// returns the file list or err
func GetCurrentFileList(filePath string) (filelist , filedir []string, err error) {

	err = filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			filelist = append(filelist, path)
			if len(filelist) > 100000 {
				err = errors.New("文件列表达到最大限制10000！")
				return err
			}
		} else {
			if strings.Compare(filePath, path) != 0 {
				filedir = append(filedir, path)
			}
		}
		return nil
	})

	return filelist, filedir, err
}

// check path is dir
// returns true/false
func checkPathIsDir(filePath string) bool {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return os.IsExist(err)
	}
	return fileInfo.IsDir()
}

// check file path is exists
// returns true/false
func checkPathIsExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}