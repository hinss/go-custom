package services

import (
	"os"

	"github.com/hinss/go-custom/framework"
	"github.com/hinss/go-custom/framework/contract"
)

// CustomConsoleLog 代表控制台输出
type CustomConsoleLog struct {
	CustomLog
}

// NewCustomConsoleLog 实例化CustomConsoleLog
func NewCustomConsoleLog(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)

	log := &CustomConsoleLog{}

	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)

	// 最重要的将内容输出到控制台
	log.SetOutput(os.Stdout)
	log.c = c
	return log, nil
}
