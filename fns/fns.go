package fns

import (
	"io"

	"github.com/bytedance/gopkg/lang/fastrand"
)

// Ptr 返回指向给定值的指针。
func Ptr[T any](v T) *T {
	return &v
}

// SliceToMap ToMap 切片转映射。
func SliceToMap[K comparable, T any](s []T, f func(T) K) map[K]T {
	m := make(map[K]T, len(s))
	for _, v := range s {
		m[f(v)] = v
	}
	return m
}

// Filter 按函数过滤并返回符合条件的切片子元素。
func Filter[T any](s []T, f func(T) bool) []T {
	var r []T
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

// FilterParam 按指定的参数及函数过滤并返回符合条件的切片子元素。
func FilterParam[T any, P any](s []T, p P, f func(T, P) bool) []T {
	var r []T
	for _, v := range s {
		if f(v, p) {
			r = append(r, v)
		}
	}
	return r
}

// Map 将给定的函数应用于切片中的每个元素，并返回一个带有结果的新切片。
func Map[T, U any](s []T, f func(T) U) []U {
	r := make([]U, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return r
}

// SelectRandom 从切片中随机选出指定数量的子元素。
func SelectRandom[T any](s []T, n int) []T {
	if len(s) <= n {
		return s
	}
	//rand.Seed(uint64(time.Now().UnixNano()))
	r := make([]T, n)
	for i := range r {
		// 从切片中选择并弹出一个随机索引
		idx := fastrand.Intn(len(s))
		// 放入新切片
		r[i] = s[idx]
		// 剔除掉选中的索引
		s = append(s[:idx], s[idx+1:]...)
	}
	return r
}

// Unique 按需返回切片中的唯一值。
func Unique[T comparable](s []T) []T {
	m := make(map[T]struct{})
	var r []T
	for _, v := range s {
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			r = append(r, v)
		}
	}
	return r
}

// Any 判断切片中是否有任一个元素能够满足指定函数的要求。
func Any[T any](s []T, f func(T) bool) bool {
	for _, v := range s {
		if f(v) {
			return true
		}
	}
	return false
}

// All 判断切片中是否有任一个元素能够满足指定函数的要求。
func All[T any](s []T, f func(T) bool) bool {
	for _, v := range s {
		if !f(v) {
			return false
		}
	}
	return true
}

func CloseIgnore(stream io.Closer) {
	_ = stream.Close()
}
