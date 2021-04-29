package util

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}

func JsonEncodeToStringMust(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func RSAEncrypt(plaintext string, publickey *rsa.PublicKey) (string, error) {
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publickey, []byte(plaintext))
	if err != nil {
		return "", err
	}
	decodedtext := base64.StdEncoding.EncodeToString(ciphertext)
	return decodedtext, nil
}

func RSAEncryptRaw(data []byte, publickey *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, publickey, data)
}

func RSAEncryptObjRaw(obj interface{}, pubKey *rsa.PublicKey) string {
	info := JsonEncodeToStringMust(obj)
	partLen := pubKey.N.BitLen()/8 - 11
	data := []byte(info)
	if len(data) <= partLen {
		str, err := RSAEncrypt(info, pubKey)
		if err != nil {
			return ""
		}
		return str
	}
	chunks := split([]byte(info), partLen)
	buffer := bytes.Buffer{}
	for _, chunk := range chunks {
		b, err := RSAEncryptRaw(chunk, pubKey)
		if err != nil {
			return ""
		}
		buffer.Write(b)
	}
	return base64.StdEncoding.EncodeToString(buffer.Bytes())
}

func RSALoadPublicKeyBase64(base64key string) (*rsa.PublicKey, error) {
	keybytes, err := base64.StdEncoding.DecodeString(base64key)
	if err != nil {
		return nil, fmt.Errorf("base64 decode failed, error=%s\n", err.Error())
	}

	pubkeyinterface, err := x509.ParsePKIXPublicKey(keybytes)
	if err != nil {
		return nil, err
	}

	publickey := pubkeyinterface.(*rsa.PublicKey)
	return publickey, nil
}

func RSALoadPrivateKeyBase64(base64key string) (*rsa.PrivateKey, error) {
	keybytes, err := base64.StdEncoding.DecodeString(base64key)
	if err != nil {
		return nil, fmt.Errorf("base64 decode failed, error=%s\n", err.Error())
	}

	privatekey, err := x509.ParsePKCS1PrivateKey(keybytes)
	if err != nil {
		return nil, errors.New("parse private key error!")
	}

	return privatekey, nil
}
