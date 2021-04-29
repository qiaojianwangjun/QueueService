package util

import (
	"crypto/rsa"
	"time"
)

type tokenData struct {
	Info      interface{} `json:"info"`
	ExpiredAt int64       `json:"expired"`
}

func CreateToken(info interface{}, expiredIn time.Duration, publicKeyRSA *rsa.PublicKey) string {
	return RSAEncryptObjRaw(&tokenData{
		Info:      info,
		ExpiredAt: time.Now().Add(expiredIn).Unix(),
	}, publicKeyRSA)
}
