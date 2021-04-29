package main

import (
	"QueueService/log"
	"QueueService/util"
	"os"
	"os/signal"
	"syscall"
)

var exitSig = make(chan struct{})
var stopSig = make(chan struct{})

func init() {
	// 处理系统信号
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	sig2 := make(chan os.Signal, 1)
	signal.Notify(sig2, syscall.Signal(10))
	go func() {
		<-sig
		log.Info("exit sig")
		util.Try(func() {
			close(exitSig)
		}, nil)
	}()
	go func() {
		<-sig2
		log.Info("stop sig")
		util.Try(func() {
			close(stopSig)
		}, nil)
	}()
}

// WaitForExit 等待终止信号
func WaitForExit() {
	<-exitSig
}

// WaitForStop 等待停止信号(graceful stop)
func WaitForStop() {
	<-stopSig
}

func ExitSig() chan struct{} {
	return exitSig
}

func StopSig() chan struct{} {
	return stopSig
}

func Stop() {
	util.Try(func() {
		close(stopSig)
	}, nil)
}

func Exit() {
	util.Try(func() {
		close(exitSig)
	}, nil)
}
