package controller

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/dzhenquan/filesync/model"
	"strconv"
)

func AdminTranLogGet(c *gin.Context) {

	page := c.Param("page")
	pageInt, _ := strconv.Atoi(page)

	curPage := pageInt
	curOffset := 15 * (curPage-1)

	totalCount := model.GetAllTranLogCount()
	totalPage := totalCount / 15
	if totalCount % 15 != 0 {
		totalPage = totalPage + 1
	}

	tranLogs, err := model.FindAllTranLogQuery(curOffset)
	if err != nil {
		tranLogs = nil
	}

	c.HTML(http.StatusOK, "admin/tran_log.html", gin.H{
		"TranLogs": tranLogs,
		"prevPage": curPage-1,
		"curPage": curPage,
		"nextPage": curPage+1,
		"totalCount":totalCount,
		"totalPage": totalPage,
	})
}