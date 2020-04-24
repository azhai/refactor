package rewrite

import (
	"bytes"
	"fmt"
	"go/ast"
	"sort"
	"strings"

	"gitea.com/azhai/refactor/utils"
)

const MODEL_EXTENDS = "`xorm:\"extends\"`"

var substituteModels = map[string]*ClassSummary{
	"base.TimeModel": {
		Name:       "base.TimeModel",
		Import:     "gitea.com/azhai/refactor/language/common",
		Alias:      "base",
		Features: []string{
			"CreatedAt time.Time",
			"UpdatedAt time.Time",
			"DeletedAt time.Time",
		},
		FieldLines: []string{
			"CreatedAt time.Time `json:\"created_at\" xorm:\"created comment('创建时间') TIMESTAMP\"`       // 创建时间",
			"UpdatedAt time.Time `json:\"updated_at\" xorm:\"updated comment('更新时间') TIMESTAMP\"`       // 更新时间",
			"DeletedAt time.Time `json:\"deleted_at\" xorm:\"deleted comment('删除时间') index TIMESTAMP\"` // 删除时间",
		},
	},
	"base.NestedModel": {
		Name:       "base.NestedModel",
		Import:     "gitea.com/azhai/refactor/language/common",
		Alias:      "base",
		Features: []string{
			"Lft int",
			"Rgt int",
			"Depth int",
		},
		FieldLines: []string{
			"Lft   int `json:\"lft\" xorm:\"not null default 0 comment('左边界') INT(10)\"`           // 左边界",
			"Rgt   int `json:\"rgt\" xorm:\"not null default 0 comment('右边界') index INT(10)\"`     // 右边界",
			"Depth int `json:\"depth\" xorm:\"not null default 1 comment('高度') index TINYINT(3)\"` // 高度",
		},
	},
}

func RegisterSubstitute(sub *ClassSummary) {
	if sub != nil {
		substituteModels[sub.Name] = sub
	}
}

type ClassSummary struct {
	Name           string
	Substitute     string
	Import, Alias  string
	Features       []string
	sortedFeatures []string
	FieldLines     []string
	IsChanged      bool
}

func NewClassSummary(name string) *ClassSummary {
	return &ClassSummary{Name: name}
}

func (s ClassSummary) GetInnerCode() string {
	var buf bytes.Buffer
	for _, line := range s.FieldLines {
		buf.WriteString(fmt.Sprintf("\t%s\n", line))
	}
	return buf.String()
}

func (s ClassSummary) GetSortedFeatures() []string {
	if len(s.sortedFeatures) == 0 {
		s.sortedFeatures = append([]string{}, s.Features...)
		sort.Strings(s.sortedFeatures)
	}
	return s.sortedFeatures
}

func (s *ClassSummary) GetSubstitute() string {
	if s.Substitute == "" {
		s.Substitute = fmt.Sprintf("*%s %s", s.Name, MODEL_EXTENDS)
	}
	return s.Substitute
}

func (s *ClassSummary) ParseFields(cp *CodeParser, node *DeclNode) int {
	size := len(node.Fields)
	s.Features = make([]string, size)
	s.FieldLines = make([]string, size)
	for i, f := range node.Fields {
		code := cp.GetNodeCode(f)
		ps := strings.Fields(code)
		if len(ps) == 0 {
			continue
		}
		if len(ps) == 1 {
			s.Features[i] = ps[0]
		} else {
			s.Features[i] = ps[0] + " " + ps[1]
		}
		if cm := cp.GetComment(f.Comment, true); len(cm) > 0 {
			code += " //" + cm
		}
		s.FieldLines[i] = code
	}
	return size
}

func ReplaceModelFields(cp *CodeParser, node *DeclNode, summary *ClassSummary) {
	var last ast.Node
	max := len(node.Fields) - 1
	first, lastField := ast.Node(node.Fields[0]), node.Fields[max]
	if lastField.Comment != nil {
		last = ast.Node(lastField.Comment)
	} else {
		last = ast.Node(lastField)
	}
	cp.AddReplace(first, last, summary.GetInnerCode())
}

func ReplaceSummary(summary, sub *ClassSummary) *ClassSummary {
	var features, lines []string
	find := false
	sted := sub.GetSortedFeatures()
	for i, ft := range summary.Features {
		if !utils.InStringList(ft, sted, utils.CMP_STRING_EQUAL) {
			features = append(features, ft)
			lines = append(lines, summary.FieldLines[i])
		} else if !find {
			subst := sub.GetSubstitute()
			features = append(features, subst)
			lines = append(lines, subst)
			find = true
			summary.IsChanged = true
		}
	}
	summary.Features, summary.FieldLines = features, lines
	return summary
}

func ScanModelDir(dir string, verbose bool) error {
	files, _ := utils.FindFiles(dir, ".go")
	for fname := range files {
		cp, err := NewFileParser(fname)
		if err != nil {
			fmt.Println(fname, " error: ", err)
			continue
		}
		var changed bool
		imports := make(map[string]string)
		for _, node := range cp.AllDeclNode("type") {
			if len(node.Fields) == 0 {
				continue
			}
			name := node.GetName()
			if strings.Contains(cp.GetNodeCode(node), MODEL_EXTENDS) {
				continue // 避免重复处理
			}

			summary := NewClassSummary(name)
			_ = summary.ParseFields(cp, node)
			for n, sub := range substituteModels {
				if n == summary.Name {
					continue // 不要替换自己
				}
				sted := sub.GetSortedFeatures()
				sorted := summary.GetSortedFeatures()
				if utils.IsStrictSubsetList(sted, sorted) {
					summary = ReplaceSummary(summary, sub)
					imports[sub.Import] = sub.Alias
					if verbose {
						fmt.Println(summary.Name, " <- ", sub.Name)
					}
				} else if strings.HasPrefix(n, "base.") || n == summary.Name {
					continue
				} else if utils.IsStrictSubsetList(sorted, sted) {
					ReplaceSummary(sub, summary)
					if verbose {
						fmt.Println(sub.Name, " <- ", summary.Name)
					}
				}
			}
			RegisterSubstitute(summary)
			if summary.IsChanged {
				changed = true
				ReplaceModelFields(cp, node, summary)
			}
		}
		if verbose {
			fmt.Println(fname, " changed: ", changed, "\n")
		}
		if changed {
			cs := cp.CodeSource
			if code, chg := cs.AltSource(); chg {
				cs.SetSource(code)
			}
			if cs, err = ResetImports(cs, imports); err != nil {
				return err
			}
			if err = cs.WriteTo(fname); err != nil {
				return err
			}
		}
	}
	return nil
}
