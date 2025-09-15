package base62

import (
	"math"
	"strings"
)

//62进制
// 0-9: 0-9
// a-z: 10-35
// A-Z: 36-61

// const (
//
//	Base62Str = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
//
// )
var (
	base62Str string
)

// MustInit 初始化base62字符串
func MustInit(bs string) {
	if len(bs) == 0 {
		panic("need base string !")
	}
	base62Str = bs
}

func IntToBase62(seq uint64) string {
	if seq == 0 {
		return string(base62Str[0])
	}
	bl := []byte{}
	for seq > 0 {
		mod := seq % 62
		div := seq / 62
		bl = append(bl, base62Str[mod])
		seq = div
	}
	// 反转
	return string(reverse(bl))
}

// String2Int 62进制字符串转10进制
func String2Int(str string) (seq uint64) {
	bl := []byte(str)
	bl = reverse(bl)
	for idx, b := range bl {
		base := math.Pow(62, float64(idx))
		seq += uint64(strings.Index(base62Str, string(b))) * uint64(base)
	}
	return seq
}
func reverse(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
