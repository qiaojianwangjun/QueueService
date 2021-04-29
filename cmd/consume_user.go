package cmd

// 获取玩家排位请求
type ConsumeUser struct {
}

// 获取玩家排位回复
type ConsumeUserResp struct {
	NowUserId int64 // 当前队首玩家
	QueueLen  int64 // 队列人数
}

// 获取玩家排位回复
type LoginGameNotify struct {
	Token string // 登录游戏token
}
