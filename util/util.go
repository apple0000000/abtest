package util

import (
	"encoding/json"
	"strconv"
)

// SliceContainsString 检查字符串切片是否包含某字符串 / Check if a string slice contains a string
func SliceContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// SliceContainsInt64 检查int64切片是否包含某int64 / Check if an int64 slice contains an int64
func SliceContainsInt64(slice []int64, item int64) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// GetAssertString 将任意类型转换为字符串 / Convert any type to string
func GetAssertString(v interface{}) string {
	switch v := v.(type) {
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case string:
		return v
	default:
		return ""
	}
}

// ToJson 将对象转换为JSON字符串 / Convert object to JSON string
func ToJson(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(bytes)
}
