package server

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"
)

func HashStr(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

func RandomNum(length int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(length)
}
