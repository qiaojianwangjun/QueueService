package queue

import (
	"testing"
)

//go test -v memory_queue.go memory_queue_test.go

func TestQueue(t *testing.T) {
	queue := NewMemoryQueue()

	// 测试压入
	for i := 1; i <= 100; i++ {
		n, err := queue.PushOne(int64(i))
		if err != nil || n != int64(i) {
			t.Error("PushOne fail")
		}
	}

	// 测试获取
	if cnt, err := queue.GetAllCount(); err != nil || cnt != 100 {
		t.Error("GetAllCount fail")
	}
	if cnt, err := queue.GetUserRank(20); err != nil || cnt != 20 {
		t.Error("GetUserRank fail")
	}
	if cnt, err := queue.GetUserRank(80); err != nil || cnt != 80 {
		t.Error("GetUserRank fail")
	}

	// 测试弹出
	for i := 1; i <= 50; i++ {
		id, err := queue.PopOne()
		if err != nil || id != int64(i) {
			t.Error("PopOne fail")
		}
	}

	// 测试获取
	if cnt, err := queue.GetAllCount(); err != nil || cnt != 50 {
		t.Error("GetAllCount fail")
	}
	if cnt, err := queue.GetUserRank(20); err == nil || cnt != 0 {
		t.Error("GetUserRank fail")
	}
	if cnt, err := queue.GetUserRank(70); err != nil || cnt != 20 {
		t.Error("GetUserRank fail")
	}
	if cnt, err := queue.GetUserRank(75); err != nil || cnt != 25 {
		t.Error("GetUserRank fail")
	}
	if cnt, err := queue.GetUserRank(110); err == nil || cnt != 0 {
		t.Error("GetUserRank fail")
	}

	// 测试存储的结果
	cnt := 0
	queue.CacheMap.Range(func(key, value interface{}) bool {
		cnt++
		return true
	})
	if cnt != 50 {
		t.Error("pop fail")
	}

	// 测试存储的结果
	if queue.CacheList.Len() != 50 {
		t.Error("pop fail")
	}
	if queue.CurIncNum != 50 {
		t.Error("CurIncNum fail")
	}
	if queue.IncOrderNum != 100 {
		t.Error("IncOrderNum fail")
	}
}
