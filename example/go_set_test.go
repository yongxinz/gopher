package main

import (
	"testing"
)

const num = int(1 << 24)

// 测试 bool 类型
func Benchmark_SetWithBoolValueWrite(b *testing.B) {
	set := make(map[int]bool)
	for i := 0; i < num; i++ {
		set[i] = true
	}
}

// 测试 interface{} 类型
func Benchmark_SetWithInterfaceValueWrite(b *testing.B) {
	set := make(map[int]interface{})
	for i := 0; i < num; i++ {
		set[i] = struct{}{}
	}
}

// 测试 int 类型
func Benchmark_SetWithIntValueWrite(b *testing.B) {
	set := make(map[int]int)
	for i := 0; i < num; i++ {
		set[i] = 0
	}
}

// 测试 struct{} 类型
func Benchmark_SetWithStructValueWrite(b *testing.B) {
	set := make(map[int]struct{})
	for i := 0; i < num; i++ {
		set[i] = struct{}{}
	}
}
