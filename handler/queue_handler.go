package handler

import (
	"QueueService/cmd"
	"QueueService/config"
	"QueueService/log"
	"QueueService/queue"
	"QueueService/server"
	"QueueService/util"
	"time"
)

type Handler struct {
	Queue *queue.MemoryQueue
}

func NewHandler() *Handler {
	h := &Handler{}
	h.Queue = queue.NewMemoryQueue(uint64(config.GetConfig().MaxQueueCnt))
	return h
}

/*
 绑定链接
*/
func (h *Handler) BindUser(sess *server.Session, req *cmd.BindUser) (code int, message string, respData interface{}) {
	log.Debug("BindUser", req)
	// TODO 验证用户真实性
	password := util.GeneratePasswd(16, util.CharSetTypeMix)
	sess.UId = req.UserId
	sess.BindParam = req
	sess.Password = []byte(password)

	// 检测是否已满
	if h.Queue.IsFull() {
		code = -1
		return
	}

	old := sess.Server.Sessions.StoreOrUpdate(sess.UId, sess)
	if old != nil {
		log.Info("BindUser other login userId", req.UserId, sess.UId)
		old.Close("SysCmdFatal", "other login")
		code = -1
		return
	}

	rank, err := h.Queue.PushOne(sess.UId)
	if err != nil {
		code = -1
		return
	}
	queueLen, _ := h.Queue.GetAllCount()
	respRank := cmd.QueryUserRankResp{
		QueueLen: queueLen,
		Rank:     rank,
	}
	defer func() {
		sess.SendResp("QueryUserRank", 0, "", respRank)
	}()
	respData = cmd.BindUserResp{
		EncryptType: "",
		Password:    password,
	}
	return
}

/*
 查询排名
*/
func (h *Handler) QueryUserRank(sess *server.Session, req *cmd.QueryUserRank) (code int, message string, respData interface{}) {
	log.Debug("QueryUserRank", req)
	rank, err := h.Queue.PushOne(sess.UId)
	if err != nil {
		code = -1
		return
	}
	queueLen, _ := h.Queue.GetAllCount()
	respData = cmd.QueryUserRankResp{
		QueueLen: queueLen,
		Rank:     rank,
	}
	return
}

/*
 消耗玩家
*/
func (h *Handler) ConsumeUser(sess *server.Session, req *cmd.ConsumeUser) (code int, message string, respData interface{}) {
	log.Debug("ConsumeUser", req)
	var userId int64
	var err error
	for {
		userId, err = h.Queue.PopOne()
		if err != nil {
			code = -1
			return
		}
		// 查找在线的玩家，不在线玩家跳过
		v, ok := sess.Server.Sessions.Load(userId)
		if !ok {
			continue
		}
		targetSess, ok := v.(*server.Session)
		if !ok {
			continue
		}

		type userInfo struct {
			UserId int64
		}

		// 通知玩家可以登录
		// TODO token生成与持久化供其他服务使用或解密
		notify := cmd.LoginGameNotify{
			Token: util.CreateToken(userInfo{UserId: userId}, 30*time.Second, config.GetConfig().PublicKeyRSA),
		}
		targetSess.SendResp("LoginGameNotify", 0, "", notify)
		break
	}

	queueLen, _ := h.Queue.GetAllCount()
	respData = cmd.ConsumeUserResp{
		NowUserId: userId,
		QueueLen:  queueLen,
	}

	return
}
