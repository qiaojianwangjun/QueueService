package tcp

import (
	"QueueService/log"
	"net"
	"sync"
	"time"
	"unsafe"
)

// 数据长度字节数
const dataLenByteNum = 4

// 每次读取的数据长度
const batchSize = 10 * 1024

type TcpConn struct {
	net.Conn
	header         [dataLenByteNum]byte
	data           []byte
	headerPtr      *int32
	closeOnce      sync.Once
	timeout        *time.Timer
	idleTimout     time.Duration
	writeHeader    [dataLenByteNum]byte
	writeHeaderPtr *int32
	writeMu        sync.Mutex
	remoteAddr     string
}

// 创建一个tcp链接结构
func NewTcpConn(tcpConn *net.Conn) (conn *TcpConn, err error) {
	ip := (*tcpConn).RemoteAddr().String()
	conn = &TcpConn{
		Conn:       *tcpConn,
		idleTimout: time.Minute,
		remoteAddr: ip,
	}
	// 接收buf头指针
	conn.headerPtr = (*int32)(unsafe.Pointer(&conn.header))
	// 发送buf头指针
	conn.writeHeaderPtr = (*int32)(unsafe.Pointer(&conn.writeHeader))
	conn.timeout = time.NewTimer(conn.idleTimout)
	go func() {
		<-conn.timeout.C
		log.Info("conn timeout")
		conn.Close()
	}()
	return
}

func (c *TcpConn) SetIdleTimout(dur time.Duration) {
	c.idleTimout = dur
	c.timeout.Reset(c.idleTimout)
}

func (c *TcpConn) RemoteAddr() string {
	return c.remoteAddr
}

func (c *TcpConn) readMsg() (data []byte, err error) {
	idx := 0
	var n int
	// 先读取4个字节的数据长度，防止粘包问题
	for idx < dataLenByteNum {
		n, err = c.Read(c.header[idx:])
		if err != nil {
			log.Info("readMsg fail", "n", n, "err:", err)
			return
		}
		// 重新记时
		c.timeout.Reset(c.idleTimout)
		idx += n
	}
	headLen := int(*c.headerPtr)
	data = make([]byte, headLen)
	idx = 0
	idx2 := 0

	log.Debug("readMsg headLen:", headLen, len(data))
	// 读取固定长度的数据，每次读取
	for idx < headLen {
		idx2 = idx + batchSize
		if idx2 > headLen {
			idx2 = headLen
		}
		n, err = c.Read(data[idx:idx2])
		if err != nil {
			log.Info("readMsg err", n, err)
			return
		}
		// 重新记时
		c.timeout.Reset(c.idleTimout)
		idx += n
	}
	log.Debug("readMsg data:", data)
	return
}

// ReadMsg 读取数据
func (c *TcpConn) ReadMsg() (data []byte, err error) {
	data, err = c.readMsg()
	if err != nil {
		log.Info("ReadMsg fail", "err", err)
		c.Close()
		return
	}
	return
}

func (c *TcpConn) writeMsg(data []byte) (err error) {
	dataLen := len(data)
	idx := 0
	var n int
	for idx < dataLen {
		n, err = c.Write(data[idx:])
		log.Debug("writeMsg n:", n, string(data))
		if err != nil {
			return
		}
		idx += n
	}
	return
}

// ReadMsg 发送数据
func (c *TcpConn) WriteMsg(data []byte) (err error) {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	// 先写数据长度
	*c.writeHeaderPtr = int32(len(data))
	err = c.writeMsg(c.writeHeader[:])

	if err != nil {
		c.Close()
		return
	}
	// 发送数据
	err = c.writeMsg(data)
	if err != nil {
		c.Close()
		return
	}
	return
}

func (c *TcpConn) Close() {
	c.closeOnce.Do(func() {
		c.Conn.Close()
	})
}
