package server

import (
	"QueueService/config"
	"QueueService/log"
	"QueueService/util"
	"fmt"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type Server struct {
	ServerId    string
	Ip          string
	HostName    string
	Config      *config.Config
	Sessions    SessionMap // session管理
	CloseSig    *util.CloseSig
	Wg          sync.WaitGroup
	Qps         *util.InOutQPS
	ConnCount   int64
	ConnChannel chan net.Conn
}

func NewServer(conf *config.Config) *Server {
	server := &Server{
		Config:      conf,
		CloseSig:    util.NewCloseSig(),
		Qps:         util.NewInOutQPS(time.Second * 5),
		ConnChannel: make(chan net.Conn),
	}
	server.ServerId = fmt.Sprintf("agent.%s", util.GetMacAddr())
	server.Ip = util.GetPrivateIPv4()
	server.HostName, _ = os.Hostname()
	server.Qps.Start()
	return server
}

// Start 服务启动监听
func (s *Server) Start() {
	listener, err := net.ListenTCP("tcp4", &net.TCPAddr{Port: s.Config.Port})
	if err != nil {
		panic("listen fail err:" + err.Error())
		return
	}
	log.Info("server start port:", s.Config.Port)
	go HandleConn(s)
	s.accept(listener)
}

func (s *Server) accept(listener *net.TCPListener) {
	for {
		connection, err := listener.AcceptTCP()
		if err != nil {
			log.Info("accept fail:", err.Error())
		} else {
			log.Info("accept success:", connection.RemoteAddr().String())
			atomic.AddInt64(&s.ConnCount, 1)
			s.ConnChannel <- connection
		}
	}
}

func (s *Server) OnSessionClosed(sess *Session) {
	s.Sessions.Delete(sess.UId, sess)
	atomic.AddInt64(&s.ConnCount, -1)
}
