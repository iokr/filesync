package util

const (
	TASK_CREATE  	= "1010"	//任务创建
	TASK_START		= "1020"	//任务开始
	TASK_SROP		= "1030"	//任务暂停
	TASK_UPDATE		= "1040"	//任务修改
	TASK_DELETE		= "1050"	//任务删除
	TASK_UNABLE		= "1060"	//无效任务
)

const (
	FILE_COPY		= iota
	FILE_CUT
)

const (
	TASK_IS_STOP	= iota
	TASK_IS_RUNED
	TASK_IS_RUNNING
)

var (
	MSG_TRAN_PORT		= 8686			//消息指令传输端口
	MAX_TRAN_FILE_NUM 	= 5   			//每次最大同时传输文件个数
	MAX_MESSAGE_LEN		= 1024			//每次接收消息的最大长度(字节)
	MAX_FILE_DATA_LEN 	= 1024*1024*16	//每次接收文件数据的最大长度(字节)
)

/*
{"taskID":"xiaodai","taskType":"1030","filePort":8787,
"srcHost":"127.0.0.1","destHost":"10.0.0.190",
"srcPath":"/root/Golang/src/filetrans/SendFile","destPath":"/dev" }
*/