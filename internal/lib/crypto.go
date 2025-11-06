package lib

import (
	"crypto/rand"
	"encoding/hex"
	mathRand "math/rand"
	"time"
)

func RandomString(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		mathRand := mathRand.New(mathRand.NewSource(int64(time.Now().UnixNano())))
		for i := range b {
			b[i] = byte(mathRand.Intn(256))
		}
	}
	return hex.EncodeToString(b)
}
