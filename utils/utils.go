package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

// Md5Encrypt md5加密
func Md5Encrypt(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// GenerateToken 生成token
func GenerateToken() string {
	return Md5Encrypt(uuid.New().String())
}

// RandString 生成随机字符串
func RandString(n int) string {
	const letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
