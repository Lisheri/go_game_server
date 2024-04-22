package utils

import "math/rand"

var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandSeq(n int) string {
	// ? 用make生成切片, 返回第一个参数指定类型的实例, 比如这里生成 rune[] 类型切片(数组), 长度为n
	b := make([]rune, n);
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}