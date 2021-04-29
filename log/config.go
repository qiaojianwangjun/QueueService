package log

import (
	"QueueService/config"
	"QueueService/env"
	"fmt"
	"os"
	"strings"
)

type LogConfig struct {
	Level    Level  `json:"level"`    // 日志等级,debug,info,waring,error,fatal
	Filename string `json:"filename"` // 日志文件名（含路径）
	Service  string `json:"service"`  // 服务名称
}

var Config = &LogConfig{
	Level:    DebugLevel,
	Filename: "./logs/log.txt",
}

func Init(configName string) {
	envValue := env.GetEnv()
	fileName := fmt.Sprintf("./config/%s/%s", envValue, configName)

	err := config.LoadConfig(fileName, Config)
	if err != nil {
		fmt.Println("load log config fail", err)
		return
	}
	hostname, _ := os.Hostname()
	Config.Filename = strings.ReplaceAll(Config.Filename, "{hostname}", hostname)
	Config.Filename = strings.ReplaceAll(Config.Filename, "{service}", Config.Service)

	lg := NewLogger(Config)
	var logger Logger = &lg
	SetLogger(logger)
}
