package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
)

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GetMD5Hash(s string) string {
	hasher := md5.New()
    hasher.Write([]byte(s))
    return hex.EncodeToString(hasher.Sum(nil))
}