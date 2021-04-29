package compress

type ICompress interface {
	// 压缩
	Compress(data []byte) ([]byte, error)
	// 解压
	UnCompress(data []byte) ([]byte, error)
}

// GetCompress 根据压缩类型获取压缩算法
func GetCompress(compressType string) ICompress {
	switch compressType {
	case "gzip", "gz":
		return GzipCompress
	default:
		return nil
	}
}
