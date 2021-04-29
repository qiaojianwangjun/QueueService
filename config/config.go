package config

import (
	"QueueService/util"
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
)

type Pprof struct {
}

type Config struct {
	Service       string          `json:"service"`     // 服务名称
	Port          int             `json:"port"`        // 端口号
	MaxQueueCnt   int64           `json:"maxQueueCnt"` // 队列最大人数
	PProfPort     int             `json:"pprofPort"`   // pprof端口号
	PublicKey     string          `json:"publicKey"`   // rsa的公钥
	PrivateKey    string          `json:"privateKey"`  // rsa的私钥
	PublicKeyRSA  *rsa.PublicKey  `json:"-"`
	PrivateKeyRSA *rsa.PrivateKey `json:"-"`
}

// 读取配置
var cfg Config

func LoadConfig(filename string, result interface{}) error {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, result)
	if err != nil {
		return err
	}

	return nil
}

func SetConfig(c Config) {
	c.PublicKeyRSA, _ = util.RSALoadPublicKeyBase64(c.PublicKey)
	c.PrivateKeyRSA, _ = util.RSALoadPrivateKeyBase64(c.PrivateKey)
	cfg = c
}

func GetConfig() Config {
	return cfg
}
