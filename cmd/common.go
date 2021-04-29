package cmd

// 请求数据基础结构
type ReqBase struct {
	Cmd     string // 消息命令
	Version uint8  // 版本号
	UserId  int64  // 用户ID
	Data    []byte // 数据
}

// 请求回复数据基础结构
type RespBase struct {
	Cmd       string // 消息命令
	Version   uint8  // 版本号
	ErrorCode int    // 错误码
	Message   string // 描述信息
	Data      []byte // 数据
}
