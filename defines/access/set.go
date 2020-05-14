package access

// 操作列表
const (
	VIEW uint16 = 2 << iota
	DISABLE
	REMOVE
	EDIT
	CREATE
	GET
	POST
	GRANT
	ALL
	NONE uint16 = 0 // 无权限
)

var (
	AccessNames = map[uint16]string{
		VIEW: "view", DISABLE: "disable", REMOVE: "remove",
		EDIT: "edit", CREATE: "create", GET: "get", POST: "post",
		GRANT: "grant", ALL: "all", NONE: "",
	}
	AccessTitles = map[uint16]string{
		VIEW: "查看", DISABLE: "禁用", REMOVE: "删除",
		EDIT: "编辑", CREATE: "新建", GET: "GET", POST: "POST",
		GRANT: "授权", ALL: "全部", NONE: "无",
	}
)

// 分解出具体权限
func ParsePermNames(perm uint16) (codes []uint16, names []string) {
	for code, name := range AccessNames {
		if code > 0 && perm&code == code {
			codes = append(codes, code)
			names = append(names, name)
		}
	}
	return
}

// 找出权限的中文名称
func GetPermTitles(codes []uint16) (titles []string) {
	title, ok := "", false
	for _, code := range codes {
		if title, ok = AccessTitles[code]; !ok {
			title = "未知"
		}
		titles = append(titles, title)
	}
	return
}
