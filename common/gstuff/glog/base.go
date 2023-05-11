package glog

import (
	"app/common/config"
	"app/utils"
)

var (
	cfg     = config.GetConfig()
	logFile *utils.LogFile
)

func init() {
	logFile = utils.NewLogFile()
}
