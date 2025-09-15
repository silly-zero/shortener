package base62

//62进制
// 0-9: 0-9
// a-z: 10-35
// A-Z: 36-61

const (
	Base62Str = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func IntToBase62(seq uint64) string {
	if seq == 0 {
		return string(Base62Str[0])
	}
	bl := []byte{}
	for seq > 0 {
		mod := seq % 62
		div := seq / 62
		bl = append(bl, Base62Str[mod])
		seq = div
	}
	// 反转
	return string(reverse(bl))
}

func reverse(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
