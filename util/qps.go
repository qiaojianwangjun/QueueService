package util

import (
	"fmt"
	"sync/atomic"
	"time"
)

type QPS struct {
	count    int64
	interval time.Duration
	ticker   *time.Ticker
	qps      float64
}

// NewQPS 新建qps，按时间间隔进行平均
func NewQPS(interval time.Duration) *QPS {
	return &QPS{
		interval: interval,
	}
}

func (qps *QPS) Start() {
	qps.Stop()
	qps.ticker = time.NewTicker(qps.interval)
	inter := qps.interval.Seconds()
	go func() {
		// 计算间隔时间内的平均qps
		for range qps.ticker.C {
			qps.qps = float64(qps.count) / inter
			atomic.StoreInt64(&qps.count, 0)
		}
	}()
}

func (qps *QPS) Stop() {
	if qps.ticker != nil {
		qps.ticker.Stop()
	}
}

// Add 添加qps计数
func (qps *QPS) Add() {
	atomic.AddInt64(&qps.count, 1)
}

// 获取qps
func (qps *QPS) QPS() int {
	return int(qps.qps)
}

type InOutQPS struct {
	inCount  int64
	outCount int64
	interval time.Duration
	ticker   *time.Ticker
	inQps    float64
	outQps   float64
}

// NewInOutQPS 新建请求与响应不相同qps，比如长链接有推送情况，按时间间隔进行平均
func NewInOutQPS(interval time.Duration) *InOutQPS {
	return &InOutQPS{
		interval: interval,
	}
}

func (qps *InOutQPS) Start() {
	qps.Stop()
	qps.ticker = time.NewTicker(qps.interval)
	inter := qps.interval.Seconds()
	go func() {
		for range qps.ticker.C {
			qps.inQps = float64(qps.inCount) / inter
			atomic.StoreInt64(&qps.inCount, 0)
			fmt.Println("QPS inQps:", qps.inQps)
			qps.outQps = float64(qps.outCount) / inter
			atomic.StoreInt64(&qps.outCount, 0)
			fmt.Println("QPS outQps:", qps.outQps)
		}
	}()
}

func (qps *InOutQPS) Stop() {
	if qps.ticker != nil {
		qps.ticker.Stop()
	}
}

// AddIn 增加请求计数
func (qps *InOutQPS) AddIn() {
	atomic.AddInt64(&qps.inCount, 1)
}

// AddOut 增加响应计数
func (qps *InOutQPS) AddOut() {
	atomic.AddInt64(&qps.outCount, 1)
}

func (qps *InOutQPS) InQPS() int {
	return int(qps.inQps)
}

func (qps *InOutQPS) OutQPS() int {
	return int(qps.outQps)
}
