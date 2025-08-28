package util

import (
	"errors"
	"strconv"
	"strings"
)

func Partition[T any](src []T, chunkSize int) ([][]T, error) {
	if chunkSize <= 0 {
		return nil, errors.New("chunkSize 必须大于 0")
	}
	if len(src) == 0 {
		return [][]T{}, nil
	}

	totalChunks := (len(src) + chunkSize - 1) / chunkSize
	result := make([][]T, 0, totalChunks)

	for i := 0; i < len(src); i += chunkSize {
		end := i + chunkSize
		if end > len(src) {
			end = len(src)
		}
		result = append(result, src[i:end])
	}
	return result, nil
}

func SplitNonEmpty(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, sep)
	res := make([]string, 0, len(parts))

	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			res = append(res, trimmed)
		}
	}
	return res
}

func SplitToInt64Slice(s, seq string) []int64 {
	if s == "" {
		return []int64{}
	}
	var res []int64
	for _, part := range strings.Split(s, seq) {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			num, err := strconv.ParseInt(trimmed, 10, 64)
			if err == nil {
				res = append(res, num)
			}
		}
	}
	return res
}

func ContainsNoCase(s, substr string) bool {
	return strings.Contains(strings.ToUpper(s), strings.ToUpper(substr))
}

func ContainsAnySplitNoCase(s, substr string) bool {
	splitStr := strings.Split(substr, ",")

	var validStr []string
	for _, ss := range splitStr {
		if trimmed := strings.TrimSpace(ss); trimmed != "" {
			validStr = append(validStr, trimmed)
		}
	}

	str := strings.ToUpper(s)
	for _, vs := range validStr {
		sub := strings.ToUpper(vs)
		if strings.Contains(str, sub) {
			return true
		}
	}
	return false
}

func Contains[T comparable](slice []T, item T) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func MapKeyToSlice[K comparable, V any](m map[K]V) []K {
	if m == nil {
		return []K{}
	}

	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

func FlattenSlice[T any](twoD [][]T) []T {
	if twoD == nil {
		return []T{}
	}
	size := 0
	for _, row := range twoD {
		size += len(row)
	}
	oneD := make([]T, 0, size)
	for _, row := range twoD {
		oneD = append(oneD, row...)
	}
	return oneD
}

func HasIntersection[T comparable](a, b []T) bool {
	// 创建集合存储第一个切片的元素
	set := make(map[T]struct{}, len(a))
	for _, v := range a {
		set[v] = struct{}{}
	}

	// 检查第二个切片的元素是否在集合中
	for _, v := range b {
		if _, exists := set[v]; exists {
			return true
		}
	}
	return false
}
