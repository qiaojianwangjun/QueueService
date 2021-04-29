package cmd

// 获取玩家排位请求
type QueryUserRank struct {
}

// 获取玩家排位回复
type QueryUserRankResp struct {
	QueueLen int64 // 队列人数
	Rank     int64 // 当前排位
}
