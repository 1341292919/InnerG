package utils

import (
	"InnerG/pkg/constants"
	"math/rand"
	"time"
)

func GenerateRandomCode(length int) string {
	// 初始化随机数生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	code := make([]byte, length)
	for i := range code {
		code[i] = constants.CharSet[r.Intn(len(constants.CharSet))]
	}

	return string(code)
}
