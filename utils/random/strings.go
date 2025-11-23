package random

import "math/rand"

var letters = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStr(n int) string {
	//定义长度为n的rune切片存储结果
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))] //中括号里面是以letters长度为界限，生成0到这个界限的随机数，作为索引，实现字符串随机
	}
	return string(b)
}
