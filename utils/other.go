package utils

import (
	"crypto/md5"
	crypto "crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
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

func GenerateRandomCode(length int) (string, error) {
    var code string
    for i := 0; i < length; i++ {
        n, err := crypto.Int(crypto.Reader, big.NewInt(10))
        if err != nil {
            return "", err
        }
        code += fmt.Sprintf("%d", n.Int64())
    }
    return code, nil
}