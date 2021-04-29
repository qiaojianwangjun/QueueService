package server

import (
	"QueueService/cmd"
	conn2 "QueueService/conn"
	"QueueService/log"
	"QueueService/util"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const SysCmdError = "_sys/error"
const SysCmdFatal = "_sys/fatal"

type Session struct {
	Server      *Server
	UId         int64
	conn        conn2.IConn
	ip          string
	BindParam   *cmd.BindUser // 绑定参数
	recChan     chan []byte
	sendChan    chan []byte
	closeSig    *util.CloseSig
	wg          sync.WaitGroup
	connTimeout *time.Timer
	ack         int32
	recIdx      int32
}

// NewSession 新建session
func NewSession(serv *Server, conn conn2.IConn) (sess *Session, err error) {
	defer func() {
		if err != nil {
			log.Error("create session error", "err", err.Error(), "query", conn)
			sess.Close(SysCmdFatal, err.Error())
		}
	}()

	sess = &Session{
		conn:     conn,
		ip:       conn.RemoteAddr(),
		Server:   serv,
		closeSig: util.NewCloseSig(),
	}
	return
}

func (s *Session) Start() {
	var err error
	defer func() {
		if err != nil {
			log.Error("start session error", "err", err.Error())
			s.Close(SysCmdFatal, err.Error())
		}
	}()
	go func() {
		s.Server.Wg.Add(1)
		defer func() {
			s.Server.Wg.Done()
			s.Server.OnSessionClosed(s)
		}()

		s.connTimeout = time.NewTimer(time.Hour)
		s.resetConnTimeout()
		s.recChan = make(chan []byte, 100)
		s.sendChan = make(chan []byte, 100)

		// 接收消息
		s.wg.Add(3)
		go s.rec()
		go s.doRec()
		go s.keepAlive()
		s.wg.Wait()
	}()
}

func (s *Session) rec() {
	defer func() {
		s.wg.Done()
	}()
	for !s.closeSig.Closed {
		data, err := s.conn.ReadMsg()
		if s.closeSig.Closed {
			return
		}
		if err != nil {
			if strings.Contains(err.Error(), "close") {
				log.Info("rec fail", "err", err.Error())
				s.Close(SysCmdError, "rec msg close")
			} else {
				log.Info("rec fail", "err", err.Error())
				s.Close(SysCmdError, "rec msg error")
			}
			return
		}
		s.resetConnTimeout()

		// qps
		s.Server.Qps.AddIn()
		util.Try(func() {
			s.recChan <- data
		}, nil)
	}
}

func (s *Session) decode(bytes []byte) (out []byte, err error) {
	//bytes, err = s.encryptor.Decrypt(bytes)
	//if err != nil {
	//	return
	//}
	out = bytes
	return
}

func (s *Session) doAck() (err error) {
	return
}

func (s *Session) doRec() {
	defer func() {
		close(s.recChan)
		s.wg.Done()
	}()
	for !s.closeSig.Closed {
		select {
		case <-s.closeSig.ClosingSig:
			return
		case bytes, ok := <-s.recChan:
			if !ok {
				return
			}
			// 简化的心跳
			if len(bytes) <= 1 {
				err := s.doAck()
				if err != nil {
					return
				}
				continue
			}

			recMsg, err := s.decode(bytes)
			if err != nil {
				s.Close(SysCmdError, "decode recMsg err")
				return
			}
			log.Debug("doRec recMsg len", len(recMsg))
			msg := cmd.ReqBase{}
			err = json.Unmarshal(recMsg, &msg)
			if err != nil {
				log.Info("doRec Unmarshal fail", err)
				s.Close(SysCmdError, "doRec Unmarshal fail")
				return
			}

			log.Debug("doRec msg:", msg)
			err = s.HandleCmd(msg.Cmd, msg.Data)
			if err != nil {
				log.Info("doRec HandleCmd fail", err)
				s.Close(SysCmdError, "doRec HandleCmd fail")
				return
			}
		}
	}
}

// HandleCmd 处理接收到的命令
func (s *Session) HandleCmd(cmd string, data []byte) error {
	handler, ok := HandlerMap[cmd]
	if !ok {
		log.Info("no this cmd:", cmd, "cmdMap:", HandlerMap)
		return errors.New("no this cmd " + cmd)
	}

	// 反射构造第二个参数（请求数据）
	req := reflect.New(handler.Method.Type.In(2).Elem())
	err := json.Unmarshal(data, req.Interface())
	if err != nil {
		log.Info("HandleCmd parse fail", err)
		return err
	}

	// 进行调用
	params := make([]reflect.Value, 0, 3)
	params = append(params, handler.Logic)
	params = append(params, reflect.ValueOf(s))
	params = append(params, req)
	values := handler.Method.Func.Call(params)
	if len(values) != 3 {
		log.Info("cmd handler error:", cmd, "values:", values)
		return errors.New("cmd handler error, cmd:" + cmd)
	}
	err = s.SendResp(cmd, int(values[0].Int()), values[1].String(), values[2])
	if err != nil {
		log.Info("cmd handler send fail:", cmd, "values:", values)
		return errors.New("cmd handler send fail, cmd:" + cmd)
	}
	return nil
}

// 重置时间
func (s *Session) resetConnTimeout() {
	s.connTimeout.Reset(time.Second * time.Duration(15))
}

func (s *Session) keepAlive() {
	defer func() {
		s.wg.Done()
	}()
	for !s.closeSig.Closed {
		select {
		case <-s.closeSig.ClosingSig:
			return
		case <-s.connTimeout.C:
			log.Debug("keepAlive timeout")
			//s.Close(SysCmdError, "keepAlive timeout")
			return
		}
	}
}

// EncodeMsg 编码数据
func (s *Session) EncodeMsg(cmdStr string, errCode int, message string, data interface{}) (respB []byte, err error) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Info("DoSend Marshal fail", err)
		return
	}
	resp := cmd.RespBase{
		Cmd:       cmdStr,
		ErrorCode: errCode,
		Message:   message,
		Data:      b,
	}
	respB, err = json.Marshal(resp)
	if err != nil {
		log.Info("DoSend Marshal fail", err)
		return
	}
	return
}

// SendResp 发送数据
func (s *Session) SendResp(cmdStr string, errCode int, message string, data interface{}) (err error) {
	respB, err := s.EncodeMsg(cmdStr, errCode, message, data)
	if err != nil {
		log.Info("DoSend Marshal fail", err)
		return
	}
	// 发送数据
	err = s.conn.WriteMsg(respB)
	if err != nil {
		log.Info("DoSend Marshal fail", err)
		return
	}
	// qps处理
	s.Server.Qps.AddOut()
	return
}

func (s *Session) Close(cmd string, reason string) {
	if !s.closeSig.Close() {
		return
	}
	if s.conn != nil {
		atomic.AddInt64(&s.Server.ConnCount, -1)
		s.conn.Close()
	}
	log.Info("close session", "cmd", cmd, "reason", reason)
}
