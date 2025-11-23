package md5

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func MD5WithFile(file io.Reader) string {
	//新哈希计算实例
	m := md5.New()
	//复制到哈希对象并自动计算哈希值字节切片
	io.Copy(m, file)
	//获取计算好的哈希结果，这里sum实际上是没有计算的，summarize的意思是获取
	sum := m.Sum(nil) //nil是不追加到任何切片后面
	//转化为十六进制字符串
	return hex.EncodeToString(sum)
}
