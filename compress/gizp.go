package compress

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"sync"
)

type gzw struct {
	*gzip.Writer
	buf *bytes.Buffer
}

// gzip写的池
var gzwPool = &sync.Pool{
	New: func() interface{} {
		gz := &gzw{
			Writer: gzip.NewWriter(nil),
			buf:    bytes.NewBuffer(make([]byte, 1024*10)), // 默认10k
		}
		return gz
	},
}

// gzip读的池
var gzrPool = &sync.Pool{
	New: func() interface{} {
		gz := &gzip.Reader{}
		return gz
	},
}

// gzip压缩/解压接口实现
type gzipCompress struct {
}

// Compress 压缩接口
func (*gzipCompress) Compress(in []byte) (out []byte, err error) {
	gz := gzwPool.Get().(*gzw)
	defer gzwPool.Put(gz)
	gz.buf.Reset()
	gz.Reset(gz.buf)
	_, err = gz.Write(in)
	if err != nil {
		gz.Close()
		return
	}
	err = gz.Close()
	if err != nil {
		return
	}
	return ioutil.ReadAll(gz.buf)
}

// UnCompress 压缩接口
func (*gzipCompress) UnCompress(in []byte) (out []byte, err error) {
	gz := gzrPool.Get().(*gzip.Reader)
	defer gzrPool.Put(gz)
	err = gz.Reset(bytes.NewReader(in))
	if err != nil {
		return
	}
	defer gz.Close()
	return ioutil.ReadAll(gz)
}

var GzipCompress = &gzipCompress{}
