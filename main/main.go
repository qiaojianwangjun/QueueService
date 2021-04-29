package main

import (
	"QueueService/config"
	"QueueService/env"
	"QueueService/handler"
	"QueueService/log"
	"QueueService/server"
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	// 读取配置
	var cfg config.Config
	envValue := env.GetEnv()
	log.Init("log_server.json")
	fileName := fmt.Sprintf("./config/%s/app.json", envValue)
	err := config.LoadConfig(fileName, &cfg)
	if err != nil {
		fmt.Println("load app config fail", err)
		return
	}
	config.SetConfig(cfg)
	// 检测是否开启性能监控
	if cfg.PProfPort > 0 {
		go func() {
			log.Info("start ListenAndServe pprof port:", cfg.PProfPort)
			err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.PProfPort), nil)
			if err != nil {
				log.Error("ListenAndServe pprof fail", err)
				return
			}
		}()
	}
	// 注册处理逻辑
	server.Register(handler.NewHandler())
	// 启动服务
	s := server.NewServer(&cfg)
	s.Start()
}
