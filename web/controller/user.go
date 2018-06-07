package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"fmt"
)

// SignIn Get  用户登录
func SigninGet(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/signin.html", nil)
}

func SigninPost(c *gin.Context) {
	useremail	:= c.PostForm("useremail")
	password	:= c.PostForm("password")

	useremail = strings.TrimSpace(useremail)

	fmt.Println("email: ", useremail)
	fmt.Println("pword: ", password)

	if strings.Compare(useremail, "admin") == 0 &&
		strings.Compare(password, "123456") == 0 {

		c.Redirect(http.StatusMovedPermanently, "/admin/index")
		return
	}

	c.HTML(http.StatusOK, "auth/signin.html", gin.H{
		"message": "登录失败!",
	})
}


func AdminIndexGet(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/index.html", nil)
}