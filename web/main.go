package main

import (
	"os"
	"io"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/dzhenquan/filesync/config"
	"github.com/dzhenquan/filesync/model"
	"github.com/dzhenquan/filesync/web/router"
)

func main() {
	fmt.Println("gin.Version: ", gin.Version)
	if config.ServerConfig.Env != model.DevelopmentMode {
		// Disable Console Color, you don`t need onsole color when writing the logs the file.
		gin.DisableConsoleColor()

		// Logging to a file.
		logFile, err := os.OpenFile(config.ServerConfig.LogFile, os.O_WRONLY| os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
		gin.DefaultWriter = io.MultiWriter(logFile)
	}

	//Init Model
	db, err := model.InitDB()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	defer db.Close()

	// Creates a router without any middleware by default
	app := gin.New()

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	maxSize := int64(config.ServerConfig.MaxMultipartMemory)
	app.MaxMultipartMemory = maxSize << 20

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	app.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	app.Use(gin.Recovery())

	router.Route(app)

	addr := fmt.Sprintf("%s:%d", config.ServerConfig.WebHost, config.ServerConfig.WebPort)

	app.Run(addr)
}