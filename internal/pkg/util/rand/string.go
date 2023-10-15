package rand

import (
	"crypto/rand"
	"encoding/hex"
)

func GenHex(l uint64) string {
	bytes := make([]byte, l)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
