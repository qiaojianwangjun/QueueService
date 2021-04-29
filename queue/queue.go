package queue

type IQueue interface {
	// 压入一个
	PushOne(userId int64) (int64, error)
	// 弹出一个
	PopOne() (int64, error)
	// 获取用户排队位置
	GetUserRank(userId int64) (int64, error)
	// 获取总用户数
	GetAllCount() (int64, error)
}
