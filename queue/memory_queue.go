package queue

import (
	"container/list"
	"errors"
	"sync"
	"sync/atomic"
)

type CacheData struct {
	UId         int64
	OrderNumber uint64
}

// 内存实现的队列
type MemoryQueue struct {
	CacheMap    sync.Map   // 用户map
	CacheList   *list.List // 用户队列
	CurIncNum   uint64     // 当前号码
	IncOrderNum uint64     // 自增号码
}

func NewMemoryQueue() *MemoryQueue {
	return &MemoryQueue{
		CacheList: list.New(),
	}
}

// calcUserRank 计算用户位置
func (m *MemoryQueue) calcUserRank(orderNumber uint64) (rank uint64) {
	return orderNumber - (m.IncOrderNum - uint64(m.CacheList.Len()))
}

// calcUserRank 计算用户位置
func (m *MemoryQueue) incNum() (rank uint64) {
	return atomic.AddUint64(&m.IncOrderNum, 1)
}

// Compress 压缩接口
func (m *MemoryQueue) PushOne(UId int64) (int64, error) {
	v, ok := m.CacheMap.Load(UId)
	// 已经存在
	if ok {
		data, ok := v.(CacheData)
		if !ok {
			return 0, errors.New("parse fail")
		}
		rank := int64(m.calcUserRank(data.OrderNumber))
		return rank, nil
	}

	data := CacheData{
		UId:         UId,
		OrderNumber: m.incNum(),
	}
	m.CacheMap.Store(UId, data)
	m.CacheList.PushBack(data)
	return int64(m.calcUserRank(data.OrderNumber)), nil
}

// PopOne 弹出一个
func (m *MemoryQueue) PopOne() (int64, error) {
	atomic.AddUint64(&m.CurIncNum, 1)
	e := m.CacheList.Front()
	if e == nil {
		return 0, errors.New("list empty")
	}
	data, ok := (e.Value).(CacheData)
	if !ok {
		return 0, errors.New("parse fail")
	}
	m.CacheList.Remove(e)
	m.CacheMap.Delete(data.UId)
	return data.UId, nil
}

// GetUserRank 获取用户排队位置
func (m *MemoryQueue) GetUserRank(UId int64) (int64, error) {
	v, ok := m.CacheMap.Load(UId)
	// 已经存在
	if ok {
		data, ok := v.(CacheData)
		if !ok {
			return 0, errors.New("parse fail")
		}
		rank := int64(m.calcUserRank(data.OrderNumber))
		return rank, nil
	}
	return 0, errors.New("no this user")
}

// 获取总用户数
func (m *MemoryQueue) GetAllCount() (int64, error) {
	return int64(m.CacheList.Len()), nil
}
