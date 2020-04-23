package utils

import (
	"fmt"
	"strings"
	"unsafe"
)

// 快速转为字符串
func ToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// 删除空格
func RemoveSpaces(s string) string {
	return strings.Join(strings.Fields(s), "")
}

// 用:号连接两个部分，如果后一部分也存在的话
func ConcatWith(master, slave string) string {
	if slave != "" {
		master += ":" + slave
	}
	return master
}

// 如果本身不为空，在左右两边添加字符
func WrapWith(s, left, right string) string {
	if s == "" {
		return ""
	}
	return fmt.Sprintf("%s%s%s", left, s, right)
}
