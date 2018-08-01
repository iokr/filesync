filesync 是使用gin框架提供web界面的文件传输工具  
用户通过web界面来创建，删除，修改，开启，暂停任务  

首先使用go get下载源代码  
```
[root@localhost ~]# go get github.com/dzhenquan/filesync
```
或者直接将源代码下载放到自己的go环境中使用  
```
git clone https://github.com/dzhenquan/filesync.git
```
*注意：* 
用户使用首先根据自己需要修改config.json配置文件

后台服务器界面如下所示:  
任务管理  
![任务管理](https://raw.githubusercontent.com/dzhenquan/filesync/master/images/task_manager.png)

任务添加  
![任务添加](https://raw.githubusercontent.com/dzhenquan/filesync/master/images/task_add.png)

任务修改  
![任务修改](https://raw.githubusercontent.com/dzhenquan/filesync/master/images/task_update.png)

任务查看  
![任务查看](https://raw.githubusercontent.com/dzhenquan/filesync/master/images/task_show.png)

日志查看  
![日志查看](https://raw.githubusercontent.com/dzhenquan/filesync/master/images/log_manager.png)
