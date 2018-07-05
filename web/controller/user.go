package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"github.com/dzhenquan/filesync/config"
)

// SignIn Get  用户登录
func SigninGet(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/signin.html", nil)
}

func SigninPost(c *gin.Context) {
	useremail	:= c.PostForm("useremail")
	password	:= c.PostForm("password")

	useremail = strings.TrimSpace(useremail)

	if strings.Compare(useremail, config.ServerConfig.WebUser) == 0 &&
		strings.Compare(password, config.ServerConfig.WebPwd) == 0 {

		c.Redirect(http.StatusMovedPermanently, "/admin/index")
		return
	}

	c.HTML(http.StatusOK, "auth/signin.html", gin.H{
		"message": "登录失败!",
	})
}

func AdminProfileGet(c *gin.Context) {


	c.HTML(http.StatusOK, "admin/profile.html", nil)
}

func AdminIndexGet(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/index.html", nil)
}