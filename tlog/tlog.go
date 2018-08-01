package tlog

import (
	"github.com/dzhenquan/dlog"
	"github.com/dzhenquan/filesync/config"
)

var Tlog *dlog.DLogger

func init() {
	dlog := dlog.NewDLogger(config.ServerConfig.LogDir,
					config.ServerConfig.LogOldDir,
					config.ServerConfig.LogPrefix,
					config.ServerConfig.LogSuffix)

	dlog.SetMaxLine(config.ServerConfig.LogMaxLine)
	dlog.SetMaxByte(config.ServerConfig.LogMaxByte)

	Tlog = dlog
}