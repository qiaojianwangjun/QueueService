package cmd

// 绑定用户，协商数据
type BindUser struct {
	UserId       int64  // 测试时用，正常通过token上第三方获取相应数据
	Token        string // sdk等登录的token，用于验证用户
	EncryptType  string // 可用于协商加密类型
	CompressType string // 可用于协商压缩类型
	EncodeType   string // 可用于协商编码类型
}

type BindUserResp struct {
}
