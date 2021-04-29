package server

import (
	"QueueService/conn"
	"QueueService/conn/tcp"
	"QueueService/log"
	"sync/atomic"
)

// 处理连接
func HandleConn(s *Server) {
	log.Info("HandleConn connection ....")
	for {
		select {
		case c := <-s.ConnChannel:
			log.Info("HandleConn " + c.RemoteAddr().String() + " connected")
			tcpConn, err := tcp.NewTcpConn(&c)
			if err == nil {
				addSession(s, tcpConn)
			} else {
				log.Info("HandleConn fail:", tcpConn, "err:", err)
			}
		}
	}
}

// 新建并添加session
func addSession(s *Server, c conn.IConn) {
	sess, err := NewSession(s, c)
	if err != nil {
		log.Info("addSession fail", "err:", err)
		return
	}
	atomic.AddInt64(&s.ConnCount, 1)
	sess.Start()
}
