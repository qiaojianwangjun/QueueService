package util

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	NUmStr  = "0123456789"
	CharStr = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	SpecStr = "+=-@#~,.[]()!%^*$"
)

type CharSetType string

const (
	CharSetTypeNum     CharSetType = "num"
	CharSetTypeChar    CharSetType = "char"
	CharSetTypeMix     CharSetType = "mix"
	CharSetTypeAdvance CharSetType = "advance"
)

func GeneratePasswd(length int, charset CharSetType) string {
	// 初始化密码切片
	var passwd = make([]byte, length, length)
	//源字符串
	var sourceStr string
	// 判断字符类型,如果是数字
	if charset == CharSetTypeNum {
		sourceStr = NUmStr
		// 如果选的是字符
	} else if charset == CharSetTypeChar {
		sourceStr = CharStr
		// 如果选的是混合模式
	} else if charset == CharSetTypeMix {
		sourceStr = fmt.Sprintf("%s%s", NUmStr, CharStr)
		// 如果选的是高级模式
	} else if charset == CharSetTypeAdvance {
		sourceStr = fmt.Sprintf("%s%s%s", NUmStr, CharStr, SpecStr)
	} else {
		sourceStr = NUmStr
	}
	// 遍历，生成一个随机index索引
	charLen := len(sourceStr)
	// 如果seed固定，那么每次程序重启后重新生成随机数会重复上一次的随机数
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		index := rand.Intn(charLen)
		passwd[i] = sourceStr[index]
	}
	return string(passwd)
}
