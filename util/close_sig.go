package util

import (
	"sync"
)

type CloseSig struct {
	Closed      bool
	ClosingSig  chan interface{} // 正在关闭的信号
	closingOnce sync.Once
	ClosedSig   chan struct{} // 关闭完成的信号
	ClosedOnce  sync.Once
}

func NewCloseSig() *CloseSig {
	c := new(CloseSig)
	c.ClosingSig = make(chan interface{}, 1)
	c.ClosedSig = make(chan struct{})
	return c
}

func (c *CloseSig) Close() bool {
	if c.Closed {
		return false
	}
	change := false
	c.closingOnce.Do(func() {
		close(c.ClosingSig)
		c.Closed = true
		change = true
	})
	return change
}

// CloseWithData 有数据的关闭
func (c *CloseSig) CloseWithData(data interface{}) bool {
	if c.Closed {
		return false
	}
	change := false
	c.closingOnce.Do(func() {
		c.ClosingSig <- data
		close(c.ClosingSig)
		c.Closed = true
		change = true
	})
	return change
}

func (c *CloseSig) Done() {
	c.ClosedOnce.Do(func() {
		close(c.ClosedSig)
	})
}

func (c *CloseSig) WaitForClosing() {
	<-c.ClosingSig
}

func (c *CloseSig) WaitForClosed() {
	<-c.ClosedSig
}
