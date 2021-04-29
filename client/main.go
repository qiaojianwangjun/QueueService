package main

import (
	"QueueService/cmd"
	"QueueService/conn/tcp"
	"QueueService/log"
	"encoding/json"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Client struct {
	UserId      int64
	Url         string
	Wg          *sync.WaitGroup
	RttAllTime  int64
	RttAllCount int64
	AveRtt      float64
	RttMap      sync.Map
}

func NewClient(userId int64, url string, wg *sync.WaitGroup) *Client {
	c := &Client{}
	c.UserId = userId
	c.Url = url
	c.Wg = wg
	c.RttMap = sync.Map{}
	return c
}

func (c *Client) Start() {
	conn, err := net.Dial("tcp", c.Url)
	if err != nil || conn == nil {
		log.Info("connect tcp conn fail", err)
		return
	}

	log.Debug("connect success, userId:", c.UserId)
	c.Wg.Add(1)
	// 链接
	tcpConn, err := tcp.NewTcpConn(&conn)
	if err != nil || tcpConn == nil {
		log.Info("new tcp conn fail", err)
		c.Wg.Done()
		return
	}

	// 发送绑定消息
	c.sendBindUser(c.UserId, tcpConn)
	go func() {
		// 程序退出前关闭连接，释放io资源
		defer func() {
			conn.Close()
			c.Wg.Done()
		}()
		for {
			time.Sleep(5 * time.Second)
			c.sendQueryUserRank(c.UserId, tcpConn)
		}
	}()

	// 接收消息
	go c.RecvMsg(tcpConn)
}

func (c *Client) RecvMsg(conn *tcp.TcpConn) {
	for {
		data, err := conn.ReadMsg()
		if err != nil {
			if strings.Contains(err.Error(), "close") {
				log.Error("接收数据错误, userId:", c.UserId, "err", err.Error())
				conn.Close()
				return
			}

			log.Error("接收数据错误, userId:", c.UserId, "err", err.Error())
			conn.Close()
			return
		}
		c.HandleMsg(data)
	}
}

func (c *Client) HandleMsg(data []byte) {
	log.Debug("HandleMsg, userId:", c.UserId, string(data))
	resp := cmd.RespBase{}
	err := json.Unmarshal(data, &resp)
	if err != nil {
		log.Info("HandleMsg Unmarshal fail, err:", err)
		return
	}
	if resp.ErrorCode != 0 {
		log.Info("HandleMsg fail, userId:", c.UserId, ", cmd:", resp.Cmd, ", ErrorCode:", resp.ErrorCode)
		return
	}

	if resp.Cmd == "BindUserResp" {
		log.Debug("HandleMsg BindUserResp, userId:", c.UserId)
	} else if resp.Cmd == "QueryUserRank" {
		d := cmd.QueryUserRankResp{}
		err := json.Unmarshal(resp.Data, &d)
		if err != nil {
			log.Info("HandleMsg QueryUserRankResp fail, err:", err)
			return
		}
		log.Debug("HandleMsg QueryUserRankResp, userId:", c.UserId, ", data", string((resp.Data)))
	}
	c.calcRtt(resp.Cmd)
}

func (c *Client) calcRtt(cmd string) {
	v, ok := c.RttMap.Load(cmd)
	if !ok {
		return
	}
	startTime, ok := v.(time.Time)
	if !ok {
		return
	}

	rttTime := time.Since(startTime).Nanoseconds() / 1e+03 //ns-->us
	atomic.AddInt64(&c.RttAllTime, rttTime)
	atomic.AddInt64(&c.RttAllCount, 1)
}

// 获取所有接口的平均RTT
func (c *Client) GetAveRtt() int {
	if c.RttAllCount == 0 {
		return 0
	}
	return int(c.RttAllTime / c.RttAllCount)
}

func (c *Client) sendBindUser(userId int64, conn *tcp.TcpConn) {
	data := cmd.BindUser{
		UserId:      userId,
		EncryptType: "sha-1",
	}
	b, _ := json.Marshal(data)
	request := cmd.ReqBase{
		Cmd:    "BindUser",
		UserId: userId,
		Data:   b,
	}
	rb, _ := json.Marshal(request)
	c.RttMap.Store(request.Cmd, time.Now())
	err := conn.WriteMsg(rb)
	if err != nil {
		log.Info("send:", len(rb), err)
		return
	}

	log.Debug("sendBindUser:", userId, len(rb))
}

func (c *Client) sendQueryUserRank(userId int64, conn *tcp.TcpConn) {
	data := cmd.BindUser{
		EncryptType: "sha-1",
	}
	b, _ := json.Marshal(data)
	request := cmd.ReqBase{
		Cmd:  "QueryUserRank",
		Data: b,
	}
	rb, _ := json.Marshal(request)
	c.RttMap.Store(request.Cmd, time.Now())
	err := conn.WriteMsg(rb)
	if err != nil {
		log.Info("send fail:", len(rb), ", err:", err)
		return
	}

	log.Debug("sendQueryUserRank:", userId, len(rb))
}

var benchCount = 10000

func main() {
	wg := &sync.WaitGroup{}
	log.Init("log_client.json")

	log.Info("benchCount:", benchCount)
	cMap := make(map[int]*Client)
	startId := 10001
	endId := startId + benchCount
	for i := startId; i <= endId; i++ {
		c := NewClient(int64(i), "127.0.0.1:8080", wg)
		c.Start()
		cMap[i] = c
	}
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		go func() {
			// 计算间隔时间内的平均rtt
			for range ticker.C {
				allRtt := 0
				allCount := 0
				for _, c := range cMap {
					if c != nil {
						allRtt += c.GetAveRtt()
						allCount += 1
					}
				}
				log.Info("All average RTT:", allRtt/allCount, "us")
			}
		}()
	}()
	wg.Wait()
}
