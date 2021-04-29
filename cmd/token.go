package cmd

// 通知玩家进行登录游戏
type ObtainTokenResp struct {
	Token string // 获取到token，用于游戏服务验证
}
