package utils

import (
	"sort"
	"strings"
)

// 字符串比较方式
const (
	CMP_STRING_OMIT             = iota // 不比较
	CMP_STRING_CONTAINS                // 包含
	CMP_STRING_STARTSWITH              // 打头
	CMP_STRING_ENDSWITH                // 结尾
	CMP_STRING_IGNORE_SPACES           // 忽略空格
	CMP_STRING_CASE_INSENSITIVE        // 不分大小写
	CMP_STRING_EQUAL                   // 相等
)

// 比较是否相符
func StringMatch(a, b string, cmp int) bool {
	switch cmp {
	case CMP_STRING_OMIT:
		return true
	case CMP_STRING_CONTAINS:
		return strings.Contains(a, b)
	case CMP_STRING_STARTSWITH:
		return strings.HasPrefix(a, b)
	case CMP_STRING_ENDSWITH:
		return strings.HasSuffix(a, b)
	case CMP_STRING_IGNORE_SPACES:
		a, b = RemoveSpaces(a), RemoveSpaces(b)
		return strings.EqualFold(a, b)
	case CMP_STRING_CASE_INSENSITIVE:
		return strings.EqualFold(a, b)
	default: // 包括 CMP_STRING_EQUAL
		return strings.Compare(a, b) == 0
	}
}

// 是否在字符串列表中，只适用于CMP_STRING_EQUAL和CMP_STRING_STARTSWITH
func compareStringList(x string, lst []string, cmp int) bool {
	size := len(lst)
	if size == 0 {
		return false
	}
	if !sort.StringsAreSorted(lst) {
		sort.Strings(lst)
	}
	i := sort.Search(size, func(i int) bool { return lst[i] >= x })
	return i < size && StringMatch(x, lst[i], cmp)
}

// 是否在字符串列表中
func InStringList(x string, lst []string) bool {
	return compareStringList(x, lst, CMP_STRING_EQUAL)
}

// 是否在字符串列表中，比较方式是有任何一个开头符合
func StartStringList(x string, lst []string) bool {
	return compareStringList(x, lst, CMP_STRING_STARTSWITH)
}

// lst1 是否 lst2 的（真）子集
func IsSubsetList(lst1, lst2 []string, strict bool) bool {
	if len(lst1) > len(lst2) {
		return false
	}
	if strict && len(lst1) == len(lst2) {
		return false
	}
	for _, x := range lst1 {
		if !InStringList(x, lst2) {
			return false
		}
	}
	return true
}
