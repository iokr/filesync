package controller

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/dzhenquan/filesync/model"
)

func AdminTranLogGet(c *gin.Context) {

	tranLogs, err := model.FindAllTranLogQuery()
	if err != nil {
		tranLogs = nil
	}

	c.HTML(http.StatusOK, "admin/tran_log.html", gin.H{
		"TranLogs": tranLogs,
	})
}