package config

import (
	"os"
	"fmt"
	"regexp"
	"io/ioutil"
	"encoding/json"
	"github.com/dzhenquan/filesync/util"
)


var jsonData map[string]interface{}

func initJSON() {
	bytes, err := ioutil.ReadFile("../config.json")
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		os.Exit(-1)
	}

	configStr := string(bytes[:])
	reg := regexp.MustCompile(`/\*.*\*/`)

	configStr = reg.ReplaceAllString(configStr, "")
	bytes = []byte(configStr)

	if err := json.Unmarshal(bytes, &jsonData); err != nil {
		fmt.Println("invalid config: ", err.Error())
		os.Exit(-1)
	}

	return
}


type dBConfig struct {
	Dialect				string
	Database			string
	User 				string
	Password 			string
	Host 				string		// 数据库ip
	Port 				int			// 数据库端口
	Charset 			string
	URL 				string
	MaxIdleConns 		int			// 空闲时最大的连接数
	MaxOpenConns 		int			// 最大的连接数
}

// DBConfig 数据库相关配置
var DBConfig dBConfig

func initDB() {
	util.SetStructByJSON(&DBConfig, jsonData["database"].(map[string]interface{}))
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		DBConfig.User, DBConfig.Password, DBConfig.Host, DBConfig.Port, DBConfig.Database, DBConfig.Charset)
	DBConfig.URL = url
}

type serverConfig struct {
	WebPort				int			// Web服务器监听端口
	FServerPort 		int			// 文件传输服务器监听端口
	MaxMultipartMemory 	int			// 上传的图片最大允许的大小，单位MB
	MaxFtsNum			int 		// 每次最大同时传输文件个数
	WebHost 			string 		// web服务器监听IP(默认监听所有)
	FServerHost 		string 		// FServer服务器监听IP(默认监听所有)
	WebUser				string 		// web界面登陆用户名
	WebPwd 				string 		// web界面登陆密码
	Env                	string		// 模式(开发，测试，产品)
	LogDir            	string		// 日志文件所在的目录，如果不设的话，默认在项目目录下
	LogOldDir 			string		// 旧日志文件所在的目录，如果不设的话，默认在项目目录下
	LogPrefix 			string 		// 日志文件前缀名(默认program)
	LogSuffix 			string 		// 日志文件后缀名(默认log)
	LogMaxLine 			int64 		// 日志文件的最大行数(默认10000行)
	LogMaxByte 			int64 		// 日志文件的最大字节数(默认52428800字节(50M))
	LogFile            	string		// 日志文件

/*
	APIPoweredBy       string
	SiteName           string
	Host               string
	ImgHost            string

	APIPrefix          string
	UploadImgDir       string
	ImgPath            string

	Port               int
	TokenSecret        string
	TokenMaxAge        int
	PassSalt           string
	LuosimaoVerifyURL  string
	LuosimaoAPIKey     string
	CrawlerName        string
	MailUser           string //域名邮箱账号
	MailPass           string //域名邮箱密码
	MailHost           string //smtp邮箱域名
	MailPort           int    //smtp邮箱端口
	MailFrom           string //邮件来源
	Github             string
	BaiduPushLink      string
*/
}

// ServerConfig 服务器相关配置
var ServerConfig serverConfig

func initServer() {
	util.SetStructByJSON(&ServerConfig, jsonData["go"].(map[string]interface{}))
}

func init() {
	initJSON()
	initDB()
	initServer()
}