package md5

import (
	"crypto/md5"
	"encoding/hex"
)

// Sum 计算md5
func Sum(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
