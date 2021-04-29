package conn

// IConn 链接接口
type IConn interface {
	RemoteAddr() string
	ReadMsg() (data []byte, err error)
	WriteMsg(data []byte) (err error)
	Close()
}
