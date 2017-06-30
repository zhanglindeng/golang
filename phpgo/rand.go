package phpgo

import (
	"math/rand"
	"time"
)

// RandFloat64 generate random float 64 bit
func RandFloat64() float64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()
}

// RandFloat generate random float 32
func RandFloat() float32 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float32()
}

// RandInt generate random intger number
func RandInt(max int) int {
	// 不设置随机种子，得到的结果是一样的
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}
